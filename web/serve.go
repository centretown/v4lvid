package web

import (
	"fmt"
	"log"
	"net/http"
	"v4lvid/video"
)

func Serve(vservers []*video.Server) {
	const url = "192.168.0.7:9000"

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
	webCam, isWebcam := source.(*video.Webcam)
	if isWebcam {
		handleCtl(webCam, "Zoom, Absolute", "/zoomin", "/zoomout", 1)
		handleCtl(webCam, "Pan, Absolute", "/panright", "/panleft", 0)
		handleCtl(webCam, "Tilt, Absolute", "/pandown", "/panup", 0)
		handleCtl(webCam, "Brightness", "/brightnessup", "/brightnessdown", 10)
		handleCtl(webCam, "Contrast", "/contrastup", "/contrastdown", 10)
		handleCtl(webCam, "Saturation", "/saturationup", "/saturationdown", 10)
	}

	http.HandleFunc("/record", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("recording...")
		vservers[0].RecordCmd(60)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "www/index.html")
	})

	server := &http.Server{
		Addr:         url,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	log.Fatal(server.ListenAndServe())

}

func handleCtl(webCam *video.Webcam, ctlKey, downUrl, upUrl string, step int32) {
	info, err := webCam.GetControlInfo(ctlKey)
	if err != nil {
		log.Println("Failed to handle", ctlKey, upUrl, downUrl)
		return
	}

	var (
		max                = info.Max
		min                = info.Min
		currentValue int32 = webCam.GetControlValue(ctlKey)
	)
	if step == 0 || (step >= info.Step && step%info.Step != 0) {
		step = info.Step
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
