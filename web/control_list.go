package web

import (
	"log"
	"net/http"
	"v4lvid/camera"
)

type ControlList struct {
	webcam camera.VideoSource
	// webcam   *camera.Webcam
	Id       int
	Handlers []*WebcamHandler
}

func NewControlList(mux *http.ServeMux, webcam camera.VideoSource, id int, handlers []*WebcamHandler) *ControlList {
	ctll := &ControlList{
		webcam:   webcam,
		Id:       id,
		Handlers: make([]*WebcamHandler, 0, len(handlers)),
	}
	for _, ctl := range handlers {
		ctll.AddHandler(mux, ctl)
	}
	return ctll
}

func (ctll *ControlList) AddHandler(mux *http.ServeMux, ctlh *WebcamHandler) {
	var err error
	if ctlh == nil {
		log.Fatalln("AddControl control is nil")
	}

	// ctlh.webcam = ctll.webcam
	// webcam, ok := ctlh.webcam.(*camera.Webcam)
	// if ok {
	// 	ctlh.Info, err = webcam.GetControlInfo(ctlh.Key)
	// 	ctlh.Value = webcam.GetControlValue(ctlh.Key)
	// }

	ctll.Handlers = append(ctll.Handlers, ctlh)
	if err != nil {
		log.Println("AddControl", err)
	}

	for _, ctl := range ctlh.Controls {
		mux.Handle(ctl.Url, ctlh)
	}

}

func (ctlh *ControlList) ResetControls() {
	webcam, ok := ctlh.webcam.(*camera.Webcam)
	if ok {
		for _, ctl := range ctlh.Handlers {
			info, err := webcam.GetControlInfo(ctl.Key)
			if err != nil {
				log.Println("ResetControls", ctl.Key, err)
				continue
			}
			webcam.SetValue(ctl.Key, info.Default)
			log.Println("ResetControls", ctl.Key, info.Default)
		}
	} else {
		//TODO
	}
}
