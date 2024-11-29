package web

import (
	"html/template"
	"log"
	"v4lvid/camera"
	"v4lvid/config"
)

func CreateNexigoHandlers(cfg *config.Config, cameras []*camera.Server,
	tmpl *template.Template) (handlers []*WebcamHandler) {

	const driverKey = "uvcvideo"
	handlers = make([]*WebcamHandler, 0)
	driver, ok := cfg.Drivers[driverKey]
	if !ok {
		log.Println("Driver not found", driverKey)
		return
	}

	for _, d := range driver {
		handlers = append(handlers,
			NewWebcamHandler(d.Key, d.Controls, cameras, tmpl))
	}
	return
}
