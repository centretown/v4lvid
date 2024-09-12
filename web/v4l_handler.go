package web

import (
	"html/template"
	"log"
	"net/http"
	"v4lvid/camera"

	"github.com/korandiz/v4l"
)

var _ http.Handler = (*V4lHandler)(nil)

// var _ MenuItem = (*V4lHandler)(nil)

type V4lHandler struct {
	webcam *camera.Webcam
	Key    string

	Info     v4l.ControlInfo
	Value    int32
	Controls []*V4lControl
	Map      map[string]*V4lControl
	tmpl     *template.Template
}

func NewV4lHandler(key string, ctls []*V4lControl, tmpl *template.Template) *V4lHandler {
	handler := &V4lHandler{
		Key:      key,
		Controls: ctls,
		Map:      make(map[string]*V4lControl),
		tmpl:     tmpl,
	}

	for _, ctl := range ctls {
		handler.Map[ctl.url] = ctl
		if ctl.controls != nil {
			for _, child := range ctl.controls {
				handler.Map[child.url] = child
			}
		}
	}
	return handler
}

func (handler *V4lHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	control, ok := handler.Map[r.RequestURI]
	if !ok {
		log.Println("RequestURI not found", handler.Key, r.RequestURI)
		return
	}

	log.Println("Handle", handler.Key, r.RequestURI)

	if len(control.controls) > 0 {
		handler.ServeMenu(control, w, r)
		return
	}

	newValue := handler.Value + handler.Info.Step*control.multiplier
	if newValue >= handler.Info.Min && newValue <= handler.Info.Max {
		handler.Value = newValue
		handler.webcam.SetValue(handler.Key, newValue)
	}
	err := handler.tmpl.ExecuteTemplate(w, "layout.response", handler.Value)
	if err != nil {
		log.Println("ControlHandler", err)
	}
}

func (handler *V4lHandler) ServeMenu(control *V4lControl, w http.ResponseWriter, r *http.Request) {
	err := handler.tmpl.ExecuteTemplate(w, "layout.item", control)
	if err != nil {
		log.Println("ControlHandler", err)
	}
}
