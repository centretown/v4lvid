package web

import (
	"log"
	"net/http"
	"time"
	"v4lvid/video"
)

var _ http.Handler = (*RecordControlHandler)(nil)

type RecordControlHandler struct {
	Server    *video.Server
	Url       string
	Icon      string
	recording bool
}

func (ctl *RecordControlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !ctl.recording {
		log.Printf("recording...")
		ctl.recording = true
		ctl.Server.RecordCmd(60)
		time.AfterFunc(60*time.Second, func() {
			log.Println("timer finished", ctl.recording)
			ctl.recording = false
		})

	} else {
		ctl.Server.StopRecordCmd()
		ctl.recording = false
	}
}
