package sync

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/controlroom/lincoln/interfaces"
	"github.com/controlroom/lincoln/metadata"
)

type Client struct {
	connection *rpc.Client
}

func GetClient() (*Client, error) {
	port := metadata.GetMeta("app:syncPort")
	if port == "" {
		return nil, errors.New("Watch client not running")
	}

	timeout := time.Millisecond * 500
	dsn := fmt.Sprintf("localhost:%v", port)

	connection, err := net.DialTimeout("tcp", dsn, timeout)
	if err != nil {
		return nil, err
	}

	return &Client{connection: rpc.NewClient(connection)}, nil
}

func (c *Client) Watch(backend interfaces.Operation, name string, path string) error {
	var added bool
	info := ProjectSyncInfo{
		Backend: backend,
		Name:    name,
		Path:    path,
	}
	err := c.connection.Call("RPCSync.Watch", info, &added)
	return err
}

func (c *Client) UnWatch(name string) error {
	var removed bool
	info := ProjectSyncInfo{
		Name: name,
	}
	err := c.connection.Call("RPCSync.UnWatch", info, &removed)
	return err
}
