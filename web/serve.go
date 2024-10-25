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
)

type ServerData struct {
	WebcamUrl string
	cfg       *config.Config
	Actions   []*config.Action
	// WebcamServers  []*camera.Server
	WebcamHandlers []*WebcamHandler
	Recorder       *RecordingHandler
	mux            *http.ServeMux
	template       *template.Template
	home           *ha.HomeData
}

func Serve(cfg *config.Config) (data *ServerData) {
	webcamServers := cfg.NewCameraServers()
	data = &ServerData{
		WebcamUrl: "http://192.168.10.7:9000/0/",
		Actions:   cfg.Actions,
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

	serveCameras(data, cfg, webcamServers)
	handleCameras(data)
	serveHomeData(data)
	handleHomeData(data)
	handleFiles(data)

	httpErr := make(chan error, 1)
	go func() {
		httpErr <- httpServer.ListenAndServe()
	}()

	wsServer := NewSocketServer(cfg.WsUrl)
	wsErr := make(chan error, 1)
	go func() {
		wsErr <- wsServer.Run()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-httpErr:
		log.Printf("failed to serve http: %v", err)
	case err := <-wsErr:
		log.Printf("failed to serve websockets: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)
	wsServer.server.Shutdown(ctx)
	return
}

func handleFiles(data *ServerData) {
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
		// buf, err := json.MarshalIndent(data.cfg, "", "  ")

		// f, err := os.Create("config.json")
		// if err != nil {
		// 	log.Println("filesave", err)
		// 	return
		// }
		// defer f.Close()
		// f.Write(buf)
		// log.Println(string(buf))
	})
}

func serveCameras(data *ServerData, cfg *config.Config, camServers []*camera.Server) {
	data.WebcamHandlers = NexigoControlList(cfg, data.template)

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

func serveHomeData(data *ServerData) (err error) {
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

	return
}

func handleHomeData(data *ServerData) {
	mux := data.mux
	mux.HandleFunc("/sun", handleSun(data))
	mux.HandleFunc("/weather", handleWeather(data))
	mux.HandleFunc("/wifi", handleWifi(data))
	mux.HandleFunc("/lights", handleLights(data))
	handleLightProperties(data)
}

func handleCameras(data *ServerData) {
	data.mux.HandleFunc("/camera",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Cache-Control", "no-cache")
			err := data.template.Lookup("layout.controls").Execute(w, data)
			if err != nil {
				log.Fatal("/camera", err)
			}
		})
}

func handleSun(data *ServerData) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		sensors := data.home.SunTimes()
		err := data.template.Lookup("layout.sun").Execute(w, sensors)
		if err != nil {
			log.Fatal("/sun", err)
		}
	}
}

func handleWeather(data *ServerData) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		forecast := data.home.Forecast()
		err := data.template.Lookup("layout.weather").Execute(w, forecast)
		if err != nil {
			log.Fatal("/weather", err)
		}
	}
}

func handleWifi(data *ServerData) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		sensors := data.home.WifiSensors()
		err := data.template.Lookup("layout.wifi").Execute(w, sensors)
		if err != nil {
			log.Fatal("/wifi", err)
		}
	}
}
func handleLights(data *ServerData) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		lights := data.home.LedLights()
		err := data.template.Lookup("layout.lights").Execute(w, lights)
		if err != nil {
			log.Fatal("/lights", err)
		}
	}
}

func handleLightProperties(data *ServerData) {
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
