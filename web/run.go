package web

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"v4lvid/camera"
	"v4lvid/config"
	"v4lvid/ha"
	"v4lvid/socket"
)

type RunData struct {
	WebcamUrl       string
	Config          *config.Config
	Actions         []*config.Action
	ActionMap       map[string]*config.Action
	WebcamHandlers  []*WebcamHandler
	WebcamServers   []*camera.Server
	CameraMap       map[string]*camera.Server
	Streamer        *Streamer
	Temperature     float64
	TemperatureUnit string
	mux             *http.ServeMux
	template        *template.Template
	home            *ha.HomeData
	HomeActive      bool
	WebSocket       *socket.Server
}

func Run(cfg *config.Config) (data *RunData) {
	data = &RunData{
		WebcamUrl: "/video0",
		// WebcamUrl: "http://192.168.10.7:9000/video0",
		Config:        cfg,
		Actions:       cfg.Actions,
		ActionMap:     cfg.NewActionMap(),
		CameraMap:     make(map[string]*camera.Server),
		WebcamServers: make([]*camera.Server, 0, len(cfg.Cameras)),
		mux:           &http.ServeMux{},
	}
	var (
		err        error
		httpServer = &http.Server{
			Handler:      data.mux,
			Addr:         cfg.HttpUrl,
			ReadTimeout:  0,
			WriteTimeout: 0,
		}
	)

	const pattern = "www/*.html"
	data.template, err = template.ParseGlob(pattern)
	if err != nil {
		log.Fatalln("ParseGlob", pattern, err)
	}

	data.WebSocket = socket.NewServer(data.template)
	data.WebSocket.LoadMessages()
	data.WebSocket.Run()

	data.buildCameraServers()

	data.mux.HandleFunc("/events", data.WebSocket.Events)
	data.mux.HandleFunc("/webhook", data.WebSocket.Webhook)

	data.WebcamHandlers = CreateNexigoHandlers(cfg, data.WebcamServers,
		data.template)

	serveCameras(data)
	handleCameras(data)

	data.home, err = ha.NewHomeData()
	if err == nil {
		serveHomeData(data)
		handleHomeData(data)
		data.HomeActive = true
	}

	handleFiles(data)

	httpErr := make(chan error, 1)
	go func() {
		httpErr <- httpServer.ListenAndServe()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-httpErr:
		log.Printf("failed to serve http: %v", err)
	// case err := <-wsErr:
	// 	log.Printf("failed to serve websockets: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	data.WebSocket.SaveMessages()

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)
	// wsServer.server.Shutdown(ctx)
	return
}

func handleFiles(data *RunData) {
	data.mux.Handle(data.Streamer.Url, data.Streamer)

	fs := http.FileServer(http.Dir("www/"))
	data.mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	})

	data.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		data.template.ExecuteTemplate(w, "index.html", data)
	})

	data.mux.HandleFunc("/filesave", func(w http.ResponseWriter, r *http.Request) {
		err := config.Save(data.Config, "config.json")
		if err != nil {
			log.Println("filesave", err)
			return
		}
	})

}

func serveHomeData(data *RunData) (err error) {
	home := data.home
	var ok bool
	ok, err = home.Authorize()
	if err != nil {
		log.Println("authorize", err)
		return
	}
	if !ok {
		err = fmt.Errorf("not authorized")
		log.Println(err)
		return
	}

	log.Println("Authorized HA")

	err = home.BuildEntities()
	if err != nil {
		log.Println("BuildEntities", err)
		return

	}
	log.Println("Build Entities")

	go home.Monitor()

	if home.Monitoring {
		log.Println("Monitor Entity States")
	}
	log.Println("Monitoring")

	return
}

func newCameraServer(id int, vcfg *camera.VideoConfig,
	indicator camera.StreamIndicator) (cameraServer *camera.Server, err error) {

	var source camera.VideoSource
	switch vcfg.CameraType {
	case camera.LOCAL_CAMERA:
		source = camera.NewWebcam(vcfg.Path)
	case camera.REMOTE_CAMERA:
		source = camera.NewIpcam(vcfg.Path)
	default:
		return
	}
	cameraServer = camera.NewVideoServer(id, source, vcfg, indicator)
	err = cameraServer.Open()
	return
}

func (data *RunData) buildCameraServers() {

	for id, vcfg := range data.Config.Cameras {
		cameraServer, err := newCameraServer(id, vcfg, data.WebSocket)
		if err != nil {
			log.Println(err)
		}
		data.WebcamServers = append(data.WebcamServers, cameraServer)
		data.CameraMap[cameraServer.Config.Path] = cameraServer
	}

	data.Streamer = &Streamer{
		Server: data.WebcamServers[0],
		Url:    "/record",
		Icon:   "radio_button_checked",
		Socket: data.WebSocket,
	}

}
