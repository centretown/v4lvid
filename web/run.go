package web

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
	"v4lvid/camera"
	"v4lvid/config"
	"v4lvid/ha"
	"v4lvid/sockserve"
)

type RunData struct {
	WebcamUrl string
	cfg       *config.Config
	Actions   []*config.Action
	ActionMap map[string]*config.Action
	// WebcamServers  []*camera.Server
	WebcamHandlers []*WebcamHandler
	Recorder       *RecordingHandler
	mux            *http.ServeMux
	template       *template.Template
	home           *ha.HomeData
	Sock           *sockserve.SockServer
}

func Run(cfg *config.Config) (data *RunData) {
	webcamServers := cfg.NewCameraServers()
	data = &RunData{
		WebcamUrl: "http://192.168.10.7:9000/0/",
		Actions:   cfg.Actions,
		ActionMap: cfg.NewActionMap(),
		Recorder: &RecordingHandler{
			Server: webcamServers[0],
			Url:    "/record",
			Icon:   "radio_button_checked",
		},
		mux:  &http.ServeMux{},
		home: ha.NewHomeData(),
		cfg:  cfg,
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
	data.Sock = sockserve.NewSockServer(data.template)
	data.Sock.LoadMessages()
	data.Sock.Run()

	data.mux.HandleFunc("/events", data.Sock.Events)
	data.mux.HandleFunc("/webhook", data.Sock.Webhook)
	// data.mux.HandleFunc("/wsstatus", data.Sock.Status)

	serveCameras(data, cfg, webcamServers)
	handleCameras(data)
	serveHomeData(data)
	handleHomeData(data)
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

	data.Sock.SaveMessages()

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)
	// wsServer.server.Shutdown(ctx)
	return
}

func handleFiles(data *RunData) {
	data.mux.Handle(data.Recorder.Url, data.Recorder)

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
		err := config.Save(data.cfg, "config.json")
		if err != nil {
			log.Println("filesave", err)
			return
		}
	})

	data.mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		err := data.template.Lookup("layout.echo").Execute(w, data)
		if err != nil {
			log.Fatal("/echo", err)
		}
	})

}

func serveCameras(data *RunData, cfg *config.Config, camServers []*camera.Server) {
	data.WebcamHandlers = CreateNexigoHandlers(cfg, data.template)

	for i, camServer := range camServers {
		path := fmt.Sprintf("/%d/", i)
		data.mux.Handle(path, camServer.Stream())

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
		log.Printf("Serving %s\n", path)
	}
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

func handleHomeData(data *RunData) {
	mux := data.mux
	mux.HandleFunc("/sun", handleSun(data))
	mux.HandleFunc("/weather", handleWeather(data))
	mux.HandleFunc("/wifi", handleWifi(data))
	mux.HandleFunc("/lights", handleLights(data))
	handleLightProperties(data)
}

type CameraData struct {
	Action         *config.Action
	WebcamHandlers []*WebcamHandler
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
}

func handleSun(data *RunData) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		sun := data.home.NewSun(data.ActionMap["sun"])
		err := data.template.Lookup("layout.sun").Execute(w, sun)
		if err != nil {
			log.Fatal("/sun", err)
		}
	}
}

func handleWeather(data *RunData) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		forecast := data.home.NewWeather(data.ActionMap["weather"])
		err := data.template.Lookup("layout.weather").Execute(w, forecast)
		if err != nil {
			log.Fatal("/weather", err)
		}
	}
}

func handleWifi(data *RunData) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		sensors := data.home.WifiSensors(data.ActionMap["wifi"])
		err := data.template.Lookup("layout.wifi").Execute(w, sensors)
		if err != nil {
			log.Fatal("/wifi", err)
		}
	}
}
func handleLights(data *RunData) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		lights := data.home.NewLedLights(data.ActionMap["lights"])
		err := data.template.Lookup("layout.lights").Execute(w, lights)
		if err != nil {
			log.Fatal("/lights", err)
		}
	}
}

func handleLightProperties(data *RunData) {
	home := data.home
	readBody := func(r *http.Request) (id string, key string, val string) {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("handleLights.readBody", err)
		}
		lines := strings.Split(string(buf), "&")
		var k, v string
		for _, line := range lines {
			blanksep := strings.Replace(line, "=", " ", 1)
			fmt.Sscan(blanksep, &k, &v)
			if k == "id" {
				id = v
			} else {
				key = k
				val = v
			}
		}
		log.Println("id:", id, "key:", key, "value:", val)
		return
	}

	data.mux.HandleFunc("/light/state",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/state")
			id, key, _ := readBody(r)
			if key == "state" {
				home.CallService(LightCmd(id))
			} else {
				home.CallService(LightCmdOff(id))
			}
		})

	data.mux.HandleFunc("/light/brightness",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/brightness")
			id, key, val := readBody(r)
			home.CallService(LightCmd(id, ServiceData{Key: key, Value: val}))
		})

	data.mux.HandleFunc("/light/color",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/color")
			id, key, val := readBody(r)
			length := len(val)
			if length > 6 {
				val := val[length-6:]
				var red, green, blue int
				fmt.Sscan(fmt.Sprintf("0x%s 0x%s 0x%s", val[:2], val[2:4], val[4:]),
					&red, &green, &blue)
				home.CallService(LightCmd(id, ServiceData{Key: key,
					Value: fmt.Sprintf("[%d,%d,%d]", red, green, blue)}))
			}
		})

	data.mux.HandleFunc("/light/effect",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/effect")
			id, key, val := readBody(r)
			home.CallService(LightCmd(id, ServiceData{Key: key,
				Value: `"` + val + `"`}))
		})
}
