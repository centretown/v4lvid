package sockserve

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type SockServer struct {
	messageChan   chan string
	Messages      []*Message
	mutex         sync.Mutex
	hub           *Hub
	statusLayout  *template.Template
	messageLayout *template.Template
}

func NewSockServer(t *template.Template) *SockServer {
	s := &SockServer{
		messageChan:   make(chan string),
		Messages:      make([]*Message, 0),
		hub:           NewHub(),
		statusLayout:  t.Lookup("layout.wsstatus"),
		messageLayout: t.Lookup("layout.wsmessage"),
	}
	return s
}

func (s *SockServer) Run() {
	go s.hub.Run()
}

func (s *SockServer) PastMessages() (past []*Message) {
	max := len(s.Messages)
	past = make([]*Message, max)
	for i := range max {
		past[max-i-1] = s.Messages[i]
	}
	return
}

const messageFile = "messages.json"

func (s *SockServer) LoadMessages() (err error) {
	var buf []byte
	buf, err = os.ReadFile(messageFile)
	if err != nil {
		log.Println("LoadMessages:ReadFile", err)
		return
	}

	err = json.Unmarshal(buf, &s.Messages)
	if err != nil {
		log.Println("LoadMessages:Unmarshal", err)
		return
	}
	return
}

func (s *SockServer) SaveMessages() (err error) {
	var buf []byte
	buf, err = json.MarshalIndent(s.Messages, "", "  ")
	if err != nil {
		log.Println("SaveMessages:MarshalIndent", err)
		return
	}
	err = os.WriteFile(messageFile, buf, os.ModePerm)
	if err != nil {
		log.Println("SaveMessages:WriteFile", err)
		return
	}
	return
}

const (
	messageOff = `<span id="streamer" hx-swap-oob="outerHTML" class="symbols">radio_button_checked</span>`
	messageOn  = `<span id="streamer" hx-swap-oob="outerHTML" class="symbols streaming">radio_button_checked</span>`
)

func (s *SockServer) StreamOn() {
	s.hub.Broadcast <- []byte(messageOn)
	// w.WriteHeader(http.StatusOK)
}

func (s *SockServer) StreamOff() {
	s.hub.Broadcast <- []byte(messageOff)
	// w.WriteHeader(http.StatusOK)
}

func (s *SockServer) Webhook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed to parse body: %v", err)
		return
	}

	name := "unknown"
	names, ok := r.PostForm["name"]
	if ok && len(names) > 0 {
		name = names[0]
	}

	message := "empty"
	messages, ok := r.PostForm["message"]
	if ok && len(messages) > 0 {
		message = messages[0]
	}

	log.Printf("Received webhook: %s %s", name, message)

	var (
		buf bytes.Buffer
		msg = &Message{Name: name, Message: message, Stamp: time.Now()}
	)

	err = s.messageLayout.Execute(&buf, msg)
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// wrap the message in a div so we can use htmx to add it to the page
	s.hub.Broadcast <- []byte("<div hx-swap-oob=\"afterbegin:#messages\">" + buf.String() + "</div>")

	s.mutex.Lock()
	s.Messages = append(s.Messages, msg)
	if len(s.Messages) > maxMessages {
		s.Messages = s.Messages[1:]
	}
	log.Printf("Now have %d past messages", len(s.Messages))
	s.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (s *SockServer) Events(w http.ResponseWriter, r *http.Request) {
	client, err := NewClient(s.hub, w, r)
	if err != nil {
		log.Printf("Failed to create WebSocket client: %v", err)
		return
	}

	s.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

// func (s *SockServer) Status(w http.ResponseWriter, r *http.Request) {
// 	log.Println("STATUS")
// 	s.statusLayout.Execute(
// 		w,
// 		struct {
// 			WebsocketHost string
// 			ClientList    string
// 			PastMessages  []string
// 		}{
// 			ClientList:    s.hub.GetClientList(),
// 			WebsocketHost: "ws://" + r.Host + "/events",
// 			PastMessages:  s.PastMessages,
// 		},
// 	)
// }
