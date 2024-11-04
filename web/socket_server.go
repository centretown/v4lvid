package web

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
	"v4lvid/config"

	"golang.org/x/time/rate"

	"github.com/coder/websocket"
)

type SocketServer struct {
	server *http.Server
	url    string
	mux    *http.ServeMux
}

func NewSocketServer(data *RunData, cfg *config.Config) *SocketServer {
	wss := &SocketServer{
		url: cfg.WsUrl,
		mux: &http.ServeMux{},
		server: &http.Server{
			ReadTimeout:  0,
			WriteTimeout: 0,
		},
	}
	wss.server.Handler = wss.mux
	echoServer := &EchoServer{
		logf: log.Printf,
	}
	wss.mux.HandleFunc("/echo", handleEcho(data, echoServer))
	// wss.mux.HandleFunc("/chat", handleChat(data, echoServer))
	return wss
}

func handleEcho(data *RunData, echoServer *EchoServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		err := data.template.Lookup("layout.echo").Execute(w, data)
		if err != nil {
			log.Fatal("/chat", err)
		}
		echoServer.ServeHTTP(w, r)
	}
}

func (wss *SocketServer) Run() error {
	l, err := net.Listen("tcp", wss.url)
	if err != nil {
		log.Println("net listen", err)
		return err
	}

	log.Printf("listening on ws://%v", l.Addr())
	return wss.server.Serve(l)
}

type EchoServer struct {
	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
}

func (s EchoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer c.CloseNow()

	// if c.Subprotocol() != "echo" {
	// 	c.Close(websocket.StatusPolicyViolation, "client must speak the echo subprotocol")
	// 	return
	// }

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		err = echo(c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			s.logf("failed to echo with %v: %v", r.RemoteAddr, err)
			return
		}
	}
}

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func echo(c *websocket.Conn, l *rate.Limiter) error {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	ctx := context.Background()
	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to io.ReadAll: %w", err)
	}

	_, err = w.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to io.Write: %w", err)
	}
	log.Println(string(buf))
	err = w.Close()
	return err
}
