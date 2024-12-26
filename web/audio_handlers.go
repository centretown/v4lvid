package web

import (
	"log"
	"net/http"
	"v4lvid/config"

	"github.com/gordonklaus/portaudio"
)

type AudioParm struct {
	Action  *config.Action
	Devices []*portaudio.DeviceInfo
	Current *portaudio.DeviceInfo
	Enabled bool
}

func (rt *RunTime) handleAudio() {
	rt.mux.HandleFunc("/audio_settings", func(w http.ResponseWriter, r *http.Request) {
		var parm = &AudioParm{
			Action:  rt.ActionMap["audio_settings"],
			Devices: rt.AudioMgr.FindDevices("usb"),
			Enabled: rt.AudioMgr.Enabled,
			Current: rt.AudioMgr.Current,
		}
		if len(parm.Devices) > 0 && parm.Current == nil {
			parm.Current = parm.Devices[0]
			rt.AudioMgr.Current = parm.Current
		}
		err := rt.template.Lookup("layout.audio.list").Execute(w, parm)
		if err != nil {
			log.Println("audio_settings", err)
			return
		}
	})

	rt.mux.HandleFunc("/audio/enable", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		selection := r.FormValue("audio_enabled")
		if len(selection) == 0 {
			rt.AudioMgr.Enabled = false
			log.Printf("audio disabled")
			return
		}
		rt.AudioMgr.Enabled = true
		log.Printf("audio enabled")
	})
	rt.mux.HandleFunc("/audio/select", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		selection := r.FormValue("audio_select")
		if len(selection) == 0 {
			log.Println("form value 'audio_select' not found")
			return
		}

		device, err := rt.AudioMgr.FindDevice(selection)
		if err != nil {
			log.Println("audio/select", err)
			return
		}
		rt.AudioMgr.Current = device
		log.Printf("Found '%s'\n", device.Name)
	})
}
