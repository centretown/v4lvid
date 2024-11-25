package web

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"v4lvid/camera"

	"github.com/korandiz/v4l"
)

var _ http.Handler = (*WebcamHandler)(nil)

type WebcamHandler struct {
	webcam camera.VideoSource
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
	var (
		err    error
		source string
	)
	err = r.ParseForm()
	if err != nil {
		log.Println("src: ", err)
		return
	}

	source = r.FormValue("source")
	last := strings.LastIndex(source, "/")
	if last >= 0 {
		source = source[last+1:]
	}

	log.Println("src:", source)
	control, ok := handler.controlMap[r.RequestURI]
	if !ok {
		log.Println("RequestURI not found", handler.Key, r.RequestURI)
		return
	}

	webcam, ok := handler.webcam.(*camera.Webcam)
	if ok {
		log.Println("Handle", handler.Key, r.RequestURI)
		newValue := handler.Value + handler.Info.Step*control.Multiplier
		if newValue >= handler.Info.Min && newValue <= handler.Info.Max {
			handler.Value = newValue
			webcam.SetValue(handler.Key, newValue)
		}
		err = handler.tmpl.ExecuteTemplate(w, "layout.response", handler.Value)
		if err != nil {
			log.Println("ControlHandler", err)
		}
	} else {
		//TODO
	}
}
