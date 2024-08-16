package web

import (
	"log"
	"net/http"
	"v4lvid/video"
)

type ControlList struct {
	webcam   *video.Webcam
	Id       int
	Controls []*ControlHandler
}

func NewControlList(webcam *video.Webcam, id int, controls []*ControlHandler) *ControlList {
	ctlh := &ControlList{
		webcam:   webcam,
		Id:       id,
		Controls: make([]*ControlHandler, 0, len(controls)),
	}
	for _, ctl := range controls {
		ctlh.AddControl(ctl)
	}
	return ctlh
}

func (ctlh *ControlList) AddControl(ctl *ControlHandler) {
	var err error
	if ctl == nil {
		log.Fatalln("AddControl control is nil")
	}
	ctl.webcam = ctlh.webcam
	ctl.Info, err = ctl.webcam.GetControlInfo(ctl.Key)

	ctlh.Controls = append(ctlh.Controls, ctl)
	if err != nil {
		log.Println("AddControl", err)
	}
	http.Handle(ctl.Url, ctl)
}
