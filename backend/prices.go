package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "time"
)

type Price struct {
    ID         int     `json:"id"`
    Region     string  `json:"region"`
    Price      float64 `json:"price"`
    Unit       string  `json:"unit"`
    Source     string  `json:"source"`
    RecordedAt string  `json:"recorded_at"`
    CreatedAt  string  `json:"created_at"`
}

// AutoFetchPrices simulates fetching prices and saves to database
func AutoFetchPrices() error {
    regions := []string{"Jember", "Malang", "Surabaya", "Bondowoso"}
    source := "Market Data API"
    
    for _, region := range regions {
        // Simulate price data (5000-8000 per kg)
        price := 5000 + rand.Intn(3000)
        recordedAt := time.Now().Format("2006-01-02 15:04:05")
        
        _, err := DB.Exec(`INSERT INTO prices (region, price, unit, source, recorded_at) VALUES (?, ?, ?, ?, ?)`,
            region, price, "per kg", source, recordedAt)
        if err != nil {
            log.Printf("Failed to insert price for %s: %v", region, err)
            return err
        }
        
        log.Printf("Inserted price for %s: Rp %d/kg", region, price)
    }
    
    return nil
}

// GetLatestPriceJSON returns the latest price for a region as JSON string
func GetLatestPriceJSON(region string) (string, error) {
    var p Price
    
    err := DB.QueryRow(`
        SELECT id, region, price, unit, source, recorded_at, created_at 
        FROM prices 
        WHERE region = ? 
        ORDER BY created_at DESC 
        LIMIT 1
    `, region).Scan(&p.ID, &p.Region, &p.Price, &p.Unit, &p.Source, &p.RecordedAt, &p.CreatedAt)
    
    if err != nil {
        return "", fmt.Errorf("no price data found for region %s: %v", region, err)
    }
    
    jsonData, err := json.Marshal(p)
    if err != nil {
        return "", err
    }
    
    return string(jsonData), nil
}