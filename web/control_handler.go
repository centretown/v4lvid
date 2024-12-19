package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
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
	if err != nil {
		return
	}

	if camsrv.Config.Driver == IPWebcam {
		err = ctlh.handleIPWebcam(camsrv, control, w, r)
		return
	}

	if camsrv.Config.Driver != UVCVideo {
		err = fmt.Errorf("wrong driver '%s' for %s", camsrv.Config.Driver, r.RequestURI)
		log.Println(err)
		return
	}

	webcam, ok := camsrv.Source.(*camera.Webcam)
	if !ok {
		err = handleRemoteV4L(camsrv, w, r)
		return
	}

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

	rt.template.ExecuteTemplate(w, "layout.response", ctlh.Value)
}

func (ctlh *ControlHandler) handleIPWebcam(camsrv *camera.Server, control *camera.Control,
	w http.ResponseWriter, r *http.Request) (err error) {
	var (
		ipcam *camera.Ipcam
		ipwc  *camera.IpWebcam
		ok    bool
	)

	ipcam, ok = camsrv.Source.(*camera.Ipcam)
	if !ok {
		err = fmt.Errorf("not an ip camera")
		log.Println("ipwcHandler", err)
		return
	}

	if ipcam.State == nil {
		ipwc = camera.NewIpWebCam()
		ipcam.State = ipwc
		err = ipwc.Load(camsrv.Config.Base, ctlh.rt.Config.IPWCCommands)
		if err != nil {
			log.Println("LoadIpWebCamStatus", err)
			return
		}
	}

	ipwc, ok = ipcam.State.(*camera.IpWebcam)
	if !ok {
		err = fmt.Errorf("not an ipwebcam camera")
		log.Println("ipwcHandler", err)
		return
	}

	log.Println("handleIPWebcam", r.RequestURI, len(ipwc.Properties), ctlh.Key)
	// "/zoomin"
	// "/zoomout"
	// "/panleft"
	// "/panright"
	// "/tiltup"
	// "/tiltdown"
	// "/brightnessup"
	// "/brightnessdown"
	// "/contrastup"
	// "/contrastdown"
	// "/saturationup"
	// "/saturationdown"
	rt := ctlh.rt
	key, ok := rt.Config.IPWCControls[r.RequestURI]
	if !ok {
		err = fmt.Errorf("not an ipwebcam camera")
		log.Println("ipwcHandler", err)
		return
	}
	log.Println(key)
	command, ok := rt.Config.IPWCCommands[key]
	if !ok {
		err = fmt.Errorf("not an ipwebcam camera")
		log.Println("ipwcHandler", err)
		return
	}

	log.Println(command.InputType, command.Max, command.Step, command.Command)
	prop, ok := ipwc.Properties[key]
	if !ok {
		err = fmt.Errorf("not an ipwebcam camera property `%s`", key)
		log.Println("ipwcHandler", err)
		return
	}

	val, _ := strconv.Atoi(prop.Value)
	val += command.Step * int(control.Multiplier)
	if val < command.Min {
		val = command.Min
	} else if val > command.Max {
		val = command.Max
	}
	log.Println(key, val, prop.Value, control.Multiplier)
	rt.template.ExecuteTemplate(w, "layout.response", val)
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
