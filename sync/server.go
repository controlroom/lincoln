package sync

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/controlroom/lincoln/backends/docker"
	"github.com/controlroom/lincoln/interfaces"
	"github.com/fsnotify/fsevents"
)

type watcherData struct {
	eventStream *fsevents.EventStream
	closer      chan bool
}

type RPCSync struct {
	watches map[string]watcherData
	res     chan ProjectSyncInfo
	mu      *sync.RWMutex
}

type ProjectSyncInfo struct {
	Backend interfaces.Operation
	Name    string
	Path    string
}

func NewSync(res chan ProjectSyncInfo) *RPCSync {
	return &RPCSync{
		watches: make(map[string]watcherData),
		res:     res,
		mu:      &sync.RWMutex{},
	}
}

func (s *RPCSync) Watch(info *ProjectSyncInfo, ack *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.watches[info.Name]; !ok {
		fmt.Printf("Recieved new project: %v\n", *info)

		closer := make(chan bool, 1)
		es, ec := getEvents(info.Path)
		s.watches[info.Name] = watcherData{
			eventStream: es,
			closer:      closer,
		}

		go func() {
		replay:
			select {
			case msg := <-ec:
				for _, event := range msg {
					s.res <- ProjectSyncInfo{
						Backend: info.Backend,
						Name:    info.Name,
						Path:    event.Path,
					}
				}
				goto replay
			case <-closer:
				fmt.Println(fmt.Sprintf("Closing %v", info.Name))
			}
		}()
	}

	return nil
}

func (s *RPCSync) UnWatch(info *ProjectSyncInfo, ack *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if data, ok := s.watches[info.Name]; ok {
		data.eventStream.Stop()
		data.closer <- true
		delete(s.watches, info.Name)
	}

	return nil
}

func listenForCleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		os.Exit(0)
	}()
}

func StartServer(port int) {
	listenForCleanup()

	res := make(chan ProjectSyncInfo, 1)
	sync := NewSync(res)
	rpc.Register(sync)

	l, e := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if e != nil {
		log.Fatal("listen error:", e)
	}

	go func() {
		for info := range res {
			info.Backend.Sync(info.Name, info.Path)
		}
	}()

	fmt.Println("Sync server started")
	rpc.Accept(l)
}

func allDirs(path string, dirs *[]string, ignored map[string]bool) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if ignored[file.Name()] != false {
			continue
		}
		if file.IsDir() {
			finalPath := filepath.Join(path, file.Name())
			*dirs = append(*dirs, finalPath)
			allDirs(finalPath, dirs, ignored)
		}
	}
}

func getEvents(path string) (*fsevents.EventStream, chan []fsevents.Event) {
	dirs := []string{path}
	allDirs(path, &dirs, map[string]bool{".git": true})

	es := &fsevents.EventStream{
		Paths:   dirs,
		Latency: 250 * time.Millisecond,
		Flags:   4,
	}
	es.Start()
	return es, es.Events
}

func init() {
	gob.Register(docker.DockerOperation{})
}
