package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"v4lvid/ha"
)

func handleHomeData(data *RunData) {
	mux := data.mux
	mux.HandleFunc("/sun", handleSun(data))
	mux.HandleFunc("/weather", handleWeather(data))
	mux.HandleFunc("/wifi", handleWifi(data))
	mux.HandleFunc("/lights", handleLights(data))
	handleLightProperties(data)
}

func handleSun(data *RunData) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		sun := data.home.NewSun(data.ActionMap["sun"])
		err := data.template.Lookup("layout.sun").Execute(w, sun)
		if err != nil {
			log.Fatal("/sun", err)
		}
	}
}

func handleWeather(data *RunData) func(http.ResponseWriter, *http.Request) {
	sub := ha.NewSubcription(&ha.Weather{}, func(c ha.Consumer) {
		w, ok := c.(*ha.Weather)
		if ok {
			log.Println("Temperature", w.Attributes.Temperature, w.Attributes.TemperatureUnit)
			data.Temperature = w.Attributes.Temperature
			data.TemperatureUnit = w.Attributes.TemperatureUnit
			text := fmt.Sprint(w.Attributes.Temperature, w.Attributes.TemperatureUnit)
			message := `<span id="clock-temp" hx-swap-oob="outerHTML">` + text + `</span>`
			data.WebSocket.Broadcast(message)
		}
	})
	data.home.Subscribe("weather.forecast_home", sub)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		forecast := data.home.NewWeather(data.ActionMap["weather"])
		err := data.template.Lookup("layout.weather").Execute(w, forecast)
		if err != nil {
			log.Fatal("/weather", err)
		}
	}
}

func handleWifi(data *RunData) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		sensors := data.home.WifiSensors(data.ActionMap["wifi"])
		err := data.template.Lookup("layout.wifi").Execute(w, sensors)
		if err != nil {
			log.Fatal("/wifi", err)
		}
	}
}
func handleLights(data *RunData) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		lights := data.home.NewLedLights(data.ActionMap["lights"])
		err := data.template.Lookup("layout.lights").Execute(w, lights)
		if err != nil {
			log.Fatal("/lights", err)
		}
	}
}

func handleLightProperties(data *RunData) {
	home := data.home
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

	data.mux.HandleFunc("/light/state",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/state")
			id, key, _ := readBody(r)
			if key == "state" {
				home.CallService(LightCmd(id))
			} else {
				home.CallService(LightCmdOff(id))
			}
		})

	data.mux.HandleFunc("/light/brightness",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/brightness")
			id, key, val := readBody(r)
			home.CallService(LightCmd(id, ServiceData{Key: key, Value: val}))
		})

	data.mux.HandleFunc("/light/color",
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

	data.mux.HandleFunc("/light/effect",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("/light/effect")
			id, key, val := readBody(r)
			home.CallService(LightCmd(id, ServiceData{Key: key,
				Value: `"` + val + `"`}))
		})
}
