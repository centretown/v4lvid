package homeasst

import (
	"fmt"
	"log"
	"v4lvid/config"
)

type WifiSensors struct {
	Entity[SensorAttributes]
}

type Wifi struct {
	Action  *config.Action
	Sensors []*WifiSensors
}

func (ws *WifiSensors) SignalIcon() string {
	signal := -100
	count, _ := fmt.Sscan(ws.State, &signal)
	if count == 0 {
		return "signal_wifi_bad"
	}
	if signal < -67 {
		return "signal_wifi_0_bar"
	}
	if signal < -60 {
		return "network_wifi_1_bar"
	}
	if signal < -50 {
		return "network_wifi_2_bar"
	}
	if signal < -40 {
		return "network_wifi_3_bar"
	}
	return "network_wifi"
}

func (home *HomeRuntime) WifiSensors(action *config.Action) (wifi *Wifi) {
	ids := ListEntitiesLike("wifi", home.EntityKeys)

	sensors := make([]*WifiSensors, 0, len(ids))
	for _, id := range ids {
		sensor := &WifiSensors{}
		e, ok := home.Entities[id]
		if ok {
			sensor.Copy(e)
		}
		sensors = append(sensors, sensor)
	}
	wifi = &Wifi{Action: action, Sensors: sensors}
	return
}

func (home *HomeRuntime) WifiSensor(entityID string) (wifi *WifiSensors) {
	wifi = &WifiSensors{}
	e, ok := home.Entities[entityID]
	if !ok {
		log.Println(entityID, "not found")
		return
	}
	wifi.Copy(e)
	return
}
