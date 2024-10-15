package websock

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

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

const BUFFER_SIZE = 1024 * 32

func (client *WebSockClient) Read() ([]byte, error) {
	var readBuffer []byte = make([]byte, BUFFER_SIZE)
	log.Println("READING")

	// _, rdrConn, err := client.conn.Reader(client.ctx)
	typ, rdrConn, err := client.conn.Reader(client.ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Type", typ)
	rdr := bufio.NewReaderSize(rdrConn, BUFFER_SIZE)

	for {
		// _, err := rdr.Peek(1)
		count, err := rdr.Read(readBuffer)
		if err != nil && err != io.EOF {
			log.Println("read error", err)
			return nil, err
		}

		if count > 0 {
			log.Println("read count", count)
			return readBuffer[:count], nil
		}

		time.Sleep(time.Millisecond)
	}

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

func listen() {
	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			io.Copy(c, c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}
