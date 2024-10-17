package websock

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/coder/websocket"
)

const (
	CONFIG = 20 + iota
	STATES
	SUBSCRIBE
)

const (
	dial = "ws://melon:8123/api/websocket"
)

type WebSockClient struct {
	conn      *websocket.Conn
	ctx       context.Context
	MessageID int
	Err       error
	Quit      chan int
	Buffer    chan []byte
}

func NewWebSockClient() (*WebSockClient, error) {
	ctx := context.Background()
	conn, resp, err := websocket.Dial(ctx, dial, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Dial", resp.Status)
	client := &WebSockClient{
		ctx:       ctx,
		conn:      conn,
		Quit:      make(chan int),
		Buffer:    make(chan []byte),
		MessageID: STATES,
	}
	return client, err
}

func (client *WebSockClient) Read() (buf []byte, err error) {
	// waits until something is there
	_, rdrConn, err := client.conn.Reader(client.ctx)
	if err != nil {
		log.Println("Reader", err)
		return nil, err
	}
	return io.ReadAll(rdrConn)
}

func (client *WebSockClient) WriteCommand(cmd string) error {
	err := client.conn.Write(client.ctx, websocket.MessageText, []byte(cmd))
	if err != nil {
		log.Println("Write", cmd, err)
	}
	return err
}

func (client *WebSockClient) WriteCommandID(cmd string) (id int, err error) {
	id = client.MessageID
	message := fmt.Sprintf(cmd, id)
	err = client.conn.Write(client.ctx, websocket.MessageText, []byte(message))
	if err != nil {
		log.Println("WriteID", message, err)
	}
	client.MessageID += 1
	return
}
