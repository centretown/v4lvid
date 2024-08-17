package web

import (
	"log"
	"net/http"
	"v4lvid/video"

	"github.com/korandiz/v4l"
)

var _ http.Handler = (*ControlHandler)(nil)

type Control struct {
	Url        string
	Icon       string
	Multiplier int32
}

type ControlHandler struct {
	webcam *video.Webcam
	Key    string

	Info     v4l.ControlInfo
	Value    int32
	Controls []*Control
	Map      map[string]*Control
}

func NewControlHandler(key string, ctls []*Control) *ControlHandler {
	ctlh := &ControlHandler{
		Key:      key,
		Controls: ctls,
		Map:      make(map[string]*Control),
	}
	for _, ctl := range ctls {
		ctlh.Map[ctl.Url] = ctl
	}
	return ctlh
}

func (ctlh *ControlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctl, ok := ctlh.Map[r.RequestURI]
	if !ok {
		log.Println("RequestURI not found", ctlh.Key, r.RequestURI)
		return
	}

	log.Println("Handle", ctlh.Key, r.RequestURI)
	newValue := ctlh.Value + ctlh.Info.Step*ctl.Multiplier
	if newValue >= ctlh.Info.Min && newValue <= ctlh.Info.Max {
		ctlh.Value = newValue
		ctlh.webcam.SetValue(ctlh.Key, newValue)
	}
}
