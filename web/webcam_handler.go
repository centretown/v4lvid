package web

import (
	"html/template"
	"log"
	"net/http"
	"v4lvid/camera"

	"github.com/korandiz/v4l"
)

var _ http.Handler = (*WebcamHandler)(nil)

type WebcamHandler struct {
	webcam *camera.Webcam
	Key    string

	Info       v4l.ControlInfo
	Value      int32
	Controls   []*camera.Control
	controlMap map[string]*camera.Control
	tmpl       *template.Template
}

func NewWebcamHandler(key string, ctls []*camera.Control, tmpl *template.Template) *WebcamHandler {
	handler := &WebcamHandler{
		Key:        key,
		Controls:   ctls,
		controlMap: make(map[string]*camera.Control),
		tmpl:       tmpl,
	}

	for _, ctl := range ctls {
		handler.controlMap[ctl.Url] = ctl
		// if ctl.Items != nil {
		// 	for _, child := range ctl.Items {
		// 		handler.Map[child.Url] = child
		// 	}
		// }
	}
	return handler
}

func (handler *WebcamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	control, ok := handler.controlMap[r.RequestURI]
	if !ok {
		log.Println("RequestURI not found", handler.Key, r.RequestURI)
		return
	}

	log.Println("Handle", handler.Key, r.RequestURI)
	newValue := handler.Value + handler.Info.Step*control.Multiplier
	if newValue >= handler.Info.Min && newValue <= handler.Info.Max {
		handler.Value = newValue
		handler.webcam.SetValue(handler.Key, newValue)
	}
	err := handler.tmpl.ExecuteTemplate(w, "layout.response", handler.Value)
	if err != nil {
		log.Println("ControlHandler", err)
	}
}
