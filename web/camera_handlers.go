package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"v4lvid/camera"
	"v4lvid/config"
)

type CameraListData struct {
	Action  *config.Action
	Cameras []*camera.Server
}

type CameraData struct {
	Action         *config.Action
	WebcamHandlers []*ControlHandler
}

type AddCamera struct {
	Action *config.Action
}

func (rt *RunTime) handleCameras() {
	rt.mux.HandleFunc("/camera", rt.controlCameraHandler())
	rt.mux.HandleFunc("/camera_add", rt.addCameraHandler())
	rt.mux.HandleFunc("/camera_post", rt.postCameraHandler())
	rt.mux.HandleFunc("/camera_list", rt.listCameraHandler())
	rt.mux.HandleFunc("/camera_connect", rt.connectCameraHandler())
	rt.mux.HandleFunc("/camera_primary", rt.setPrimaryCamera())
	rt.mux.HandleFunc("/ipwc/", rt.ipwcCameraHandler())

}

func (rt *RunTime) ipwcCameraHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cam, err := rt.parseSourceId(r)
		if err != nil {
			log.Println("ipcwCameraHandler", err)
			return
		}

		i := strings.LastIndex(r.RequestURI, "/")
		if i == -1 {
			log.Println("LastIndex /")
			return
		}

		fld := r.RequestURI[i+1:]
		log.Println(fld, r.FormValue(fld), r.RequestURI)
		log.Println("ipcwCameraHandler", cam.Config.Path)
	}
}

func (rt *RunTime) parseCameraPath(r *http.Request) (cam *camera.Server,
	path string, err error) {

	err = r.ParseForm()
	if err != nil {
		err = fmt.Errorf("parse form: %v", err)
		return
	}

	path = r.FormValue("path")
	cam, ok := rt.CameraMap[path]
	if !ok {
		err = fmt.Errorf("path not found: %s", path)
		return
	}
	return
}

func wrapStatus(id, msg string) []byte {
	return []byte(fmt.Sprintf(`<div id="%s" class="status">%s</div>`, id, msg))
}

func (rt *RunTime) controlCameraHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cam, err := rt.parseSourceId(r)
		if err != nil {
			log.Println("controlCameraHandler", err)
			return
		}

		log.Println("controlCameraHandler", cam.Config.Driver)
		w.Header().Add("Cache-Control", "no-cache")

		if cam.Config.Driver == "uvcvideo" {
			err = rt.template.Lookup("layout.controls").Execute(w,
				&CameraData{
					Action:         rt.ActionMap["camera"],
					WebcamHandlers: rt.ControlHandlers})
			if err != nil {
				log.Fatal("/camera", err)
			}
		} else if cam.Config.Driver == "ipwebcam" {
			rt.ipwcHandler(cam, w, r)
		}
	}
}

func (rt *RunTime) ipwcHandler(cam *camera.Server, w http.ResponseWriter, r *http.Request) {
	var (
		ipcam *camera.Ipcam
		ipwc  *camera.IPWebcam
		err   error
		ok    bool
	)

	ipcam, ok = cam.Source.(*camera.Ipcam)
	if !ok {
		log.Println("ipwcHandler", "not an ip camera")
		return
	}

	if ipcam.State != nil {
		ipwc, ok = ipcam.State.(*camera.IPWebcam)
		if !ok {
			log.Println("ipwcHandler", "not an ipwebcam camera")
			return
		}
	} else {
		ipwc = camera.NewIpWebCam()
		ipcam.State = ipwc
	}

	err = ipwc.Load(cam.Config.Base, rt.Config.IPWCCommands)
	if err != nil {
		log.Println("LoadIpWebCamStatus", err)
		return
	}

	err = rt.template.Lookup("layout.ipwebcam").Execute(w, &config.IPWCCameraData{
		Action:   rt.ActionMap["camera"],
		IPWebcam: ipwc,
	})
	if err != nil {
		log.Fatal("Lookup", err)
	}

}

func (rt *RunTime) addCameraHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		err := rt.template.Lookup("layout.camera.add").Execute(w,
			&AddCamera{
				Action: rt.ActionMap["camera_add"],
			})

		if err != nil {
			log.Fatal("/camera_add", err)
		}
	}
}

func (rt *RunTime) setPrimaryCamera() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const statusID = "camera_list_status"
		const sourceID = "source"
		wrapSource := func(id, src string) []byte {
			return []byte(fmt.Sprintf(`<img id="%s" src="%s">`, id, src))
		}

		cam, path, err := rt.parseCameraPath(r)
		if err != nil {
			msg := fmt.Sprintf("Error occured.<br>  %v", err)
			w.Write(wrapStatus(statusID, msg))
			return
		}

		if !cam.Source.IsOpened() {
			msg := fmt.Sprintf("%s as %s is not connected", path, cam.Url())
			w.Write(wrapStatus(statusID, msg))
			return
		}

		msg := fmt.Sprintf("%s is connected as %s", path, cam.Url())
		w.Write(wrapStatus(statusID, msg))
		w.Write(wrapSource(sourceID, cam.Url()))

		// `<img id="source" src="{{.WebcamUrl}}">`

	}
}

func (rt *RunTime) connectCameraHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cam, path, err := rt.parseCameraPath(r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error occured.<br>  %v", err)))
			return
		}

		if cam.Source.IsOpened() {
			w.Write([]byte(fmt.Sprintf("%s is already connected as %s",
				path, cam.Url())))
			return
		}

		err = cam.Open()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Problem connecting to %s.<br>  %v",
				path, err)))
			return
		}

		go cam.Serve()
		msg := fmt.Sprintf("Connected to path %s as %s", path, cam.Url())
		w.Write([]byte(msg))
		log.Println(msg)
	}
}

func (rt *RunTime) listCameraHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		camData := &CameraListData{
			Action:  rt.ActionMap["camera_list"],
			Cameras: rt.CameraServers,
		}
		err := rt.template.Lookup("layout.camera.list").Execute(w, camData)
		if err != nil {
			log.Fatal("/camera_list", err)
		}
	}
}

func (rt *RunTime) postCameraHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("camera_post")
		w.Header().Add("Cache-Control", "no-cache")
		err := r.ParseForm()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Camera error parsing form. %v", err)))
			return
		}

		path := fmt.Sprintf("%s.%s:%s%s",
			r.FormValue("camera_net"),
			r.FormValue("camera_suffix"),
			r.FormValue("camera_port"),
			r.FormValue("camera_url"))

		ws, ok := rt.CameraMap[path]
		if ok {
			w.Write([]byte(fmt.Sprintf("Camera already on file. %s", path)))
			err = ws.Open()
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Camera error connecting. %v", err)))
				return
			}

			go ws.Serve()
			w.Write([]byte(fmt.Sprintf("Camera error parsing form. %v", err)))
			return
		}

		width, _ := strconv.Atoi(r.FormValue("camera_width"))
		height, _ := strconv.Atoi(r.FormValue("camera_height"))
		fps, _ := strconv.Atoi(r.FormValue("camera_fps"))
		vc := &camera.VideoConfig{
			Path:       path,
			CameraType: camera.REMOTE_CAMERA,
			Codec:      r.FormValue("camera_codec"),
			Width:      width,
			Height:     height,
			FPS:        uint32(fps),
		}

		id := rt.Config.AddCamera(vc)
		ws, err = newCameraServer(id, vc, rt.WebSocket)
		// add even if error
		rt.CameraMap[path] = ws
		rt.CameraServers = append(rt.CameraServers, ws)
		if err != nil {
			msg := fmt.Sprintf("Camera Added %s.<br>The following error was reported:<br>%v", path, err)
			w.Write([]byte(msg))
			return
		}

		rt.serveCamera(ws)
		w.Write([]byte(fmt.Sprintf("Connected to %s as %s", path, ws.Url())))
	}
}

func (rt *RunTime) serveCamera(camServer *camera.Server) {
	rt.mux.Handle(camServer.Url(), camServer.Stream())
	go camServer.Serve()
}

func (rt *RunTime) serveCameras() {
	for _, camServer := range rt.CameraServers {
		rt.serveCamera(camServer)
	}

	for _, handler := range rt.ControlHandlers {
		for _, ctl := range handler.Controls {
			rt.mux.Handle(ctl.Url, handler)
		}
	}

	rt.mux.HandleFunc("/resetcontrols",
		func(w http.ResponseWriter, r *http.Request) {
			camsrv, err := rt.parseSourceId(r)
			if err != nil {
				return
			}

			switch camsrv.Config.Driver {
			case UVCVideo:
				break
			case IPWebcam:
				log.Printf("net yet implemented '%s' for %s", camsrv.Config.Driver, r.RequestURI)
				return
			default:
				log.Printf("wrong driver '%s' for %s", camsrv.Config.Driver, r.RequestURI)
				return
			}

			_, ok := camsrv.Source.(*camera.Ipcam)
			if ok && camsrv.Config.Driver == UVCVideo {
				err = handleRemoteV4L(camsrv, w, r)
				if err != nil {
					log.Println(err)
				}
				return
			}

			webcam, ok := camsrv.Source.(*camera.Webcam)
			if ok {
				rt.ResetControls(webcam)
			}
		})
}

func (rt *RunTime) ResetControls(webcam *camera.Webcam) {
	for _, ctlh := range rt.ControlHandlers {
		info, err := webcam.GetControlInfo(ctlh.Key)
		if err != nil {
			log.Println("ResetControls", ctlh.Key, err)
			continue
		}
		ctlh.Value = info.Default
		webcam.SetControlValue(ctlh.Key, ctlh.Value)
		log.Println("ResetControls", ctlh.Key, ctlh.Value)
	}
}

func (rt *RunTime) handleStreamer() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		camsrv, err := rt.parseSourceId(r)
		if err != nil {
			return
		}

		if !camsrv.Recording {
			log.Printf("recording...")
			camsrv.RecordCmd(300)
		} else {
			log.Printf("stop recording...")
			camsrv.StopRecordCmd()
		}
	}
}
