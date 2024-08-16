package web

import (
	"log"
	"net/http"
	"v4lvid/video"
)

var _ http.Handler = (*RecordControlHandler)(nil)

type RecordControlHandler struct {
	Server *video.Server
	Url    string
	Icon   string
}

func (ctl *RecordControlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("recording...")
	ctl.Server.RecordCmd(60)
}
