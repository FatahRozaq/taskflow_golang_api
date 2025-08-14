package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Name string `json:"name"`
}

var (
	weatherCache      *WeatherResponse
	weatherCacheMutex = &sync.RWMutex{}
	lastSyncTime      time.Time
)

func SyncWeatherData() {
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	city := os.Getenv("WEATHER_CITY")
	if apiKey == "" || city == "" {
		log.Println("Weather worker: API key or city is not set. Skipping.")
		return
	}

	log.Println("Weather worker: Starting data sync...")

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Weather worker: Failed to fetch data: %v", err)

		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Weather worker: API returned non-200 status: %s", resp.Status)
		return
	}

	var data WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Weather worker: Failed to decode JSON response: %v", err)
		return
	}

	weatherCacheMutex.Lock()
	weatherCache = &data
	lastSyncTime = time.Now()
	weatherCacheMutex.Unlock()

	log.Printf("Weather worker: Successfully synced weather data for %s. Temp: %.1f°C", data.Name, data.Main.Temp)
}

func GetCachedWeatherData() (*WeatherResponse, time.Time) {
	weatherCacheMutex.RLock()
	defer weatherCacheMutex.RUnlock()
	return weatherCache, lastSyncTime
}
