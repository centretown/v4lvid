package web

import (
	"log"
	"net/http"
	"v4lvid/camera"
	"v4lvid/sockserve"
)

var _ http.Handler = (*Streamer)(nil)

type Streamer struct {
	Server *camera.Server
	Url    string
	Icon   string
	Sock   *sockserve.SockServer
}

func (strm *Streamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strm.Server.Recording {
		log.Printf("recording...")
		strm.Server.RecordCmd(300)
	} else {
		log.Printf("stop recording...")
		strm.Server.StopRecordCmd()
	}
}

func (strm *Streamer) IsRecording() bool {
	return strm.Server.Recording
}
