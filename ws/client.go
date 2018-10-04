package ws

import (
	"github.com/gorilla/websocket"
	"time"
)

type Client struct {
	Conn           *websocket.Conn
	LastUpdateTime time.Time
	closed         bool
}

func (c *Client) UpdateTime() {
	c.LastUpdateTime = time.Now()
}

func (c *Client) Close() {
	c.Conn.Close()
	c.closed = true
}

func (c Client) Closed() bool {
	return c.closed
}
