package ha

import (
	"slices"
	"strings"
	"time"
	"v4lvid/config"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const ShortTime = "Monday at 15:04:05"

type TimeStampSensor struct {
	Sun       Entity[TimeStampAttributes]
	TimeStamp time.Time
	Action    *config.Action
}

func (tss *TimeStampSensor) ShortName() string {
	c := cases.Title(language.Und, cases.NoLower)
	return c.String(strings.
		TrimPrefix(tss.Sun.Attributes.Name, "Sun Next "))
}

func (tss *TimeStampSensor) FormatTime() string {
	return tss.TimeStamp.Local().Format(ShortTime)
}

type Sun struct {
	Action  *config.Action
	Sensors []*TimeStampSensor
}

func (data *HomeData) NewSun(action *config.Action) (sun *Sun) {
	sunlist := ListEntitiesLike("sensor.sun_next", data.EntityKeys)
	sensors := make([]*TimeStampSensor, 0, len(sunlist))
	for _, s := range sunlist {
		sensor := &TimeStampSensor{}
		sensor.Sun.Copy(data.Entities[s])
		sensor.TimeStamp, _ = time.Parse(time.RFC3339, sensor.Sun.State)
		sensors = append(sensors, sensor)
	}
	slices.SortFunc(sensors, func(a, b *TimeStampSensor) int {
		return a.TimeStamp.Compare(b.TimeStamp)
	})
	sun = &Sun{Action: action, Sensors: sensors}
	return
}
