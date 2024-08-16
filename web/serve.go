package web

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"v4lvid/video"
)

type ServerData struct {
	Url      string
	Controls []*ControlHandler
	Record   *RecordControlHandler
}

func Serve(vservers []*video.Server) {
	const url = "192.168.0.7:9000"
	data := &ServerData{
		Url:      "http://192.168.0.7:9000/0/",
		Controls: NexigoControls,
		Record: &RecordControlHandler{
			Server: vservers[0],
			Url:    "/record",
			Icon:   "radio_button_checked",
		},
	}

	// mux := http.NewServeMux()

	for i, vserver := range vservers {
		path := fmt.Sprintf("/%d/", i)
		http.Handle(path, vserver.StreamHook.Stream)

		source := vserver.Source
		webcam, isWebcam := source.(*video.Webcam)
		if isWebcam {
			NewControlList(webcam, 0, data.Controls)
		}

		go vserver.Serve()
		log.Printf("Serving %s%s\n", url, path)
	}

	http.Handle(data.Record.Url, data.Record)

	fs := http.FileServer(http.Dir("www/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		tmpl := template.Must(template.ParseGlob("www/*.html"))
		tmpl.ExecuteTemplate(w, "index.html", data)
	})

	server := &http.Server{
		Addr:         url,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	log.Fatal(server.ListenAndServe())
}
