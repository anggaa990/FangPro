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

// FetchWeather mengambil data cuaca dari OpenWeatherMap
func FetchWeather(region string) (*WeatherData, error) {
    apiKey := os.Getenv("OWM_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("API key belum diset")
    }

    // contoh pakai kota sebagai region
    url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", region, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    var apiResp map[string]interface{}
    if err := json.Unmarshal(body, &apiResp); err != nil {
        return nil, err
    }

    // parsing temperatur & humidity
    mainData := apiResp["main"].(map[string]interface{})
    temp := mainData["temp"].(float64)
    humidity := int(mainData["humidity"].(float64))

    // parsing curah hujan (opsional)
    rain := 0.0
    if r, ok := apiResp["rain"].(map[string]interface{}); ok {
        if h1, ok := r["1h"].(float64); ok {
            rain = h1
        }
    }

    // Simpan ke database secara ASYNC (non-blocking)
    // Ini mencegah API response terlambat karena menunggu database write
    go func() {
        _, err := DB.Exec(`INSERT INTO weather_history (region, temp_c, humidity, rain_mm, fetched_at)
            VALUES (?, ?, ?, ?, ?)`, region, temp, humidity, rain, time.Now())
        if err != nil {
            log.Println("Warning - Gagal menyimpan history cuaca:", err)
        } else {
            log.Printf("✓ Weather history saved: %s (%.1f°C, %d%%, %.1fmm)", region, temp, humidity, rain)
        }
    }()

    return &WeatherData{
        Temp:     temp,
        Humidity: humidity,
        Rain:     rain,
    }, nil
}