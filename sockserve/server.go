package sockserve

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"sync"
)

const maxMessages = 10

type SockServer struct {
	messageChan   chan string
	PastMessages  []string
	mutex         sync.Mutex
	hub           *Hub
	statusLayout  *template.Template
	messageLayout *template.Template
}

func NewSockServer(t *template.Template) *SockServer {
	s := &SockServer{
		messageChan:   make(chan string),
		PastMessages:  []string{},
		hub:           NewHub(),
		statusLayout:  t.Lookup("layout.wsstatus"),
		messageLayout: t.Lookup("layout.wsmessage"),
	}
	return s
}

func (s *SockServer) Run() {
	go s.hub.Run()
}

func (s *SockServer) Webhook(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		return
	}

	log.Printf("Received webhook: %s", string(b))

	var buf bytes.Buffer
	err = s.messageLayout.Execute(
		&buf,
		struct {
			Raw string
		}{
			Raw: string(b),
		},
	)
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// wrap the message in a div so we can use htmx to add it to the page
	s.hub.Broadcast <- []byte("<ul hx-swap-oob=\"afterbegin:#messages\"><li class=\"message\">" + buf.String() + "</li></ul>")

	s.mutex.Lock()
	s.PastMessages = append(s.PastMessages, buf.String())
	if len(s.PastMessages) > maxMessages {
		s.PastMessages = s.PastMessages[1:]
	}
	log.Printf("Now have %d past messages", len(s.PastMessages))
	s.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (s *SockServer) Events(w http.ResponseWriter, r *http.Request) {
	log.Println("EVENTS")
	client, err := NewClient(s.hub, w, r)
	if err != nil {
		log.Printf("Failed to create WebSocket client: %v", err)
		return
	}

	s.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (s *SockServer) Status(w http.ResponseWriter, r *http.Request) {
	s.statusLayout.Execute(
		w,
		struct {
			WebsocketHost string
			ClientList    string
			PastMessages  []string
		}{
			ClientList:    s.hub.GetClientList(),
			WebsocketHost: "ws://" + r.Host + "/events",
			PastMessages:  s.PastMessages,
		},
	)
}
