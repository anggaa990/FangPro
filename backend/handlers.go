package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// ============================================
// FUNCTIONAL PROGRAMMING TYPES
// ============================================

type HandlerFunc func(http.ResponseWriter, *http.Request)
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// ============================================
// PURE FUNCTIONS - No Side Effects
// ============================================

// Get region with default value
func getRegionOrDefault(region string) string {
	if region == "" {
		return "Jember"
	}
	return region
}

// Build recommendation response
func buildRecommendationResponse(result, region string, temp, humidity, rain float64) map[string]interface{} {
	return map[string]interface{}{
		"recommendation": result,
		"region":         region,
		"temperature":    temp,
		"humidity":       humidity,
		"rain_mm":        rain,
	}
}

// Build status response
func buildStatusResponse(status, message string) map[string]string {
	return map[string]string{
		"status":  status,
		"message": message,
	}
}

// ============================================
// RESPONSE HELPER FUNCTIONS
// ============================================

// Generic JSON response writer
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// Generic error response writer
func respondError(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}

// ============================================
// HIGHER ORDER FUNCTIONS - MIDDLEWARE
// ============================================

// Logging middleware
func withLogging(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.URL.RawQuery)
		next(w, r)
	}
}

// Recovery middleware - panic handler
func withRecovery(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				respondError(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// JSON Content-Type middleware
func withJSONContentType(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// Method validation middleware - Higher Order Function that returns middleware
func withMethodValidation(allowedMethods ...string) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for _, method := range allowedMethods {
				if r.Method == method {
					next(w, r)
					return
				}
			}
			respondError(w, "Method tidak didukung", http.StatusMethodNotAllowed)
		}
	}
}

// Error handling middleware wrapper
func withErrorHandling(handler func(http.ResponseWriter, *http.Request) error) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			log.Printf("Handler error: %v", err)
			respondError(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// ============================================
// FUNCTION COMPOSITION
// ============================================

// Chain multiple middlewares - applies from right to left
func chain(handler HandlerFunc, middlewares ...MiddlewareFunc) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// ============================================
// CLOSURE - HANDLER FACTORY
// ============================================

// Factory function for creating handlers with dependency injection
func makeWeatherHandler(fetchWeather func(string) (*WeatherData, error)) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		region := getRegionOrDefault(r.URL.Query().Get("region"))

		data, err := fetchWeather(region)
		if err != nil {
			respondError(w, "Gagal mengambil data cuaca", http.StatusInternalServerError)
			return
		}

		respondJSON(w, http.StatusOK, data)
	}
}

// ============================================
// MAP/FILTER/REDUCE - DATA TRANSFORMATION
// ============================================

// Generic Map function
func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Generic Filter function
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := []T{}
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Generic Reduce function
func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = fn(result, v)
	}
	return result
}

// ============================================
// FUNCTIONAL ERROR HANDLING - Result Type
// ============================================

type Result[T any] struct {
	Value T
	Error error
}

// Create new Result
func NewResult[T any](value T, err error) Result[T] {
	return Result[T]{Value: value, Error: err}
}

// Map transforms the value if no error
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.Error != nil {
		return r
	}
	return Result[T]{Value: fn(r.Value), Error: nil}
}

// OrElse provides a default value if error
func (r Result[T]) OrElse(defaultValue T) T {
	if r.Error != nil {
		return defaultValue
	}
	return r.Value
}

// ============================================
// RECURSION 🔄
// ============================================

// Factorial - Basic recursion example
func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}

// FactorialTailRecursive - Tail recursion (more efficient)
func FactorialTailRecursive(n int) int {
	return factorialHelper(n, 1)
}

func factorialHelper(n, acc int) int {
	if n <= 1 {
		return acc
	}
	return factorialHelper(n-1, n*acc)
}

// Fibonacci - Recursive fibonacci
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// FibonacciMemoized - Optimized with memoization
func FibonacciMemoized(n int) int {
	memo := make(map[int]int)
	return fibHelper(n, memo)
}

func fibHelper(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}
	if val, exists := memo[n]; exists {
		return val
	}
	memo[n] = fibHelper(n-1, memo) + fibHelper(n-2, memo)
	return memo[n]
}

// SumSliceRecursive - Sum array elements recursively
func SumSliceRecursive(slice []int) int {
	if len(slice) == 0 {
		return 0
	}
	return slice[0] + SumSliceRecursive(slice[1:])
}

// FilterRecursive - Recursive filter implementation
func FilterRecursive[T any](slice []T, predicate func(T) bool) []T {
	if len(slice) == 0 {
		return []T{}
	}

	head := slice[0]
	tail := slice[1:]

	if predicate(head) {
		return append([]T{head}, FilterRecursive(tail, predicate)...)
	}
	return FilterRecursive(tail, predicate)
}

// DeepCalculatePriceStats - Recursive calculation for nested price structures
func DeepCalculatePriceStats(prices []Price, depth int) map[string]interface{} {
	if depth <= 0 || len(prices) == 0 {
		return map[string]interface{}{
			"count": 0,
			"sum":   0.0,
		}
	}

	if len(prices) == 1 {
		return map[string]interface{}{
			"count": 1,
			"sum":   prices[0].Price,
		}
	}

	mid := len(prices) / 2
	left := DeepCalculatePriceStats(prices[:mid], depth-1)
	right := DeepCalculatePriceStats(prices[mid:], depth-1)

	return map[string]interface{}{
		"count": left["count"].(int) + right["count"].(int),
		"sum":   left["sum"].(float64) + right["sum"].(float64),
	}
}

// ============================================
// FUNCTIONAL CONCURRENCY 🚀
// ============================================

// ParallelMap - Map dengan goroutines untuk processing parallel
func ParallelMap[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, v := range slice {
		wg.Add(1)
		go func(index int, value T) {
			defer wg.Done()
			transformed := fn(value)
			mu.Lock()
			result[index] = transformed
			mu.Unlock()
		}(i, v)
	}

	wg.Wait()
	return result
}

// ParallelFilter - Filter dengan concurrent processing
func ParallelFilter[T any](slice []T, predicate func(T) bool) []T {
	resultChan := make(chan T, len(slice))
	var wg sync.WaitGroup

	for _, v := range slice {
		wg.Add(1)
		go func(value T) {
			defer wg.Done()
			if predicate(value) {
				resultChan <- value
			}
		}(v)
	}

	// Close channel after all goroutines finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	result := []T{}
	for v := range resultChan {
		result = append(result, v)
	}

	return result
}

// ParallelReduce - Concurrent reduce operation
func ParallelReduce[T any](slice []T, initial T, fn func(T, T) T, workers int) T {
	if len(slice) == 0 {
		return initial
	}

	chunkSize := (len(slice) + workers - 1) / workers
	resultChan := make(chan T, workers)
	var wg sync.WaitGroup

	// Process chunks in parallel
	for i := 0; i < workers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		if start >= len(slice) {
			break
		}

		wg.Add(1)
		go func(chunk []T) {
			defer wg.Done()
			result := initial
			for _, item := range chunk {
				result = fn(result, item)
			}
			resultChan <- result
		}(slice[start:end])
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Combine results
	finalResult := initial
	for partialResult := range resultChan {
		finalResult = fn(finalResult, partialResult)
	}

	return finalResult
}

// FetchMultipleRegionsWeather - Concurrent weather fetching
func FetchMultipleRegionsWeather(regions []string) map[string]*WeatherData {
	results := make(map[string]*WeatherData)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, region := range regions {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			data, err := FetchWeather(r)
			if err != nil {
				log.Printf("Failed to fetch weather for %s: %v", r, err)
				return
			}

			mu.Lock()
			results[r] = data
			mu.Unlock()
		}(region)
	}

	wg.Wait()
	return results
}

// FetchMultiplePricesSources - Concurrent price fetching from multiple sources
func FetchMultiplePricesSources(sources []func() error) []error {
	errorChan := make(chan error, len(sources))
	var wg sync.WaitGroup

	for _, source := range sources {
		wg.Add(1)
		go func(fetchFunc func() error) {
			defer wg.Done()
			if err := fetchFunc(); err != nil {
				errorChan <- err
			}
		}(source)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// Collect errors
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	return errors
}

// Pipeline - Functional pipeline for data processing
type Pipeline[T any] struct {
	input chan T
}

// NewPipeline creates a new pipeline
func NewPipeline[T any](data []T) *Pipeline[T] {
	p := &Pipeline[T]{
		input: make(chan T, len(data)),
	}

	go func() {
		for _, item := range data {
			p.input <- item
		}
		close(p.input)
	}()

	return p
}

// PipeMap - Map operation in pipeline
func PipeMap[T, U any](input chan T, fn func(T) U) chan U {
	output := make(chan U)

	go func() {
		for item := range input {
			output <- fn(item)
		}
		close(output)
	}()

	return output
}

// PipeFilter - Filter operation in pipeline
func PipeFilter[T any](input chan T, predicate func(T) bool) chan T {
	output := make(chan T)

	go func() {
		for item := range input {
			if predicate(item) {
				output <- item
			}
		}
		close(output)
	}()

	return output
}

// CollectFromChannel - Collect all items from channel
func CollectFromChannel[T any](ch chan T) []T {
	result := []T{}
	for item := range ch {
		result = append(result, item)
	}
	return result
}

// WorkerPool - Generic worker pool pattern
type WorkerPool[T, U any] struct {
	workers int
	jobs    chan T
	results chan U
	wg      sync.WaitGroup
}

// NewWorkerPool creates a worker pool
func NewWorkerPool[T, U any](workers int, fn func(T) U) *WorkerPool[T, U] {
	pool := &WorkerPool[T, U]{
		workers: workers,
		jobs:    make(chan T, workers*2),
		results: make(chan U, workers*2),
	}

	// Start workers
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			for job := range pool.jobs {
				pool.results <- fn(job)
			}
		}()
	}

	// Close results channel when all workers done
	go func() {
		pool.wg.Wait()
		close(pool.results)
	}()

	return pool
}

// Submit adds job to the pool
func (wp *WorkerPool[T, U]) Submit(job T) {
	wp.jobs <- job
}

// Close closes the job channel
func (wp *WorkerPool[T, U]) Close() {
	close(wp.jobs)
}

// Results returns the results channel
func (wp *WorkerPool[T, U]) Results() <-chan U {
	return wp.results
}

// ============================================
// REFACTORED HANDLERS
// ============================================

// RecommendationHandler - Basic recommendation endpoint
func RecommendationHandler(w http.ResponseWriter, r *http.Request) {
	handler := chain(
		func(w http.ResponseWriter, r *http.Request) {
			region := getRegionOrDefault(r.URL.Query().Get("region"))

			data, err := FetchWeather(region)
			if err != nil {
				respondError(w, "Gagal mengambil data cuaca", http.StatusInternalServerError)
				return
			}

			result := Recommend(data.Temp, data.Humidity, data.Rain)
			response := buildRecommendationResponse(result, region, data.Temp, float64(data.Humidity), data.Rain)

			respondJSON(w, http.StatusOK, response)
		},
		withJSONContentType,
		withLogging,
		withRecovery,
	)
	handler(w, r)
}

// AdvancedRecommendationHandler - Advanced recommendation endpoint
func AdvancedRecommendationHandler(w http.ResponseWriter, r *http.Request) {
	handler := chain(
		func(w http.ResponseWriter, r *http.Request) {
			region := getRegionOrDefault(r.URL.Query().Get("region"))

			data, err := FetchWeather(region)
			if err != nil {
				respondError(w, "Gagal mengambil data cuaca", http.StatusInternalServerError)
				return
			}

			result := GetAdvancedRecommendation(data.Temp, data.Humidity, data.Rain, region)
			respondJSON(w, http.StatusOK, result)
		},
		withJSONContentType,
		withLogging,
		withRecovery,
	)
	handler(w, r)
}

// WeatherAPIHandler - Using factory pattern with closure
func WeatherAPIHandler(w http.ResponseWriter, r *http.Request) {
	handler := chain(
		makeWeatherHandler(FetchWeather),
		withJSONContentType,
		withLogging,
		withRecovery,
	)
	handler(w, r)
}

// MultiRegionWeatherHandler - Fetch weather for multiple regions concurrently
func MultiRegionWeatherHandler(w http.ResponseWriter, r *http.Request) {
	handler := chain(
		withErrorHandling(func(w http.ResponseWriter, r *http.Request) error {
			regions := []string{"Jember", "Surabaya", "Malang", "Banyuwangi"}

			// Fetch weather concurrently for all regions
			results := FetchMultipleRegionsWeather(regions)

			return respondJSON(w, http.StatusOK, results)
		}),
		withJSONContentType,
		withLogging,
		withRecovery,
	)
	handler(w, r)
}

// AddPriceHandler - With method validation and error handling
func AddPriceHandler(w http.ResponseWriter, r *http.Request) {
	handler := chain(
		withErrorHandling(func(w http.ResponseWriter, r *http.Request) error {
			var p Price
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				respondError(w, "Request body tidak valid", http.StatusBadRequest)
				return nil
			}

			_, err := DB.Exec(`INSERT INTO prices (region, price, unit, source, recorded_at) VALUES (?, ?, ?, ?, ?)`,
				p.Region, p.Price, p.Unit, p.Source, p.RecordedAt)

			if err != nil {
				return err
			}

			response := buildStatusResponse("ok", "Data harga berhasil ditambahkan")
			return respondJSON(w, http.StatusOK, response)
		}),
		withMethodValidation(http.MethodPost),
		withJSONContentType,
		withLogging,
		withRecovery,
	)
	handler(w, r)
}

// FetchPricesHandler - With functional fallback pattern
func FetchPricesHandler(w http.ResponseWriter, r *http.Request) {
	handler := chain(
		withErrorHandling(func(w http.ResponseWriter, r *http.Request) error {
			tryFetch := func() error {
				if err := AutoFetchPricesFromScraper(); err != nil {
					log.Printf("Scraping failed, fallback to simulation: %v", err)
					return AutoFetchPrices()
				}
				return nil
			}

			if err := tryFetch(); err != nil {
				return err
			}

			response := buildStatusResponse("ok", "Berhasil fetch dan simpan harga (Web Scraping + Market Data)")
			return respondJSON(w, http.StatusOK, response)
		}),
		withMethodValidation(http.MethodPost),
		withJSONContentType,
		withLogging,
		withRecovery,
	)
	handler(w, r)
}

// GetCurrentPriceHandler - Get latest price by region
func GetCurrentPriceHandler(w http.ResponseWriter, r *http.Request) {
	handler := chain(
		withErrorHandling(func(w http.ResponseWriter, r *http.Request) error {
			region := getRegionOrDefault(r.URL.Query().Get("region"))

			jsonData, err := GetLatestPriceJSON(region)
			if err != nil {
				return err
			}

			w.Write([]byte(jsonData))
			return nil
		}),
		withJSONContentType,
		withLogging,
		withRecovery,
	)
	handler(w, r)
}

// PricesHandler - Get all prices with functional data transformation
func PricesHandler(w http.ResponseWriter, r *http.Request) {
	handler := chain(
		withErrorHandling(func(w http.ResponseWriter, r *http.Request) error {
			rows, err := DB.Query("SELECT id, region, price, unit, source, recorded_at, created_at FROM prices ORDER BY created_at DESC")
			if err != nil {
				log.Println("DB error:", err)
				return err
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

			if data == nil {
				data = []Price{}
			}

			return respondJSON(w, http.StatusOK, data)
		}),
		withJSONContentType,
		withLogging,
		withRecovery,
	)
	handler(w, r)
}

// ============================================
// BONUS: FUNCTIONAL UTILITIES
// ============================================

// Filter prices by region
func FilterPricesByRegion(prices []Price, region string) []Price {
	return Filter(prices, func(p Price) bool {
		return p.Region == region
	})
}

// Calculate average price
func CalculateAveragePrice(prices []Price) float64 {
	if len(prices) == 0 {
		return 0
	}

	sum := Reduce(prices, 0.0, func(acc float64, p Price) float64 {
		return acc + p.Price
	})

	return sum / float64(len(prices))
}

// Transform prices to simple format
func TransformPricesToSimple(prices []Price) []map[string]interface{} {
	return Map(prices, func(p Price) map[string]interface{} {
		return map[string]interface{}{
			"region": p.Region,
			"price":  p.Price,
			"unit":   p.Unit,
		}
	})
}