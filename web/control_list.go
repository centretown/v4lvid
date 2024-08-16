package web

import (
	"log"
	"net/http"
	"v4lvid/video"
)

type ControlList struct {
	webcam   *video.Webcam
	Id       int
	Controls []*Control
}

func NewControlList(webcam *video.Webcam, id int, controls []*Control) *ControlList {
	ctlh := &ControlList{
		webcam:   webcam,
		Id:       id,
		Controls: make([]*Control, 0, len(controls)),
	}
	for _, ctl := range controls {
		ctl.webcam = webcam
		ctlh.AddControl(ctl)
		http.Handle(ctl.Url, ctl)
	}
	return ctlh
}

func (ctlh *ControlList) AddControl(ctl *Control) {
	var err error
	if ctl == nil {
		log.Fatalln("AddControl control is nil")
	}
	ctlh.Controls = append(ctlh.Controls, ctl)
	ctl.Info, err = ctl.webcam.GetControlInfo(ctl.Key)
	if err != nil {
		log.Println("AddControl", err)
	}
}
