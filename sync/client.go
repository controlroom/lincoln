package sync

import (
	"net"
	"net/rpc"
	"time"
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

func (c *Client) Watch(name string, path string) (bool, error) {
	var added bool
	info := ProjectSyncInfo{
		Name: name,
		Path: path,
	}
	err := c.connection.Call("RPCSync.Watch", info, &added)
	return added, err
}