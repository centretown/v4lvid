package ha

import (
	"slices"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const ShortTime = "Monday at 15:04:05"

type TimeStampSensor struct {
	Sun       Entity[TimeStampAttributes]
	TimeStamp time.Time
}

func (tss *TimeStampSensor) ShortName() string {
	c := cases.Title(language.Und, cases.NoLower)
	return c.String(strings.
		TrimPrefix(tss.Sun.Attributes.Name, "Sun Next "))
}

func (tss *TimeStampSensor) FormatTime() string {
	return tss.TimeStamp.Local().Format(ShortTime)
}

func (data *HomeData) SunTimes() (suntimes []*TimeStampSensor) {
	sunlist := ListEntitiesWithPrefix("sensor.sun_next", data.EntityKeys)
	suntimes = make([]*TimeStampSensor, 0, len(sunlist))
	for _, s := range sunlist {
		sensor := &TimeStampSensor{}
		sensor.Sun.Copy(data.Entities[s])
		sensor.TimeStamp, _ = time.Parse(time.RFC3339, sensor.Sun.State)
		suntimes = append(suntimes, sensor)
	}
	slices.SortFunc(suntimes, func(a, b *TimeStampSensor) int {
		return a.TimeStamp.Compare(b.TimeStamp)
	})
	return
}
