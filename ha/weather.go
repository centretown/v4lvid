package ha

import "strings"

type Weather struct {
	Entity[WeatherAttributes]
}

func (data *HomeData) Forecast() *Weather {
	wthr := &Weather{}
	entity, ok := data.Entities["weather.forecast_home"]
	if ok {
		wthr.Copy(entity)
	}
	return wthr
}

// func (wthr *Weather) FormatTime() string {
// 	return wthr.LastUpdated.Local().Format(ShortTime)
// }

var weatherIcons = [][]string{
	{"partly_cloudy_day", "partly"},
	{"cloud", "cloud"},
	{"sunny", "sunny"},
	{"rainy", "rainy"},
	{"thunderstorm", "thunder"},
	{"storm", "storm"},
	{"clear_day", "clear"},
}

func (wthr *Weather) Icon() (icon string) {
	for _, keys := range weatherIcons {
		icon = keys[0]
		for _, k := range keys {
			if strings.Contains(wthr.State, k) {
				return
			}
		}
	}
	icon = "routine"
	return
}
