package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"v4lvid/camera"
	"v4lvid/ha"
)

type ServerData struct {
	Url             string
	ControlHandlers []*V4lHandler
	Record          *RecordControlHandler
	Template        *template.Template
	HasControls     bool
	Home            *ha.HomeData
}

func Serve(vservers []*camera.Server) (data *ServerData) {
	const url = "192.168.10.7:9000"
	data = &ServerData{
		Url: "http://192.168.10.7:9000/0/",
		Record: &RecordControlHandler{
			Server: vservers[0],
			Url:    "/record",
			Icon:   "radio_button_checked",
		},
	}
	var err error
	data.Template, err = template.ParseGlob("www/*.html")
	if err != nil {
		log.Fatalln("Parse", err)
	}
	data.ControlHandlers = NexigoControlList(data.Template)

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
		log.Printf("Serving %s%s\n", url, path)
	}

	handleCameras(data)

	data.Home, err = initializeHome()
	if err == nil {
		handleHome(data)
	}

	http.Handle(data.Record.Url, data.Record)

	fs := http.FileServer(http.Dir("www/"))
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		data.Template.ExecuteTemplate(w, "index.html", data)
	})

	server := &http.Server{
		Addr:         url,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	log.Fatal(server.ListenAndServe())
	return
}

func initializeHome() (home *ha.HomeData, err error) {
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
			err := data.Template.Lookup("layout.controls").Execute(w, data)
			if err != nil {
				log.Fatal("/camera", err)
			}
		})
}

func handleHome(data *ServerData) {
	home := data.Home

	http.HandleFunc("/sun",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Cache-Control", "no-cache")
			sensors := home.SunTimes()
			err := data.Template.Lookup("layout.sun").Execute(w, sensors)
			if err != nil {
				log.Fatal("/sun", err)
			}
		})
	http.HandleFunc("/weather",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Cache-Control", "no-cache")
			forecast := home.Forecast()
			err := data.Template.Lookup("layout.weather").Execute(w, forecast)
			if err != nil {
				log.Fatal("/sun", err)
			}
		})
	http.HandleFunc("/wifi",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Cache-Control", "no-cache")
			sensors := home.WifiSensors()
			err := data.Template.Lookup("layout.wifi").Execute(w, sensors)
			if err != nil {
				log.Fatal("/wifi", err)
			}
		})
}
