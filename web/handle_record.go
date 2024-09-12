package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"v4lvid/camera"
)

var _ http.Handler = (*HandleRecord)(nil)

type HandleRecord struct {
	vserve *camera.Server
}

func NewHandleRecord(vserve *camera.Server) *HandleRecord {
	hr := &HandleRecord{
		vserve: vserve,
	}
	return hr
}

func (hr *HandleRecord) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if hr.vserve == nil {
		log.Println("no camera")
		return
	}
	log.Println("record requested", r.Host)
	if !hr.vserve.Busy {
		log.Println("cam is idle", r.Host)
		return
	}

	duration := 5
	values := r.URL.Query()
	parm, ok := values["duration"]
	if ok && len(parm) > 0 {
		i, err := strconv.Atoi(parm[0])
		if err == nil {
			duration = i
		}
	}

	log.Println("request values", len(values), values)
	for k, v := range values {
		log.Println("request values", k, v)
	}

	var (
		cmd     camera.ServerCmd
		message string
	)
	if hr.vserve.Recording {
		cmd.Action = camera.RECORD_STOP
		message = "stop"
	} else {
		cmd.Action = camera.RECORD_START
		message = fmt.Sprintln("record for", duration, "seconds")
	}

	cmd.Value = duration
	hr.vserve.Cmd <- cmd

	w.Write([]byte(message))
}
