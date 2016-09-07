package sync

import (
	"net"
	"net/rpc"
	"time"

	"github.com/controlroom/lincoln/interfaces"
)

type Client struct {
	connection *rpc.Client
}

func NewClient(dsn string, timeout time.Duration) (*Client, error) {
	connection, err := net.DialTimeout("tcp", dsn, timeout)
	if err != nil {
		return nil, err
	}
	return &Client{connection: rpc.NewClient(connection)}, nil
}

func (c *Client) Watch(backend interfaces.Operation, name string, path string) (bool, error) {
	var added bool
	info := ProjectSyncInfo{
		Backend: backend,
		Name:    name,
		Path:    path,
	}
	err := c.connection.Call("RPCSync.Watch", info, &added)
	return added, err
}

func (c *Client) UnWatch(name string) (bool, error) {
	var removed bool
	info := ProjectSyncInfo{
		Name: name,
	}
	err := c.connection.Call("RPCSync.UnWatch", info, &removed)
	return removed, err
}
