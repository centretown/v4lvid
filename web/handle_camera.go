package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"v4lvid/camera"
	"v4lvid/config"
)

type CameraListData struct {
	Action  *config.Action
	Cameras []*camera.Server
}

type CameraData struct {
	Action         *config.Action
	WebcamHandlers []*WebcamHandler
}

type AddCamera struct {
	Action *config.Action
}

func handleCameras(data *RunData) {
	data.mux.HandleFunc("/camera",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Cache-Control", "no-cache")
			err := data.template.Lookup("layout.controls").Execute(w,
				&CameraData{
					Action:         data.ActionMap["camera"],
					WebcamHandlers: data.WebcamHandlers})
			if err != nil {
				log.Fatal("/camera", err)
			}
		})
	data.mux.HandleFunc("/camera_add",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Cache-Control", "no-cache")
			err := data.template.Lookup("layout.camera.add").Execute(w,
				&AddCamera{
					Action: data.ActionMap["camera_add"],
				})

			if err != nil {
				log.Fatal("/camera_add", err)
			}
		})
	data.mux.HandleFunc("/camera_post", addCameraHandler(data))
	data.mux.HandleFunc("/camera_list", listCameraHandler(data))
	data.mux.HandleFunc("/camera_connect", connectCameraHandler(data))

}

func connectCameraHandler(data *RunData) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Camera error parsing form. %v", err)))
			return
		}

		path := r.FormValue("path")
		w.Write([]byte(fmt.Sprintf("Camera Path: %s", path)))
		log.Println(path)
	}
}

func listCameraHandler(data *RunData) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		camData := &CameraListData{
			Action:  data.ActionMap["camera_list"],
			Cameras: data.WebcamServers,
		}
		err := data.template.Lookup("layout.camera.list").Execute(w, camData)
		if err != nil {
			log.Fatal("/camera_list", err)
		}
	}
}

func addCameraHandler(data *RunData) func(w http.ResponseWriter, r *http.Request) {
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

		ws, ok := data.CameraMap[path]
		if ok {
			w.Write([]byte(fmt.Sprintf("Camera already on file. %s", path)))
			err = ws.Open()
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Camera error parsing form. %v", err)))
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
			CameraType: camera.IP_CAMERA,
			Codec:      r.FormValue("camera_codec"),
			Width:      width,
			Height:     height,
			FPS:        uint32(fps),
		}

		id := data.Config.AddCamera(vc)
		ws, err = NewCameraServer(id, vc, data.WebSocket)
		data.CameraMap[path] = ws
		data.WebcamServers = append(data.WebcamServers, ws)
		if err != nil {
			msg := fmt.Sprintf("Camera Added %s.<br>The following error was reported:<br>%v", path, err)
			w.Write([]byte(msg))
			return
		}

		serveCamera(data, ws)
		w.Write([]byte(fmt.Sprintf("Connected to %s as %s", path, ws.Path())))
	}
}

func serveCamera(data *RunData, camServer *camera.Server) {
	log.Println("serveCamera", camServer.Path())
	data.mux.Handle(camServer.Path(), camServer.Stream())
	source := camServer.Source
	webcam, isWebcam := source.(*camera.Webcam)
	if isWebcam {
		ctll := NewControlList(data.mux, webcam, 0, data.WebcamHandlers)
		data.mux.HandleFunc("/resetcontrols",
			func(w http.ResponseWriter, r *http.Request) {
				ctll.ResetControls()
			})
	}

	go camServer.Serve()
	log.Printf("Serving %s\n", camServer.Path())
}

func serveCameras(data *RunData) {
	for _, camServer := range data.WebcamServers {
		serveCamera(data, camServer)
	}
}

func NewCameraServer(id int, vcfg *camera.VideoConfig,
	indicator camera.StreamIndicator) (cameraServer *camera.Server, err error) {

	var source camera.VideoSource
	switch vcfg.CameraType {
	case camera.V4L_CAMERA:
		source = camera.NewWebcam(vcfg.Path)
	case camera.IP_CAMERA:
		source = camera.NewIpcam(vcfg.Path)
	default:
		return
	}
	cameraServer = camera.NewVideoServer(id, source, vcfg, indicator)
	err = cameraServer.Open()
	return
}

func NewCameraServers(cfg *config.Config, indicator camera.StreamIndicator) (cameraServers []*camera.Server) {
	cameraServers = make([]*camera.Server, 0, len(cfg.Cameras))
	var (
		cameraServer *camera.Server
		err          error
	)

	for id, vcfg := range cfg.Cameras {
		cameraServer, err = NewCameraServer(id, vcfg, indicator)
		if err != nil {
			log.Println(err)
		}
		cameraServers = append(cameraServers, cameraServer)
	}
	return
}
