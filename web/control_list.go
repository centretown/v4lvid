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
	ctll := &ControlList{
		webcam:   webcam,
		Id:       id,
		Controls: make([]*ControlHandler, 0, len(controls)),
	}
	for _, ctl := range controls {
		ctll.AddControl(ctl)
	}
	return ctll
}

func (ctll *ControlList) AddControl(ctlh *ControlHandler) {
	var err error
	if ctlh == nil {
		log.Fatalln("AddControl control is nil")
	}
	ctlh.webcam = ctll.webcam
	ctlh.Info, err = ctlh.webcam.GetControlInfo(ctlh.Key)
	ctlh.Value = ctlh.webcam.GetControlValue(ctlh.Key)

	ctll.Controls = append(ctll.Controls, ctlh)
	if err != nil {
		log.Println("AddControl", err)
	}

	for _, ctl := range ctlh.Controls {
		http.Handle(ctl.Url, ctlh)
	}

}
