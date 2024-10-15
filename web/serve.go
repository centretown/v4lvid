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

	home := initializeHome()

	var (
		camera_active  bool
		sun_active     bool
		weather_active bool
	)

	http.HandleFunc("/camera",
		func(w http.ResponseWriter, r *http.Request) {
			if camera_active {
				camera_active = false
				return
			}
			w.Header().Add("Cache-Control", "no-cache")
			// err := data.Template.ExecuteTemplate(w, "layout.items", data)
			err := data.Template.Lookup("layout.items").Execute(w, data)
			if err != nil {
				log.Fatal("/camera", err)
			}
			camera_active = true
		})
	http.HandleFunc("/sun",
		func(w http.ResponseWriter, r *http.Request) {
			if sun_active {
				sun_active = false
				return
			}
			w.Header().Add("Cache-Control", "no-cache")
			sensors := home.SunTimes()
			err := data.Template.Lookup("layout.sun").Execute(w, sensors)
			if err != nil {
				log.Fatal("/sun", err)
			}
			sun_active = true
		})
	http.HandleFunc("/weather",
		func(w http.ResponseWriter, r *http.Request) {
			if weather_active {
				weather_active = false
				return
			}
			w.Header().Add("Cache-Control", "no-cache")
			forecast := home.Forecast()
			err := data.Template.Lookup("layout.weather").Execute(w, forecast)
			if err != nil {
				log.Fatal("/sun", err)
			}
			weather_active = true
		})

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

func initializeHome() *ha.HomeData {
	home := ha.NewHomeData()
	ok, err := home.Authorize()
	if err != nil {
		log.Fatal("authorize", err)
	}
	if !ok {
		log.Fatal("not authorized")

	}
	log.Println("Authorized HA")

	err = home.BuildEntities()
	if err != nil {
		log.Fatal("BuildEntities", err)
	}
	log.Println("Build Entities")

	return home
}
