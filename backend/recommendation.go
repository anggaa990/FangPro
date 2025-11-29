package main

import (
    "strings"
)

type RecommendationResult struct {
    Status           string   `json:"status"`            // "optimal", "good", "caution", "not_recommended"
    MainAdvice       string   `json:"main_advice"`
    DetailedAdvice   []string `json:"detailed_advice"`
    PlantingAdvice   string   `json:"planting_advice"`
    HarvestAdvice    string   `json:"harvest_advice"`
    DryingAdvice     string   `json:"drying_advice"`
    PestWarning      string   `json:"pest_warning"`
    IrrigationAdvice string   `json:"irrigation_advice"`
    Temperature      float64  `json:"temperature"`
    Humidity         int      `json:"humidity"`
    RainMM           float64  `json:"rain_mm"`
    Region           string   `json:"region"`
}

// Recommend memberikan rekomendasi berdasarkan data cuaca
func Recommend(temp float64, humidity int, rain float64) string {
    var recommendations []string

    // Analisis Suhu
    if temp >= 20 && temp <= 30 {
        recommendations = append(recommendations, "✅ Suhu optimal untuk pertumbuhan tembakau (20-30°C)")
    } else if temp < 20 {
        recommendations = append(recommendations, "⚠️ Suhu terlalu dingin, pertumbuhan mungkin terhambat")
    } else {
        recommendations = append(recommendations, "⚠️ Suhu terlalu panas, tingkatkan irigasi")
    }

    // Analisis Kelembaban
    if humidity >= 60 && humidity <= 80 {
        recommendations = append(recommendations, "✅ Kelembaban ideal untuk tembakau (60-80%)")
    } else if humidity < 60 {
        recommendations = append(recommendations, "⚠️ Kelembaban rendah, tingkatkan irigasi")
    } else {
        recommendations = append(recommendations, "⚠️ Kelembaban tinggi, risiko penyakit jamur meningkat")
    }

    // Analisis Curah Hujan
    if rain < 1 {
        recommendations = append(recommendations, "☀️ Cuaca kering, cocok untuk pengeringan daun tembakau")
    } else if rain >= 1 && rain < 5 {
        recommendations = append(recommendations, "🌦️ Hujan ringan, cocok untuk pertumbuhan")
    } else if rain >= 5 && rain < 10 {
        recommendations = append(recommendations, "🌧️ Hujan sedang, pastikan drainase baik")
    } else {
        recommendations = append(recommendations, "⛈️ Hujan lebat, tunda pemanenan, risiko busuk tinggi")
    }

    return strings.Join(recommendations, " | ")
}

// GetAdvancedRecommendation memberikan rekomendasi detail
func GetAdvancedRecommendation(temp float64, humidity int, rain float64, region string) RecommendationResult {
    result := RecommendationResult{
        Temperature: temp,
        Humidity:    humidity,
        RainMM:      rain,
        Region:      region,
    }

    var advice []string
    
    // Determine overall status
    optimalTemp := temp >= 20 && temp <= 30
    optimalHumidity := humidity >= 60 && humidity <= 80
    optimalRain := rain >= 1 && rain < 5

    if optimalTemp && optimalHumidity && optimalRain {
        result.Status = "optimal"
        result.MainAdvice = "🌟 Kondisi OPTIMAL untuk budidaya tembakau!"
    } else if optimalTemp || optimalHumidity {
        result.Status = "good"
        result.MainAdvice = "✅ Kondisi BAIK untuk budidaya tembakau"
    } else if temp > 35 || humidity > 90 || rain > 15 {
        result.Status = "not_recommended"
        result.MainAdvice = "❌ Kondisi TIDAK DISARANKAN untuk aktivitas pertanian"
    } else {
        result.Status = "caution"
        result.MainAdvice = "⚠️ Kondisi CUKUP - perhatikan faktor risiko"
    }

    // Temperature Analysis
    if temp < 15 {
        advice = append(advice, "Suhu terlalu dingin (<15°C) - pertumbuhan sangat terhambat")
        result.PlantingAdvice = "❌ TIDAK disarankan menanam. Tunggu suhu naik minimal 18°C"
    } else if temp >= 15 && temp < 20 {
        advice = append(advice, "Suhu sejuk (15-20°C) - pertumbuhan lambat")
        result.PlantingAdvice = "⚠️ Penanaman dimungkinkan tapi pertumbuhan akan lambat"
    } else if temp >= 20 && temp <= 30 {
        advice = append(advice, "Suhu optimal (20-30°C) - pertumbuhan ideal")
        result.PlantingAdvice = "✅ SANGAT COCOK untuk penanaman bibit baru"
    } else if temp > 30 && temp <= 35 {
        advice = append(advice, "Suhu hangat (30-35°C) - perlu irigasi ekstra")
        result.PlantingAdvice = "⚠️ Bisa menanam tapi pastikan irigasi mencukupi"
    } else {
        advice = append(advice, "Suhu sangat panas (>35°C) - stres tanaman tinggi")
        result.PlantingAdvice = "❌ TIDAK disarankan menanam. Tanaman akan stres"
    }

    // Humidity Analysis
    if humidity < 40 {
        advice = append(advice, "Kelembaban sangat rendah (<40%) - tanaman bisa layu")
        result.IrrigationAdvice = "💧 PENTING: Tingkatkan irigasi 2-3x sehari, gunakan mulsa"
    } else if humidity >= 40 && humidity < 60 {
        advice = append(advice, "Kelembaban rendah (40-60%) - perlu irigasi rutin")
        result.IrrigationAdvice = "💧 Irigasi 1-2x sehari, pantau kondisi tanah"
    } else if humidity >= 60 && humidity <= 80 {
        advice = append(advice, "Kelembaban ideal (60-80%) - kondisi sempurna")
        result.IrrigationAdvice = "✅ Irigasi normal sesuai jadwal standar"
    } else if humidity > 80 && humidity <= 90 {
        advice = append(advice, "Kelembaban tinggi (80-90%) - risiko penyakit jamur")
        result.IrrigationAdvice = "⚠️ Kurangi irigasi, pastikan drainase baik"
        result.PestWarning = "⚠️ PERINGATAN: Risiko penyakit jamur tinggi! Semprot fungisida preventif, tingkatkan sirkulasi udara"
    } else {
        advice = append(advice, "Kelembaban sangat tinggi (>90%) - bahaya penyakit")
        result.IrrigationAdvice = "❌ STOP irigasi, perbaiki drainase segera"
        result.PestWarning = "🚨 BAHAYA: Risiko penyakit jamur sangat tinggi! Aplikasi fungisida darurat, cek tanaman busuk"
    }

    // Rain Analysis
    if rain < 0.5 {
        advice = append(advice, "Cuaca kering - ideal untuk pengeringan")
        result.HarvestAdvice = "✅ SANGAT COCOK untuk panen dan pengeringan daun"
        result.DryingAdvice = "☀️ Kondisi SEMPURNA untuk penjemuran tembakau. Maksimalkan pengeringan hari ini!"
    } else if rain >= 0.5 && rain < 2 {
        advice = append(advice, "Hujan ringan - aman untuk pertumbuhan")
        result.HarvestAdvice = "✅ Bisa panen pagi hari sebelum hujan"
        result.DryingAdvice = "⚠️ Penjemuran bisa dilakukan dengan pengawasan ketat"
    } else if rain >= 2 && rain < 5 {
        advice = append(advice, "Hujan sedang - baik untuk vegetatif")
        result.HarvestAdvice = "⚠️ Tunda panen jika memungkinkan, atau panen cepat sebelum hujan lebat"
        result.DryingAdvice = "❌ Tidak disarankan menjemur hari ini. Gunakan pengering mekanis jika mendesak"
    } else if rain >= 5 && rain < 10 {
        advice = append(advice, "Hujan lebat - pastikan drainase baik")
        result.HarvestAdvice = "❌ TUNDA panen! Daun basah tidak layak dipanen"
        result.DryingAdvice = "❌ STOP penjemuran. Pindahkan tembakau ke tempat kering"
    } else {
        advice = append(advice, "Hujan sangat lebat - risiko genangan")
        result.HarvestAdvice = "❌ JANGAN panen. Cek kondisi tanaman setelah hujan reda"
        result.DryingAdvice = "❌ Penjemuran tidak memungkinkan. Pastikan gudang kering dan ventilasi baik"
        if result.PestWarning == "" {
            result.PestWarning = "⚠️ Cek tanaman setelah hujan reda - risiko busuk batang dan akar tinggi"
        }
    }

    // Combined Analysis for Harvesting
    if temp >= 25 && temp <= 32 && rain < 1 && humidity < 75 {
        result.HarvestAdvice = "🌟 KONDISI PANEN SEMPURNA! Suhu, kelembaban, dan cuaca mendukung"
    }

    // Pest and Disease Warnings
    if humidity > 80 && temp > 25 {
        if result.PestWarning == "" {
            result.PestWarning = "🚨 Kombinasi panas + lembab: Risiko tinggi embun tepung, busuk daun, dan serangan ulat"
        }
    } else if temp < 18 && rain > 5 {
        if result.PestWarning == "" {
            result.PestWarning = "⚠️ Kondisi dingin + basah: Waspadai penyakit busuk akar dan batang"
        }
    }

    // Default messages if not set
    if result.PlantingAdvice == "" {
        result.PlantingAdvice = "Evaluasi kondisi lebih lanjut sebelum penanaman"
    }
    if result.HarvestAdvice == "" {
        result.HarvestAdvice = "Pantau perkembangan cuaca untuk menentukan waktu panen"
    }
    if result.DryingAdvice == "" {
        result.DryingAdvice = "Sesuaikan metode pengeringan dengan kondisi cuaca"
    }
    if result.IrrigationAdvice == "" {
        result.IrrigationAdvice = "Lakukan irigasi sesuai kebutuhan tanaman"
    }
    if result.PestWarning == "" {
        result.PestWarning = "✅ Risiko hama dan penyakit dalam batas normal. Lakukan monitoring rutin"
    }

    result.DetailedAdvice = advice

    return result
}

// GetRecommendationSummary untuk backward compatibility
func GetRecommendationSummary(temp float64, humidity int, rain float64) string {
    return Recommend(temp, humidity, rain)
}