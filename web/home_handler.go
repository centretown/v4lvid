package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"v4lvid/homeasst"
)

func (rt *RunTime) homeHandler() {
	mux := rt.mux
	mux.HandleFunc("/sun", rt.handleSun())
	mux.HandleFunc("/weather", rt.handleWeather())
	mux.HandleFunc("/wifi", rt.handleWifi())
	mux.HandleFunc("/lights", rt.handleLights())
	rt.handleLightProperties()
}

func (rt *RunTime) handleSun() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		sun := rt.Home.NewSun(rt.ActionMap["sun"])
		err := rt.template.Lookup("layout.sun").Execute(w, sun)
		if err != nil {
			log.Fatal("/sun", err)
		}
	}
}

func (rt *RunTime) handleWeather() func(http.ResponseWriter, *http.Request) {
	sub := homeasst.NewSubcription(&homeasst.Weather{},
		func(c homeasst.Consumer) {
			w, ok := c.(*homeasst.Weather)
			if ok {
				log.Println("Temperature", w.Attributes.Temperature, w.Attributes.TemperatureUnit)
				rt.Home.Temperature = w.Attributes.Temperature
				rt.Home.TemperatureUnit = w.Attributes.TemperatureUnit
				text := fmt.Sprint(w.Attributes.Temperature, w.Attributes.TemperatureUnit)
				message := `<span id="clock-temp" hx-swap-oob="outerHTML">` + text + `</span>`
				rt.WebSocket.Broadcast(message)
			}
		})
	rt.Home.Subscribe("weather.forecast_home", sub)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		forecast := rt.Home.NewWeather(rt.ActionMap["weather"])
		err := rt.template.Lookup("layout.weather").Execute(w, forecast)
		if err != nil {
			log.Fatal("/weather", err)
		}
	}
}

func (rt *RunTime) handleWifi() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		sensors := rt.Home.WifiSensors(rt.ActionMap["wifi"])
		err := rt.template.Lookup("layout.wifi").Execute(w, sensors)
		if err != nil {
			log.Fatal("/wifi", err)
		}
	}
}

func (rt *RunTime) handleLights() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		lights := rt.Home.NewLedLights(rt.ActionMap["lights"])
		err := rt.template.Lookup("layout.lights").Execute(w, lights)
		if err != nil {
			log.Fatal("/lights", err)
		}
	}
}

func (rt *RunTime) handleLightProperties() {
	home := rt.Home
	readBody := func(r *http.Request) (id string, key string, val string) {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("handleLights.readBody", err)
		}
		lines := strings.Split(string(buf), "&")
		var k, v string
		for _, line := range lines {
			blanksep := strings.Replace(line, "=", " ", 1)
			fmt.Sscan(blanksep, &k, &v)
			if k == "id" {
				id = v
			} else {
				key = k
				val = v
			}
		}
		log.Println("id:", id, "key:", key, "value:", val)
		return
	}

	rt.mux.HandleFunc("/light/state",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/state")
			id, key, _ := readBody(r)
			if key == "state" {
				home.CallService(LightCmd(id))
			} else {
				home.CallService(LightCmdOff(id))
			}
		})

	rt.mux.HandleFunc("/light/brightness",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/brightness")
			id, key, val := readBody(r)
			home.CallService(LightCmd(id, ServiceData{Key: key, Value: val}))
		})

	rt.mux.HandleFunc("/light/color",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/color")
			id, key, val := readBody(r)
			length := len(val)
			if length > 6 {
				val := val[length-6:]
				var red, green, blue int
				fmt.Sscan(fmt.Sprintf("0x%s 0x%s 0x%s", val[:2], val[2:4], val[4:]),
					&red, &green, &blue)
				home.CallService(LightCmd(id, ServiceData{Key: key,
					Value: fmt.Sprintf("[%d,%d,%d]", red, green, blue)}))
			}
		})

	rt.mux.HandleFunc("/light/effect",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/effect")
			id, key, val := readBody(r)
			home.CallService(LightCmd(id, ServiceData{Key: key,
				Value: `"` + val + `"`}))
		})
}
