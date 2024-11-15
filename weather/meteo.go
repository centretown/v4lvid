package weather

const params = `{
	"latitude": 45.4208777,
	"longitude": -75.6901106,
	"hourly": "temperature_2m"
}`

const path = "https://api.open-meteo.com/v1/forecast"

// const responses = await fetchWeatherApi(url, params);

func Forcast() {
	// conn, err := net.Dial("tcp", path)
	// if err != nil {
	// 	log.Println("Forcast error", err)
	// 	return
	// }

	// conn.Write()

	// rdr := bytes.NewReader([]byte(params))
	// req, err := http.NewRequest("GET", path, rdr)
	// if err != nil {
	// 	log.Println("Forcast error", err)
	// 	return
	// }
	// log.Println(req.Method)
}
