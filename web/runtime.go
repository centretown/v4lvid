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
	"v4lvid/homeasst"
	"v4lvid/socket"
)

type RunTime struct {
	WebcamUrl       string
	Config          *config.Config
	ActionsCamera   []*config.Action
	ActionsHome     []*config.Action
	ActionsChat     []*config.Action
	ActionMap       map[string]*config.Action
	ControlHandlers []*ControlHandler
	CameraServers   []*camera.Server
	CameraMap       map[string]*camera.Server
	Streamer        *Streamer
	mux             *http.ServeMux
	template        *template.Template
	Home            *homeasst.HomeRuntime
	WebSocket       *socket.Server
}

func Run(cfg *config.Config) (rt *RunTime) {
	rt = &RunTime{
		WebcamUrl: "/video0",
		// WebcamUrl: "http://192.168.10.7:9000/video0",
		Config:        cfg,
		ActionsCamera: cfg.ActionsCamera,
		ActionsHome:   cfg.ActionsHome,
		ActionsChat:   cfg.ActionsChat,
		ActionMap:     cfg.NewActionMap(),
		CameraMap:     make(map[string]*camera.Server),
		CameraServers: make([]*camera.Server, 0, len(cfg.Cameras)),
		mux:           &http.ServeMux{},
	}
	var (
		err        error
		httpServer = &http.Server{
			Handler:      rt.mux,
			Addr:         cfg.HttpUrl,
			ReadTimeout:  0,
			WriteTimeout: 0,
		}
	)

	const pattern = "www/*.html"
	rt.template, err = template.ParseGlob(pattern)
	if err != nil {
		log.Fatalln("ParseGlob", pattern, err)
	}

	rt.WebSocket = socket.NewServer(rt.template)
	rt.WebSocket.LoadMessages()
	rt.WebSocket.Run()

	rt.mux.HandleFunc("/events", rt.WebSocket.Events)
	rt.mux.HandleFunc("/webhook", rt.WebSocket.Webhook)

	rt.buildCameraServers()
	rt.ControlHandlers = CreateNexigoHandlers(rt)

	rt.serveCameras()
	rt.handleCameras()
	rt.Home, err = homeasst.NewHomeRuntime()
	if err == nil {
		rt.serveHomeData()
		rt.handleHomeData()
	}

	rt.handleFiles()

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

	rt.WebSocket.SaveMessages()

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)
	// wsServer.server.Shutdown(ctx)
	return
}

func (rt *RunTime) handleFiles() {
	rt.mux.Handle(rt.Streamer.Url, rt.Streamer)

	fs := http.FileServer(http.Dir("www/"))
	rt.mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	})

	rt.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		rt.template.ExecuteTemplate(w, "index.html", rt)
	})

	rt.mux.HandleFunc("/filesave", func(w http.ResponseWriter, r *http.Request) {
		err := config.Save(rt.Config, "config.json")
		if err != nil {
			log.Println("filesave", err)
			return
		}
	})

}

func (rt *RunTime) serveHomeData() (err error) {
	home := rt.Home
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

func (rt *RunTime) buildCameraServers() {

	for id, vcfg := range rt.Config.Cameras {
		cameraServer, err := newCameraServer(id, vcfg, rt.WebSocket)
		if err != nil {
			log.Println(err)
		}
		rt.CameraServers = append(rt.CameraServers, cameraServer)
		rt.CameraMap[cameraServer.Config.Path] = cameraServer
	}

	rt.Streamer = &Streamer{
		Server: rt.CameraServers[0],
		Url:    "/record",
		Icon:   "radio_button_checked",
		Socket: rt.WebSocket,
	}

}
