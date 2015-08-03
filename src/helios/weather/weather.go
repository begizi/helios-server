package weather

import (
	"encoding/json"
	"helios/helios"
	"net/http"
	"os"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"
)

var WeatherChan chan interface{}

func Service() helios.ServiceHandler {
	return func(h *helios.Engine) error {

		apiURL := os.Getenv("WEATHER_API")

		WeatherChan = h.NewBroadcastChannel("weather")

		go initWeatherFetch(apiURL)

		return nil
	}
}

func initWeatherFetch(apiURL string) {
	for {
		weather := getWeather(apiURL)

		WeatherChan <- weather

		// Timeout for 60 seconds before next request
		time.Sleep(60 * time.Second)
	}

	panic("Shouldn't be here")
}

func getWeather(apiURL string) interface{} {
	var weatherBody interface{}

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Warn("Failed to fetch weather", "error", err.Error())
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&weatherBody)
	if err != nil {
		log.Warn("Failed to parse weather body", "error", err.Error())
	}

	return weatherBody
}
