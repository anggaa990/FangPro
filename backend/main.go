package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/joho/godotenv"
    "os"
)

// CORS Middleware
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Handle preflight request
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}

func main() {
    // Load file .env
    err := godotenv.Load()
    if err != nil {
        log.Println("Gagal load .env file, pastikan file ada:", err)
    } else {
        log.Println(".env berhasil di-load")
    }

    // Debug: cek apakah API key terbaca
    log.Println("OWM_API_KEY =", os.Getenv("OWM_API_KEY"))

    // Inisialisasi database
    InitDB()
    defer DB.Close()

    mux := http.NewServeMux()

    // Handler dengan CORS
    // Price endpoints
    mux.HandleFunc("/harga", enableCORS(PricesHandler))
    mux.HandleFunc("/harga/add", enableCORS(AddPriceHandler))
    mux.HandleFunc("/harga/fetch", enableCORS(FetchPricesHandler))        // NEW: Auto fetch prices
    mux.HandleFunc("/harga/current", enableCORS(GetCurrentPriceHandler))  // NEW: Get current price
    
    // Weather endpoints
    mux.HandleFunc("/cuaca", enableCORS(WeatherAPIHandler))
    mux.HandleFunc("/weather", enableCORS(WeatherAPIHandler))
    
    // Recommendation endpoints
    mux.HandleFunc("/rekomendasi", enableCORS(RecommendationHandler))                  // Simple
    mux.HandleFunc("/rekomendasi/advanced", enableCORS(AdvancedRecommendationHandler)) // Detailed

    fmt.Println("Server berjalan di http://localhost:8080")
    fmt.Println("Endpoints tersedia:")
    fmt.Println("  GET  /harga           - Lihat semua harga")
    fmt.Println("  POST /harga/add       - Tambah harga manual")
    fmt.Println("  POST /harga/fetch     - Fetch harga otomatis (scraping)")
    fmt.Println("  GET  /harga/current   - Lihat harga terkini by region")
    fmt.Println("  GET  /cuaca           - Data cuaca")
    fmt.Println("  GET  /rekomendasi     - Rekomendasi sederhana")
    fmt.Println("  GET  /rekomendasi/advanced - Rekomendasi detail")
    
    log.Fatal(http.ListenAndServe(":8080", mux))
}