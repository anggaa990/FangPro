package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type WeatherData struct {
	Temp     float64 `json:"temp"`
	Humidity int     `json:"humidity"`
	Rain     float64 `json:"rain_mm"`
}

// Struct untuk parsing response OpenWeatherMap yang LENGKAP
type OpenWeatherResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Rain struct {
		OneHour   float64 `json:"1h"`
		ThreeHour float64 `json:"3h"`
	} `json:"rain"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Name string `json:"name"`
}

// FetchWeather mengambil data cuaca dari OpenWeatherMap
func FetchWeather(region string) (*WeatherData, error) {
	apiKey := os.Getenv("OWM_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key belum diset")
	}

	// Build URL dengan region sebagai query
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", region, apiKey)

	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("❌ API Error for %s (status %d): %s", region, resp.StatusCode, string(body))
		return nil, fmt.Errorf("API returned status %d for %s", resp.StatusCode, region)
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 🔍 DEBUG: Print raw response
	log.Printf("📡 Raw API response for %s: %s", region, string(body))

	// Parse JSON response
	var apiResp OpenWeatherResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Extract rain data (prioritas 1h, fallback ke 3h)
	rain := apiResp.Rain.OneHour
	if rain == 0 && apiResp.Rain.ThreeHour > 0 {
		rain = apiResp.Rain.ThreeHour / 3.0
	}

	// 🔍 DEBUG: Print parsed rain data
	log.Printf("☔ Rain data for %s: 1h=%.2fmm, 3h=%.2fmm, final=%.2fmm", 
		region, apiResp.Rain.OneHour, apiResp.Rain.ThreeHour, rain)

	// Get weather condition
	weatherCondition := ""
	if len(apiResp.Weather) > 0 {
		weatherCondition = apiResp.Weather[0].Main
	}

	// Log weather summary
	log.Printf("🌤️  Weather fetched: %s - temp=%.1f°C, humidity=%d%%, rain=%.2fmm, condition=%s", 
		region, apiResp.Main.Temp, apiResp.Main.Humidity, rain, weatherCondition)

	// Simpan ke database secara ASYNC (non-blocking)
	go func() {
		_, err := DB.Exec(`INSERT INTO weather_history (region, temp_c, humidity, rain_mm, fetched_at)
			VALUES (?, ?, ?, ?, ?)`, region, apiResp.Main.Temp, apiResp.Main.Humidity, rain, time.Now())
		if err != nil {
			log.Printf("⚠️  Warning - Gagal menyimpan history cuaca untuk %s: %v", region, err)
		} else {
			log.Printf("✅ Weather history saved: %s (%.1f°C, %d%%, %.2fmm)", 
				region, apiResp.Main.Temp, apiResp.Main.Humidity, rain)
		}
	}()

	return &WeatherData{
		Temp:     apiResp.Main.Temp,
		Humidity: apiResp.Main.Humidity,
		Rain:     rain,
	}, nil
}

// FetchWeatherForecast - Bonus: ambil data forecast untuk cek rain prediction
func FetchWeatherForecast(region string) ([]WeatherData, error) {
	apiKey := os.Getenv("OWM_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key belum diset")
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric", region, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var forecastResp struct {
		List []struct {
			Main struct {
				Temp     float64 `json:"temp"`
				Humidity int     `json:"humidity"`
			} `json:"main"`
			Rain struct {
				ThreeHour float64 `json:"3h"`
			} `json:"rain"`
		} `json:"list"`
	}

	if err := json.Unmarshal(body, &forecastResp); err != nil {
		return nil, err
	}

	var forecasts []WeatherData
	for _, item := range forecastResp.List {
		forecasts = append(forecasts, WeatherData{
			Temp:     item.Main.Temp,
			Humidity: item.Main.Humidity,
			Rain:     item.Rain.ThreeHour,
		})
	}

	log.Printf("📊 Forecast data retrieved for %s: %d entries", region, len(forecasts))

	return forecasts, nil
}