// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package openmeteo_sdk

import "strconv"

type Unit byte

const (
	Unitundefined                           Unit = 0
	Unitcelsius                             Unit = 1
	Unitcentimetre                          Unit = 2
	Unitcubic_metre_per_cubic_metre         Unit = 3
	Unitcubic_metre_per_second              Unit = 4
	Unitdegree_direction                    Unit = 5
	Unitdimensionless_integer               Unit = 6
	Unitdimensionless                       Unit = 7
	Uniteuropean_air_quality_index          Unit = 8
	Unitfahrenheit                          Unit = 9
	Unitfeet                                Unit = 10
	Unitfraction                            Unit = 11
	Unitgdd_celsius                         Unit = 12
	Unitgeopotential_metre                  Unit = 13
	Unitgrains_per_cubic_metre              Unit = 14
	Unitgram_per_kilogram                   Unit = 15
	Unithectopascal                         Unit = 16
	Unithours                               Unit = 17
	Unitinch                                Unit = 18
	Unitiso8601                             Unit = 19
	Unitjoule_per_kilogram                  Unit = 20
	Unitkelvin                              Unit = 21
	Unitkilopascal                          Unit = 22
	Unitkilogram_per_square_metre           Unit = 23
	Unitkilometres_per_hour                 Unit = 24
	Unitknots                               Unit = 25
	Unitmegajoule_per_square_metre          Unit = 26
	Unitmetre_per_second_not_unit_converted Unit = 27
	Unitmetre_per_second                    Unit = 28
	Unitmetre                               Unit = 29
	Unitmicrograms_per_cubic_metre          Unit = 30
	Unitmiles_per_hour                      Unit = 31
	Unitmillimetre                          Unit = 32
	Unitpascal                              Unit = 33
	Unitper_second                          Unit = 34
	Unitpercentage                          Unit = 35
	Unitseconds                             Unit = 36
	Unitunix_time                           Unit = 37
	Unitus_air_quality_index                Unit = 38
	Unitwatt_per_square_metre               Unit = 39
	Unitwmo_code                            Unit = 40
	Unitparts_per_million                   Unit = 41
)

var EnumNamesUnit = map[Unit]string{
	Unitundefined:                           "undefined",
	Unitcelsius:                             "celsius",
	Unitcentimetre:                          "centimetre",
	Unitcubic_metre_per_cubic_metre:         "cubic_metre_per_cubic_metre",
	Unitcubic_metre_per_second:              "cubic_metre_per_second",
	Unitdegree_direction:                    "degree_direction",
	Unitdimensionless_integer:               "dimensionless_integer",
	Unitdimensionless:                       "dimensionless",
	Uniteuropean_air_quality_index:          "european_air_quality_index",
	Unitfahrenheit:                          "fahrenheit",
	Unitfeet:                                "feet",
	Unitfraction:                            "fraction",
	Unitgdd_celsius:                         "gdd_celsius",
	Unitgeopotential_metre:                  "geopotential_metre",
	Unitgrains_per_cubic_metre:              "grains_per_cubic_metre",
	Unitgram_per_kilogram:                   "gram_per_kilogram",
	Unithectopascal:                         "hectopascal",
	Unithours:                               "hours",
	Unitinch:                                "inch",
	Unitiso8601:                             "iso8601",
	Unitjoule_per_kilogram:                  "joule_per_kilogram",
	Unitkelvin:                              "kelvin",
	Unitkilopascal:                          "kilopascal",
	Unitkilogram_per_square_metre:           "kilogram_per_square_metre",
	Unitkilometres_per_hour:                 "kilometres_per_hour",
	Unitknots:                               "knots",
	Unitmegajoule_per_square_metre:          "megajoule_per_square_metre",
	Unitmetre_per_second_not_unit_converted: "metre_per_second_not_unit_converted",
	Unitmetre_per_second:                    "metre_per_second",
	Unitmetre:                               "metre",
	Unitmicrograms_per_cubic_metre:          "micrograms_per_cubic_metre",
	Unitmiles_per_hour:                      "miles_per_hour",
	Unitmillimetre:                          "millimetre",
	Unitpascal:                              "pascal",
	Unitper_second:                          "per_second",
	Unitpercentage:                          "percentage",
	Unitseconds:                             "seconds",
	Unitunix_time:                           "unix_time",
	Unitus_air_quality_index:                "us_air_quality_index",
	Unitwatt_per_square_metre:               "watt_per_square_metre",
	Unitwmo_code:                            "wmo_code",
	Unitparts_per_million:                   "parts_per_million",
}

var EnumValuesUnit = map[string]Unit{
	"undefined":                           Unitundefined,
	"celsius":                             Unitcelsius,
	"centimetre":                          Unitcentimetre,
	"cubic_metre_per_cubic_metre":         Unitcubic_metre_per_cubic_metre,
	"cubic_metre_per_second":              Unitcubic_metre_per_second,
	"degree_direction":                    Unitdegree_direction,
	"dimensionless_integer":               Unitdimensionless_integer,
	"dimensionless":                       Unitdimensionless,
	"european_air_quality_index":          Uniteuropean_air_quality_index,
	"fahrenheit":                          Unitfahrenheit,
	"feet":                                Unitfeet,
	"fraction":                            Unitfraction,
	"gdd_celsius":                         Unitgdd_celsius,
	"geopotential_metre":                  Unitgeopotential_metre,
	"grains_per_cubic_metre":              Unitgrains_per_cubic_metre,
	"gram_per_kilogram":                   Unitgram_per_kilogram,
	"hectopascal":                         Unithectopascal,
	"hours":                               Unithours,
	"inch":                                Unitinch,
	"iso8601":                             Unitiso8601,
	"joule_per_kilogram":                  Unitjoule_per_kilogram,
	"kelvin":                              Unitkelvin,
	"kilopascal":                          Unitkilopascal,
	"kilogram_per_square_metre":           Unitkilogram_per_square_metre,
	"kilometres_per_hour":                 Unitkilometres_per_hour,
	"knots":                               Unitknots,
	"megajoule_per_square_metre":          Unitmegajoule_per_square_metre,
	"metre_per_second_not_unit_converted": Unitmetre_per_second_not_unit_converted,
	"metre_per_second":                    Unitmetre_per_second,
	"metre":                               Unitmetre,
	"micrograms_per_cubic_metre":          Unitmicrograms_per_cubic_metre,
	"miles_per_hour":                      Unitmiles_per_hour,
	"millimetre":                          Unitmillimetre,
	"pascal":                              Unitpascal,
	"per_second":                          Unitper_second,
	"percentage":                          Unitpercentage,
	"seconds":                             Unitseconds,
	"unix_time":                           Unitunix_time,
	"us_air_quality_index":                Unitus_air_quality_index,
	"watt_per_square_metre":               Unitwatt_per_square_metre,
	"wmo_code":                            Unitwmo_code,
	"parts_per_million":                   Unitparts_per_million,
}

func (v Unit) String() string {
	if s, ok := EnumNamesUnit[v]; ok {
		return s
	}
	return "Unit(" + strconv.FormatInt(int64(v), 10) + ")"
}
