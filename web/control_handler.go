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
	wbch := &ControlHandler{
		Key:        key,
		Controls:   ctls,
		controlMap: make(map[string]*camera.Control),
		rt:         rt,
	}

	for _, ctl := range ctls {
		wbch.controlMap[ctl.Url] = ctl
	}
	return wbch
}

func (wbch *ControlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		info v4l.ControlInfo
		rt   = wbch.rt
	)

	defer func() {
		if err != nil {
			wbch.rt.template.ExecuteTemplate(w, "layout.response", 0)
		}
	}()

	control, ok := wbch.controlMap[r.RequestURI]
	if !ok {
		log.Println("Control not found", wbch.Key, r.RequestURI)
		return
	}

	camsrv, err := rt.parseSourceId(r)

	if camsrv.Config.Driver == IPWebcam {
		err = wbch.handleIpWebcam(camsrv, w, r)
		return
	}

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
		if !wbch.validValue {
			wbch.Value = webcam.GetControlValue(wbch.Key)
			wbch.validValue = true
		}
		info, err = webcam.GetControlInfo(wbch.Key)
		if err != nil {
			return
		}
		val := wbch.Value + info.Step*control.Multiplier

		// log.Println(cur, wbch.Value, val, info.Min, info.Max)
		if val >= info.Min && val <= info.Max {
			wbch.Value = val
			webcam.SetControlValue(wbch.Key, val)
		}
		// log.Println(control.Multiplier, info.Step)
		// log.Println(wbch.Value, val, info.Min, info.Max)
	}

	rt.template.ExecuteTemplate(w, "layout.response", wbch.Value)
}

func (wbch *ControlHandler) handleIpWebcam(camsrv *camera.Server, w http.ResponseWriter, r *http.Request) (err error) {
	log.Println("wbch.Key", wbch.Key, r.URL.RequestURI())

	var (
		cmd = camsrv.Config.Base
		val = -1
	// 	client = &http.Client{}
	// 	req    *http.Request
	// 	resp   *http.Response
	// 	buf    []byte
	)

	constrain := func(val, min, max int) int {
		if val < min {
			return min
		}
		if val > max {
			return max
		}
		return val
	}

	request := r.URL.RequestURI()
	switch request {
	case "/zoomin":
		val = constrain(val+10, 0, 100)
		cmd += fmt.Sprintf("/ptz?zoom=%d", val)
	case "/zoomout":
		val = constrain(val-10, 0, 100)
		cmd += fmt.Sprintf("/ptz?zoom=%d", val)
	default:
		err = fmt.Errorf("'%s' not yet implemented for %s", camsrv.Config.Driver, r.RequestURI)
		return
	}
	log.Println(cmd, val)
	err = fmt.Errorf("'%s' not yet implemented for %s", camsrv.Config.Driver, r.RequestURI)
	return
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
