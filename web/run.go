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
	"strconv"
	"strings"
	"time"
	"v4lvid/camera"
	"v4lvid/config"
	"v4lvid/ha"
	"v4lvid/socket"
)

type RunData struct {
	WebcamUrl string
	Config    *config.Config
	Actions   []*config.Action
	ActionMap map[string]*config.Action
	// WebcamServers  []*camera.Server
	WebcamHandlers  []*WebcamHandler
	CameraMap       map[string]*camera.VideoConfig
	Streamer        *Streamer
	Temperature     float64
	TemperatureUnit string
	mux             *http.ServeMux
	template        *template.Template
	home            *ha.HomeData
	HomeActive      bool
	Socket          *socket.Server
}

func Run(cfg *config.Config) (data *RunData) {
	data = &RunData{
		WebcamUrl: "http://192.168.10.7:9000/0/",
		Config:    cfg,
		Actions:   cfg.Actions,
		ActionMap: cfg.NewActionMap(),
		CameraMap: cfg.NewCameraMap(),
		mux:       &http.ServeMux{},
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

	data.WebcamHandlers = CreateNexigoHandlers(cfg, data.template)

	data.Socket = socket.NewServer(data.template)
	data.Socket.LoadMessages()
	data.Socket.Run()

	webcamServers := cfg.NewCameraServers(data.Socket)
	data.Streamer = &Streamer{
		Server: webcamServers[0],
		Url:    "/record",
		Icon:   "radio_button_checked",
		Sock:   data.Socket,
	}

	data.mux.HandleFunc("/events", data.Socket.Events)
	data.mux.HandleFunc("/webhook", data.Socket.Webhook)

	serveCameras(data, webcamServers)
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

	data.Socket.SaveMessages()

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

	data.mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		err := data.template.Lookup("layout.echo").Execute(w, data)
		if err != nil {
			log.Fatal("/echo", err)
		}
	})

}

func serveCameras(data *RunData, camServers []*camera.Server) {

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

type CameraListData struct {
	Action  *config.Action
	Cameras []*camera.VideoConfig
}

func listCameraHandler(data *RunData) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		camData := &CameraListData{
			Action:  data.ActionMap["camera_list"],
			Cameras: data.Config.Cameras,
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
		// w.Header().Add("Cache-Control", "no-cache")
		err := r.ParseForm()
		if err != nil {
			log.Println("/camera_add", err)
		}
		for k, v := range r.Form {
			log.Println(k, v)
		}
		path := fmt.Sprintf("%s.%s:%s",
			r.FormValue("camera_net"),
			r.FormValue("camera_suffix"),
			r.FormValue("camera_port"))

		cam, ok := data.CameraMap[path]
		if ok {
			log.Println("Camera Found", cam.Path)
		} else {
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
			data.CameraMap[path] = vc
			data.Config.AddCamera(vc)
			log.Println("Camera Added", vc)
		}
	}
}

func handleWeather(data *RunData) func(http.ResponseWriter, *http.Request) {
	sub := ha.NewSubcription(&ha.Weather{}, func(c ha.Consumer) {
		w, ok := c.(*ha.Weather)
		if ok {
			log.Println("Temperature", w.Attributes.Temperature, w.Attributes.TemperatureUnit)
			data.Temperature = w.Attributes.Temperature
			data.TemperatureUnit = w.Attributes.TemperatureUnit
			text := fmt.Sprint(w.Attributes.Temperature, w.Attributes.TemperatureUnit)
			message := `<span id="clock-temp" hx-swap-oob="outerHTML">` + text + `</span>`
			data.Socket.Broadcast(message)
		}
	})
	data.home.Subscribe("weather.forecast_home", sub)

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
