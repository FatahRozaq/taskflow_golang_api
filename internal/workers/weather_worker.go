package workers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func FetchWeather() {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	city := "Jakarta"
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching weather:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	fmt.Println("Weather Data:", data)
}

func StartWeatherScheduler() {
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			<-ticker.C
			FetchWeather()
		}
	}()
}
