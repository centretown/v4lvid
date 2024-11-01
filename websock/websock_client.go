package websock

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/gorilla/websocket"
)

const (
	CONFIG = 20 + iota
	STATES
	SUBSCRIBE
)

const (
	dial = "ws://melon:8123/api/websocket"
)

//	var upgrader = ws.Upgrader{
//		ReadBufferSize:  1024,
//		WriteBufferSize: 1024,
//	}
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
	dialer := &websocket.Dialer{
		ReadBufferSize: 10_000,
	}
	// conn, resp, err := websocket.DefaultDialer.Dial(dial, nil)
	conn, resp, err := dialer.DialContext(ctx, dial, nil)
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
	_, rdrConn, err := client.conn.NextReader()
	if err != nil {
		log.Println("NextReader", err)
		return nil, err
	}
	buf, err = io.ReadAll(rdrConn)
	if err != nil {
		log.Println("ReadAll", err)
	}
	return
}

func (client *WebSockClient) WriteCommand(cmd string) error {
	w, err := client.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Println("WriteCommand NextWriter", cmd, err)
		return err
	}
	defer w.Close()
	_, err = w.Write([]byte(cmd))
	log.Println("Write", cmd, err)
	return err
}

func (client *WebSockClient) WriteCommandID(cmd string) (id int, err error) {
	id = client.MessageID
	message := fmt.Sprintf(cmd, id)
	var w io.WriteCloser
	w, err = client.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Println("WriteCommandID NextWriter", message, err)
		return
	}
	defer w.Close()
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Println("WriteID", message, err)
		return
	}
	client.MessageID += 1
	return
}
