package ha

import (
	"fmt"
	"log"
)

type Wifi struct {
	Entity[SensorAttributes]
}

func (wifi *Wifi) SignalIcon() string {
	signal := -100
	count, _ := fmt.Sscan(wifi.State, &signal)
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

func (data *HomeData) WifiSensors() (wifis []*Wifi) {
	ids := ListEntitiesLike("wifi", data.EntityKeys)
	wifis = make([]*Wifi, 0, len(ids))
	for _, id := range ids {
		wifi := &Wifi{}
		e, ok := data.Entities[id]
		if ok {
			wifi.Copy(e)
		}
		wifis = append(wifis, wifi)
	}
	return
}

func (data *HomeData) WifiSensor(entityID string) (wifi *Wifi) {
	wifi = &Wifi{}
	e, ok := data.Entities[entityID]
	if !ok {
		log.Println(entityID, "not found")
		return
	}
	wifi.Copy(e)
	return
}
