package ha

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

func (wthr *Weather) FormatTime() string {
	return wthr.LastUpdated.Local().Format(ShortTime)
}
