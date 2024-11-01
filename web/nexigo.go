package web

import (
	"html/template"
	"log"
	"v4lvid/config"
)

func CreateNexigoHandlers(cfg *config.Config, tmpl *template.Template) (handlers []*WebcamHandler) {
	handlers = make([]*WebcamHandler, 0)
	const key = "uvcvideo"
	driver, ok := cfg.Drivers[key]
	if !ok {
		log.Println("Driver not found", key)
		return
	}

	for _, d := range driver {
		handlers = append(handlers,
			NewWebcamHandler(d.Key, d.Controls, tmpl))
	}
	return
}
