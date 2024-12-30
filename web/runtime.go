package web

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
	"v4lvid/config"
	"v4lvid/homeasst"
	"v4lvid/socket"

	"github.com/centretown/avcam"
)

type RunTime struct {
	WebcamUrl       string
	Config          *config.Config
	ActionsCamera   []*config.Action
	ActionsHome     []*config.Action
	ActionsChat     []*config.Action
	ActionMap       map[string]*config.Action
	ControlHandlers []*ControlHandler
	CameraServers   []*avcam.Server
	CameraMap       map[string]*avcam.Server
	mux             *http.ServeMux
	template        *template.Template
	Home            *homeasst.HomeRuntime
	WebSocket       *socket.Server
	AudioMgr        *avcam.AudioMgr
}

func Run(cfg *config.Config) (rt *RunTime) {
	rt = &RunTime{
		WebcamUrl:     "/video0",
		Config:        cfg,
		ActionsCamera: cfg.ActionsCamera,
		ActionsHome:   cfg.ActionsHome,
		ActionsChat:   cfg.ActionsChat,
		ActionMap:     cfg.NewActionMap(),
		CameraMap:     make(map[string]*avcam.Server),
		CameraServers: make([]*avcam.Server, 0, len(cfg.Cameras)),
		mux:           &http.ServeMux{},
		AudioMgr:      avcam.NewAudio(),
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
	rt.ControlHandlers = rt.CreateV4LHandlers()

	rt.serveCameras()
	rt.handleCameras()
	rt.Home, err = homeasst.NewHomeRuntime()
	if err == nil {
		rt.serveHomeData()
		rt.homeHandler()
	}

	rt.handleAudio()

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
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	rt.WebSocket.SaveMessages()

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)
	return
}

func (rt *RunTime) handleFiles() {
	rt.mux.HandleFunc("/record", rt.handleStreamer())

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

func newCameraServer(id int, vcfg *avcam.VideoConfig, audio avcam.AudioSource,
	indicator avcam.StreamListener) (cameraServer *avcam.Server, err error) {

	var source avcam.VideoSource
	switch vcfg.CameraType {
	case avcam.LOCAL_CAMERA:
		source = avcam.NewWebcam(vcfg.Path)
	case avcam.REMOTE_CAMERA:
		source = avcam.NewIpcam(vcfg.Path)
	default:
		return
	}
	cameraServer = avcam.NewVideoServer(id, source, vcfg, audio, indicator)
	err = cameraServer.Open()
	return
}

func (rt *RunTime) buildCameraServers() {

	for id, vcfg := range rt.Config.Cameras {
		cameraServer, err := newCameraServer(id, vcfg, rt.AudioMgr, rt.WebSocket)
		if err != nil {
			log.Println(err)
		}
		rt.CameraServers = append(rt.CameraServers, cameraServer)
		rt.CameraMap[cameraServer.Config.Path] = cameraServer
	}
}

func (rt *RunTime) CreateV4LHandlers() (handlers []*ControlHandler) {

	handlers = make([]*ControlHandler, 0)
	driver, ok := rt.Config.Drivers[UVCVideo]
	if !ok {
		log.Println("Driver not found", UVCVideo)
		return
	}

	for _, d := range driver {
		handlers = append(handlers,
			NewControlHandler(d.Key, d.Controls, rt))
	}
	return
}

func (rt *RunTime) parseSourceId(r *http.Request) (camsrv *avcam.Server, err error) {
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
