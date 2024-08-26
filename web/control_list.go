package web

import (
	"log"
	"net/http"
	"v4lvid/video"
)

type ControlList struct {
	webcam   *video.Webcam
	Id       int
	Handlers []*V4lHandler
}

func NewControlList(webcam *video.Webcam, id int, handlers []*V4lHandler) *ControlList {
	ctll := &ControlList{
		webcam:   webcam,
		Id:       id,
		Handlers: make([]*V4lHandler, 0, len(handlers)),
	}
	for _, ctl := range handlers {
		ctll.AddHandler(ctl)
	}
	return ctll
}

func (ctll *ControlList) AddHandler(ctlh *V4lHandler) {
	var err error
	if ctlh == nil {
		log.Fatalln("AddControl control is nil")
	}
	ctlh.webcam = ctll.webcam
	ctlh.Info, err = ctlh.webcam.GetControlInfo(ctlh.Key)
	ctlh.Value = ctlh.webcam.GetControlValue(ctlh.Key)

	ctll.Handlers = append(ctll.Handlers, ctlh)
	if err != nil {
		log.Println("AddControl", err)
	}

	for _, ctl := range ctlh.Controls {
		http.Handle(ctl.url, ctlh)
	}

}

func (ctlh *ControlList) ResetControls() {
	for _, ctl := range ctlh.Handlers {
		ctl.webcam.SetValue(ctl.Key, ctl.Info.Default)
		log.Println("ResetControls", ctl.Key, ctl.Info.Default)
		ctl.Value = ctl.Info.Default
	}
}
