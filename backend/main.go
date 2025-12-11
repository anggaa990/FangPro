package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// ============================================
// FUNCTIONAL MIDDLEWARE - CORS
// ============================================

// CORS Middleware - Higher Order Function
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

// ============================================
// PURE FUNCTIONS - CONFIGURATION
// ============================================

// Load environment variables (with side effect isolation)
func loadEnvironment() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("Gagal load .env file, pastikan file ada:", err)
		return err
	}
	log.Println("✓ .env berhasil di-load")
	log.Println("✓ OWM_API_KEY =", os.Getenv("OWM_API_KEY"))
	return nil
}

// ============================================
// FUNCTIONAL ROUTER SETUP
// ============================================

// Route definition type
type Route struct {
	Pattern string
	Handler http.HandlerFunc
	Method  string
}

// Register routes functionally
func registerRoutes(mux *http.ServeMux, routes []Route) {
	for _, route := range routes {
		// Apply CORS to all handlers
		mux.HandleFunc(route.Pattern, enableCORS(route.Handler))
		log.Printf("✓ Registered: %-8s %s", route.Method, route.Pattern)
	}
}

// Define all routes in a declarative way
func getRoutes() []Route {
	return []Route{
		// Price endpoints
		{Pattern: "/harga", Handler: http.HandlerFunc(PricesHandler), Method: "GET"},
		{Pattern: "/harga/add", Handler: http.HandlerFunc(AddPriceHandler), Method: "POST"},
		{Pattern: "/harga/fetch", Handler: http.HandlerFunc(FetchPricesHandler), Method: "POST"},
		{Pattern: "/harga/current", Handler: http.HandlerFunc(GetCurrentPriceHandler), Method: "GET"},
		
		// Weather endpoints
		{Pattern: "/cuaca", Handler: http.HandlerFunc(WeatherAPIHandler), Method: "GET"},
		{Pattern: "/weather", Handler: http.HandlerFunc(WeatherAPIHandler), Method: "GET"},
		{Pattern: "/weather/multi", Handler: http.HandlerFunc(MultiRegionWeatherHandler), Method: "GET"},
		
		// Recommendation endpoints
		{Pattern: "/rekomendasi", Handler: http.HandlerFunc(RecommendationHandler), Method: "GET"},
		{Pattern: "/rekomendasi/advanced", Handler: http.HandlerFunc(AdvancedRecommendationHandler), Method: "GET"},
	}
}

// Print available endpoints
func printEndpoints() {
	separator := "============================================================"
	
	fmt.Println("\n" + separator)
	fmt.Println("🚀 Server berjalan di http://localhost:8080")
	fmt.Println(separator)
	fmt.Println("\n📋 Endpoints tersedia:\n")
	
	endpoints := []struct {
		method      string
		path        string
		description string
	}{
		{"GET", "/harga", "Lihat semua harga"},
		{"POST", "/harga/add", "Tambah harga manual"},
		{"POST", "/harga/fetch", "Fetch harga otomatis (scraping)"},
		{"GET", "/harga/current", "Lihat harga terkini by region"},
		{"GET", "/cuaca", "Data cuaca single region"},
		{"GET", "/weather/multi", "🆕 Data cuaca multiple regions (concurrent)"},
		{"GET", "/rekomendasi", "Rekomendasi sederhana"},
		{"GET", "/rekomendasi/advanced", "Rekomendasi detail"},
	}
	
	for _, ep := range endpoints {
		fmt.Printf("  %-6s %-30s - %s\n", ep.method, ep.path, ep.description)
	}
	
	fmt.Println("\n" + separator)
	fmt.Println("✨ Functional Programming Features:")
	fmt.Println("  ✓ Higher-Order Functions (Middleware)")
	fmt.Println("  ✓ Pure Functions & Immutability")
	fmt.Println("  ✓ Function Composition (chain)")
	fmt.Println("  ✓ Closure & Factory Pattern")
	fmt.Println("  ✓ Map/Filter/Reduce (Generic)")
	fmt.Println("  ✓ Recursion (Factorial, Fibonacci)")
	fmt.Println("  ✓ Functional Concurrency (Goroutines)")
	fmt.Println("  ✓ Pipeline Pattern")
	fmt.Println("  ✓ Worker Pool Pattern")
	fmt.Println(separator + "\n")
}

// ============================================
// MAIN FUNCTION - COMPOSITION
// ============================================

func main() {
	// 1. Load environment (side effect)
	loadEnvironment()
	
	// 2. Initialize database (side effect)
	InitDB()
	defer DB.Close()
	log.Println("✓ Database initialized")
	
	// 3. Setup router
	mux := http.NewServeMux()
	
	// 4. Register routes functionally
	routes := getRoutes()
	registerRoutes(mux, routes)
	
	// 5. Print server info
	printEndpoints()
	
	// 6. Start server
	log.Fatal(http.ListenAndServe(":8080", mux))
}