package web

import (
	"log"
)

func CreateNexigoHandlers(rt *RunTime) (handlers []*ControlHandler) {

	handlers = make([]*ControlHandler, 0)
	driver, ok := rt.Config.Drivers[UVCVideo]
	if !ok {
		log.Println("Driver not found", UVCVideo)
		return
	}

	for _, d := range driver {
		handlers = append(handlers,
			NewControlHandler(d.Key, d.Controls, rt))
	}
	return
}
