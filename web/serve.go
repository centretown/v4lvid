package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"v4lvid/camera"
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

	http.HandleFunc("/togglemenu",
		func(w http.ResponseWriter, r *http.Request) {
		})

	http.Handle(data.Record.Url, data.Record)

	fs := http.FileServer(http.Dir("www/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

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
