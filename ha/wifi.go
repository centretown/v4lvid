package ha

import "log"

type Wifi struct {
	Entity[SensorAttributes]
}

// const doorWifi = "sensor.doorstop_wifi"

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
