package websock

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

// go get github.com/gorilla/websocket

const (
	CONFIG = 20 + iota
	STATES
	SUBSCRIBE
)

const (
	melonUrl = "ws://melon:8123/api/websocket"
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
	wsc := &WebSockClient{
		ctx: context.Background(),
		// conn:      conn,
		Quit:      make(chan int),
		Buffer:    make(chan []byte),
		MessageID: STATES,
	}
	// ctx := context.Background()
	conn, resp, err := websocket.DefaultDialer.Dial(melonUrl, nil)
	// conn, resp, err := websocket.Dial(ctx, dial, nil)
	if err != nil {
		return nil, err
	}
	wsc.conn = conn
	log.Println("Dial", resp.Status)
	// hs := &WebSockClient{
	// 	ctx:       ctx,
	// 	conn:      conn,
	// 	Quit:      make(chan int),
	// 	Buffer:    make(chan []byte),
	// 	MessageID: STATES,
	// }
	return wsc, err
}

const BUFFER_SIZE = 1024 * 32

func (wsc *WebSockClient) Read() (readBuffer []byte, err error) {
	var typ int
	typ, readBuffer, err = wsc.conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Type", typ)
	return
	// // var readBuffer []byte = make([]byte, BUFFER_SIZE)

	// _, rdrConn, err := hs.conn.ReadMessage() //.Reader(hs.ctx)
	// // typ, rdrConn, err := hs.conn.Reader(hs.ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// // log.Println("Type", typ)
	// rdr := bufio.NewReaderSize(rdrConn, BUFFER_SIZE)

	// for {

	// 	count, err := rdr.Read(readBuffer)
	// 	if err != nil && err != io.EOF {
	// 		return nil, err
	// 	}

	// 	if count > 0 {
	// 		// log.Println("read count", count)
	// 		return readBuffer[:count], nil
	// 	}
	// }

}

func (wsc *WebSockClient) Write(cmd string) error {
	err := wsc.conn.WriteMessage(1, []byte(cmd))
	// err := hs.conn.Write(hs.ctx, websocket.MessageText, []byte(cmd))
	if err != nil {
		log.Println(cmd)
		log.Println(err)
	}
	return err
}

func (wsc *WebSockClient) WriteID(cmd string) (id int, err error) {
	// id = hs.MessageID
	// message := fmt.Sprintf(cmd, id)
	// hs.MessageID++
	// err = hs.conn.Write(hs.ctx, websocket.MessageText, []byte(message))
	// if err != nil {
	// 	log.Println(message)
	// 	log.Println(err)
	// }
	return
}
