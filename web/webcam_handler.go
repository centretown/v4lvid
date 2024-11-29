package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"v4lvid/camera"
)

var _ http.Handler = (*WebcamHandler)(nil)

type WebcamHandler struct {
	Key string

	// Info       v4l.ControlInfo
	// Value      int32
	Controls   []*camera.Control
	Cameras    []*camera.Server
	controlMap map[string]*camera.Control
	tmpl       *template.Template
}

func NewWebcamHandler(key string, ctls []*camera.Control,
	cameras []*camera.Server,
	tmpl *template.Template) *WebcamHandler {
	handler := &WebcamHandler{
		Key:        key,
		Controls:   ctls,
		Cameras:    cameras,
		controlMap: make(map[string]*camera.Control),
		tmpl:       tmpl,
	}

	for _, ctl := range ctls {
		handler.controlMap[ctl.Url] = ctl
	}
	return handler
}

func (handler *WebcamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	control, ok := handler.controlMap[r.RequestURI]
	if !ok {
		log.Println("Control not found", handler.Key, r.RequestURI)
		return
	}

	id, err := parseSource(r)
	if err != nil {
		log.Println("unable to parse request", r.RequestURI, err)
		return
	}

	if id >= len(handler.Cameras) {
		log.Printf("Camera id = %d in source id out of range limit\n", id)
		return
	}

	camsrv := handler.Cameras[id]
	if camsrv.Config.Driver != "uvcvideo" {
		log.Println("not uvcvideo driver", r.RequestURI)
		return
	}

	// remote?
	_, ok = camsrv.Source.(*camera.Ipcam)
	if ok {
		handler.handleRemote(camsrv, w, r)
		return
	}

	var value int32
	webcam, ok := camsrv.Source.(*camera.Webcam)
	if ok {
		ctl, _ := webcam.GetControlInfo(handler.Key)
		value = webcam.GetControlValue(handler.Key)
		newValue := value + ctl.Step*control.Multiplier
		if newValue >= ctl.Min && newValue <= ctl.Max {
			value = newValue
			webcam.SetValue(handler.Key, newValue)
		}
	}

	handler.tmpl.ExecuteTemplate(w, "layout.response", value)
}

func (handler *WebcamHandler) handleRemote(camsrv *camera.Server, w http.ResponseWriter, r *http.Request) {
	var (
		url    = camsrv.Config.Base + r.RequestURI
		client = &http.Client{}
		err    error
		req    *http.Request
		resp   *http.Response
		buf    []byte
	)

	defer func() {
		if err != nil {
			handler.tmpl.ExecuteTemplate(w, "layout.response", 0)
		}
	}()

	req, err = http.NewRequest(r.Method, url, nil)
	if err != nil {
		log.Println("handleRemote NewRequest", url, err)
		return
	}

	resp, err = client.Do(req)
	if err != nil {
		log.Println("client.Do", err)
		return
	}

	buf, err = io.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		log.Println("ReadAll", err)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		log.Println("Write", err)
		return
	}
}

const prefix = "/video"

func parseSource(r *http.Request) (id int, err error) {
	err = r.ParseForm()
	if err != nil {
		log.Println("ParseForm", err)
		return
	}

	source := r.FormValue("source")
	last := strings.LastIndex(source, prefix)
	if last == -1 {
		err = fmt.Errorf("invalid source %s", source)
		log.Println(err)
		return
	}

	source = source[last+len(prefix):]
	count, err := fmt.Sscan(source, &id)
	if err != nil {
		log.Println("scan", err)
		return
	}

	if count != 1 {
		err = fmt.Errorf("not a valid source, count = %d '%s'", count, source)
		log.Println(err)
		return
	}
	return
}
