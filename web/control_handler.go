package web

import (
	"log"
	"net/http"
	"v4lvid/video"

	"github.com/korandiz/v4l"
)

var _ http.Handler = (*ControlHandler)(nil)

type ControlHandler struct {
	webcam     *video.Webcam
	Key        string
	Url        string
	Multiplier int32
	Icon       string
	Info       v4l.ControlInfo
	Value      int32
}

func (ctl *ControlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle", ctl.Key, ctl.Url)

	var (
		currentValue int32 = ctl.webcam.GetControlValue(ctl.Key)
		step               = ctl.Info.Step * ctl.Multiplier
	)

	newValue := currentValue + step
	if newValue >= ctl.Info.Min && newValue <= ctl.Info.Max {
		currentValue = newValue
		ctl.webcam.SetValue(ctl.Key, currentValue)
	}
}
