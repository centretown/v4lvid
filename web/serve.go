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
	Controls []*Control
}

func Serve(vservers []*video.Server) {
	const url = "192.168.0.7:9000"
	data := &ServerData{
		Url:      "http://192.168.0.7:9000/0/",
		Controls: NexigoControls,
	}

	// mux := http.NewServeMux()

	for i, vserver := range vservers {
		path := fmt.Sprintf("/%d/", i)
		http.Handle(path, vserver.StreamHook.Stream)
		// hr := NewHandleRecord(vserver)
		// http.Handle(path+"record/", hr)

		go vserver.Serve()
		log.Printf("Serving %s%s\n", url, path)
	}

	source := vservers[0].Source
	webcam, isWebcam := source.(*video.Webcam)
	if isWebcam {
		NewControlList(webcam, 0, data.Controls)
	}

	http.HandleFunc("/record", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("recording...")
		vservers[0].RecordCmd(60)
	})

	fs := http.FileServer(http.Dir("www/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// http.ServeFile(w, r, "www/index.html")
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

func handleCtl(webCam *video.Webcam, ctlKey, downUrl, upUrl string, multiplier int32) {
	info, err := webCam.GetControlInfo(ctlKey)
	if err != nil {
		log.Println("Failed to handle", ctlKey, upUrl, downUrl)
		return
	}

	var (
		max                = info.Max
		min                = info.Min
		currentValue int32 = webCam.GetControlValue(ctlKey)
		step               = info.Step
	)
	if multiplier > 0 {
		step *= multiplier
	}

	http.HandleFunc(upUrl, func(w http.ResponseWriter, r *http.Request) {
		if currentValue-step >= min {
			currentValue -= step
			webCam.SetValue(ctlKey, currentValue)
		}
	})
	http.HandleFunc(downUrl, func(w http.ResponseWriter, r *http.Request) {
		if currentValue+step <= max {
			currentValue += step
			webCam.SetValue(ctlKey, currentValue)
		}
	})

}
