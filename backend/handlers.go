package main

import (
    "encoding/json"
    "log"
    "net/http"
)

func RecommendationHandler(w http.ResponseWriter, r *http.Request) {
    region := r.URL.Query().Get("region")
    if region == "" {
        region = "Jember"
    }

    // Fetch weather data
    data, err := FetchWeather(region)
    if err != nil {
        http.Error(w, "Gagal mengambil data cuaca", http.StatusInternalServerError)
        return
    }

    result := Recommend(data.Temp, data.Humidity, data.Rain)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "recommendation": result,
        "region":         region,
        "temperature":    data.Temp,
        "humidity":       data.Humidity,
        "rain_mm":        data.Rain,
    })
}

func AdvancedRecommendationHandler(w http.ResponseWriter, r *http.Request) {
    region := r.URL.Query().Get("region")
    if region == "" {
        region = "Jember"
    }

    // Fetch weather data
    data, err := FetchWeather(region)
    if err != nil {
        http.Error(w, "Gagal mengambil data cuaca", http.StatusInternalServerError)
        return
    }

    result := GetAdvancedRecommendation(data.Temp, data.Humidity, data.Rain, region)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func WeatherAPIHandler(w http.ResponseWriter, r *http.Request) {
    region := r.URL.Query().Get("region")
    if region == "" {
        region = "Jember"
    }

    data, err := FetchWeather(region)
    if err != nil {
        http.Error(w, "Gagal mengambil data cuaca", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(data)
}

func AddPriceHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method tidak didukung", http.StatusMethodNotAllowed)
        return
    }

    var p Price
    err := json.NewDecoder(r.Body).Decode(&p)
    if err != nil {
        http.Error(w, "Request body tidak valid", http.StatusBadRequest)
        return
    }

    _, err = DB.Exec(`INSERT INTO prices (region, price, unit, source, recorded_at) VALUES (?, ?, ?, ?, ?)`,
        p.Region, p.Price, p.Unit, p.Source, p.RecordedAt)
    if err != nil {
        http.Error(w, "Gagal menyimpan data", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status":  "ok",
        "message": "Data harga berhasil ditambahkan",
    })
}

func FetchPricesHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method tidak didukung", http.StatusMethodNotAllowed)
        return
    }

    // Try scraping first
    err := AutoFetchPricesFromScraper()
    if err != nil {
        log.Printf("Scraping failed, fallback to simulation: %v", err)
        // Fallback to simulation
        err = AutoFetchPrices()
        if err != nil {
            http.Error(w, "Gagal fetch harga: "+err.Error(), http.StatusInternalServerError)
            return
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status":  "ok",
        "message": "Berhasil fetch dan simpan harga (Web Scraping + Market Data)",
    })
}

func GetCurrentPriceHandler(w http.ResponseWriter, r *http.Request) {
    region := r.URL.Query().Get("region")
    if region == "" {
        region = "Jember"
    }

    jsonData, err := GetLatestPriceJSON(region)
    if err != nil {
        http.Error(w, "Gagal mendapatkan harga: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(jsonData))
}

func PricesHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := DB.Query("SELECT id, region, price, unit, source, recorded_at, created_at FROM prices ORDER BY created_at DESC")
    if err != nil {
        http.Error(w, "DB error", 500)
        log.Println("DB error:", err)
        return
    }
    defer rows.Close()

    var data []Price

    for rows.Next() {
        var p Price
        err := rows.Scan(&p.ID, &p.Region, &p.Price, &p.Unit, &p.Source, &p.RecordedAt, &p.CreatedAt)
        if err != nil {
            log.Println("Scan error:", err)
            continue
        }
        data = append(data, p)
    }

    // Return empty array instead of null
    if data == nil {
        data = []Price{}
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}