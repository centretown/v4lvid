package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"v4lvid/camera"

	"github.com/korandiz/v4l"
)

var _ http.Handler = (*ControlHandler)(nil)

const (
	Prefix   = "/video"
	UVCVideo = "uvcvideo"
	IPWebcam = "ipwebcam"
)

type ControlHandler struct {
	Key        string
	Value      int32
	validValue bool
	Controls   []*camera.Control
	controlMap map[string]*camera.Control
	rt         *RunTime
}

func NewControlHandler(key string, ctls []*camera.Control,
	rt *RunTime) *ControlHandler {
	ctlh := &ControlHandler{
		Key:        key,
		Controls:   ctls,
		controlMap: make(map[string]*camera.Control),
		rt:         rt,
	}

	for _, ctl := range ctls {
		ctlh.controlMap[ctl.Url] = ctl
	}
	return ctlh
}

func (ctlh *ControlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		info v4l.ControlInfo
		rt   = ctlh.rt
	)

	defer func() {
		if err != nil {
			ctlh.rt.template.ExecuteTemplate(w, "layout.response", 0)
		}
	}()

	control, ok := ctlh.controlMap[r.RequestURI]
	if !ok {
		log.Println("Control not found", ctlh.Key, r.RequestURI)
		return
	}

	camsrv, err := rt.parseSourceId(r)

	if camsrv.Config.Driver != UVCVideo {
		err = fmt.Errorf("wrong driver '%s' for %s", camsrv.Config.Driver, r.RequestURI)
		log.Println(err)
		return
	}

	_, ok = camsrv.Source.(*camera.Ipcam)
	if ok {
		err = handleRemoteV4L(camsrv, w, r)
		return
	}

	webcam, ok := camsrv.Source.(*camera.Webcam)
	if ok {
		if !ctlh.validValue {
			ctlh.Value = webcam.GetControlValue(ctlh.Key)
			ctlh.validValue = true
		}
		info, err = webcam.GetControlInfo(ctlh.Key)
		if err != nil {
			return
		}
		// some steps don't take on the device (tilt) so we assume we have
		// the correct value even when we don't to skip the "holes"
		val := ctlh.Value + info.Step*control.Multiplier
		if val >= info.Min && val <= info.Max {
			ctlh.Value = val
			webcam.SetControlValue(ctlh.Key, val)
		}
	}

	rt.template.ExecuteTemplate(w, "layout.response", ctlh.Value)
}

func handleRemoteV4L(camsrv *camera.Server, w http.ResponseWriter, r *http.Request) (err error) {
	var (
		url    = camsrv.Config.Base + r.RequestURI
		client = &http.Client{}
		req    *http.Request
		resp   *http.Response
		buf    []byte
	)

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
	return
}

func (rt *RunTime) parseSourceId(r *http.Request) (camsrv *camera.Server, err error) {
	var id int
	err = r.ParseForm()
	if err != nil {
		log.Println("ParseForm", err)
		return
	}

	source := r.FormValue("source")
	last := strings.LastIndex(source, Prefix)
	if last == -1 {
		err = fmt.Errorf("prefix '%s' not found in source %s", Prefix, source)
		log.Println(err)
		return
	}
	offset := last + len(Prefix)
	if offset >= len(source) {
		err = fmt.Errorf("source too short %s", source)
		log.Println(err)
		return
	}

	count, err := fmt.Sscan(source[offset:], &id)
	if err != nil {
		log.Println("scan", err)
		return
	}
	if count != 1 {
		err = fmt.Errorf("not a valid source, count = %d '%s'", count, source)
		log.Println(err)
		return
	}

	if id >= len(rt.CameraServers) {
		log.Printf("Camera id = %d in source id out of range limit (%d)\n",
			id, len(rt.CameraServers))
		return
	}

	camsrv = rt.CameraServers[id]
	return
}
