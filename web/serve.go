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
	"v4lvid/ha"
)

type Action struct {
	Name string
	Icon string
}

type ServerData struct {
	Url             string
	ControlHandlers []*V4lHandler
	Record          *RecordingHandler
	Actions         []Action

	template *template.Template
	home     *ha.HomeData
}

func Serve(vservers []*camera.Server) (data *ServerData) {
	const httpUrl = "192.168.10.7:9000"
	const wsUrl = "192.168.10.7:9900"
	data = &ServerData{
		Url: "http://192.168.10.7:9000/0/",
		Record: &RecordingHandler{
			Server: vservers[0],
			Url:    "/record",
			Icon:   "radio_button_checked",
		},
		Actions: []Action{
			{Name: "camera", Icon: "settings_video_camera"},
			{Name: "sun", Icon: "wb_twilight"},
			{Name: "weather", Icon: "routine"},
			{Name: "wifi", Icon: "wifi"},
			{Name: "lights", Icon: "backlight_high"},
		},
	}
	var (
		err     error
		pattern = "www/*.html"
	)
	data.template, err = template.ParseGlob(pattern)
	if err != nil {
		log.Fatalln("ParseGlob", pattern, err)
	}
	data.ControlHandlers = NexigoControlList(data.template)

	for i, vserver := range vservers {
		path := fmt.Sprintf("/%d/", i)
		http.Handle(path, vserver.StreamHook.Stream)

		source := vserver.Source
		webcam, isWebcam := source.(*camera.Webcam)
		if isWebcam {
			ctll := NewControlList(webcam, 0, data.ControlHandlers)
			http.HandleFunc("/resetcontrols",
				func(w http.ResponseWriter, r *http.Request) {
					ctll.ResetControls()
				})
		}

		go vserver.Serve()
		log.Printf("Serving %s%s\n", httpUrl, path)
	}

	handleCameras(data)

	data.home, err = setupHomeData()
	if err == nil {
		http.HandleFunc("/sun", handleSun(data))
		http.HandleFunc("/weather", handleWeather(data))
		http.HandleFunc("/wifi", handleWifi(data))
		http.HandleFunc("/lights", handleLights(data))
		handleLightProperties(data.home)
	}

	http.Handle(data.Record.Url, data.Record)

	fs := http.FileServer(http.Dir("www/"))
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		data.template.ExecuteTemplate(w, "index.html", data)
	})

	httpServer := &http.Server{
		Addr:         httpUrl,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	httpErr := make(chan error, 1)
	go func() {
		httpErr <- httpServer.ListenAndServe()
	}()

	wsServer := NewSocketServer(wsUrl)
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

func setupHomeData() (home *ha.HomeData, err error) {
	var ok bool
	home = ha.NewHomeData()
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

func handleCameras(data *ServerData) {
	http.HandleFunc("/camera",
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
			log.Fatal("/sun", err)
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

func handleLightProperties(home *ha.HomeData) {
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

	http.HandleFunc("/light/state",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/state")
			id, key, _ := readBody(r)
			if key == "state" {
				home.CallService(LightCmd(id))
			} else {
				home.CallService(LightCmdOff(id))
			}
		})

	http.HandleFunc("/light/brightness",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/brightness")
			id, key, val := readBody(r)
			home.CallService(LightCmd(id, ServiceData{Key: key, Value: val}))
		})

	http.HandleFunc("/light/color",
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

	http.HandleFunc("/light/effect",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/effect")
			id, key, val := readBody(r)
			home.CallService(LightCmd(id, ServiceData{Key: key,
				Value: `"` + val + `"`}))
		})
}
