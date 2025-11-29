# 🕷️ Web Scraping Guide - TobaccoTrack

## 📋 Overview

Aplikasi TobaccoTrack menggunakan **multi-source scraping strategy** dengan fallback mechanism untuk mendapatkan harga tembakau real-time.

---

## 🎯 Sumber Data (Prioritas)

### **1. BAPPEBTI Info Harga** (Primary)

**URL**: https://infoharga.bappebti.go.id/harga_komoditi_pedagang

**Komoditas Tersedia:**
- Tembakau Boyolali
- Tembakau Burley  
- Tembakau Kasturi

**Metode**: HTML Table Scraping

**Update Frequency**: 2x sehari (pagi & sore)

**Data Quality**: ✅ Official Government Data

**Kelebihan:**
- Data resmi dari BAPPEBTI
- Update rutin dari kontributor lapangan
- Multiple regions

**Kekurangan:**
- Tidak semua region tersedia
- Website bisa down
- Struktur HTML bisa berubah

---

### **2. Manual Research + Market Data** (Fallback) ⭐⭐⭐

**Sumber**: Riset manual dari:
- Portal berita (InfoPublik, ANTARA, iNews)
- Laporan pemerintah daerah
- Survey pasar langsung

**Data Base (Update Manual):**

| Region | Harga (Rp/kg) | Tanggal Riset | Sumber |
|--------|---------------|---------------|--------|
| Jember | 85,000 | 2024-09-15 | DPRD Jember Report |
| Temanggung | 150,000 | 2024-09-18 | InfoPublik + ANTARA |
| Lombok | 78,000 | 2024-08-01 | Market Survey |
| Klaten | 88,000 | 2024-07-15 | Regional Market |
| Pamekasan | 95,000 | 2024-08-20 | Madura Survey |

**Metode**: Base price + daily variation (±2%)

**Kelebihan:**
- Always available (offline mode)
- Berdasarkan riset real
- Reliable fallback

**Kekurangan:**
- Perlu update manual berkala
- Simulasi variation (bukan real-time exact)

---

## 🔧 Implementasi

### **Arsitektur Scraper**

```
ScraperManager
├── BAPPEBTIScraper (Try first)
│   ├── Success → Use data
│   └── Fail → Next scraper
│
└── MockScraperWithRealData (Fallback)
    └── Always success (offline data)
```

### **Flow Diagram**

```
User Click "Fetch Harga"
    ↓
Try BAPPEBTIScraper
    ↓
    ├── Success → Save to DB → Done ✅
    │
    └── Fail
        ↓
    Try MockScraperWithRealData
        ↓
        Save to DB → Done ✅
```

---

## 📝 Code Structure

### **File: `scraper.go`**

**Interfaces:**
```go
type TobaccoScraper interface {
    Scrape() ([]ScrapedPrice, error)
    GetName() string
}
```

**Implementations:**
- `BAPPEBTIScraper` - Scrape BAPPEBTI website
- `MockScraperWithRealData` - Fallback dengan data riset

**Manager:**
- `ScraperManager` - Koordinasi multiple scrapers

---

## 🚀 Usage

### **Install Dependencies**

```bash
cd backend
go get github.com/PuerkitoBio/goquery
go mod tidy
```

### **Run Server**

```bash
go run .
```

### **Test Scraping**

**Via cURL:**
```bash
# Fetch & save all prices
curl -X POST http://localhost:8080/harga/fetch

# Preview scraped data
curl http://localhost:8080/harga/current?region=Jember
```

**Via UI:**
Klik tombol "⚡ Fetch & Simpan Semua Region"

---

## 🔄 Update Manual Data

### **Cara Update Base Price:**

1. **Cari berita terbaru** tentang harga tembakau
2. **Edit file `scraper.go`**:

```go
LastResearch: map[string]PriceResearch{
    "Jember": {
        BasePrice:   90000,  // ← Update harga baru
        DateChecked: time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC),  // ← Update tanggal
        Source:      "Portal Berita X",
        Notes:       "Harga naik karena permintaan tinggi",
    },
    // ... dst
}
```

3. **Restart server**

---

## 🛡️ Error Handling

### **Skenario 1: BAPPEBTI Down**
```
✓ Fallback otomatis ke MockScraperWithRealData
✓ Log: "Scraping failed, fallback to simulation"
✓ Data tetap tersimpan ke database
```

### **Skenario 2: Network Error**
```
✓ Timeout setelah 10 detik
✓ Try next scraper
✓ Always have fallback data
```

### **Skenario 3: HTML Structure Changed**
```
✓ Parsing error → Skip to fallback
✓ Log error untuk debugging
✓ Manual fix needed (update selector)
```

---

## 📊 Data Quality Indicators

Setiap data yang tersimpan di database memiliki label source:

| Source Label | Meaning |
|--------------|---------|
| `BAPPEBTI Info Harga (Scraped: Standard)` | Real scraping dari BAPPEBTI |
| `InfoPublik + ANTARA News (Last checked: 2024-09-18)` | Manual research |
| `Market Survey (Scraped: Standard)` | Data dari survey pasar |

---

## ⚠️ Legal & Ethics

### **✅ Allowed:**
- Scraping public government websites (BAPPEBTI, BPS)
- Reasonable request frequency (< 1 request/minute)
- For educational/research purposes
- Proper User-Agent identification

### **❌ Not Allowed:**
- Aggressive scraping (DDoS-like behavior)
- Bypassing CAPTCHA/authentication
- Commercial resale of scraped data
- Ignoring robots.txt

### **robots.txt Check:**
```bash
curl https://infoharga.bappebti.go.id/robots.txt
```

---

## 🔮 Future Improvements

### **Short Term:**
- [ ] Add more news portal scrapers
- [ ] Implement caching (avoid duplicate scraping)
- [ ] Better HTML parsing (XPath)
- [ ] Retry logic with exponential backoff

### **Medium Term:**
- [ ] Headless browser (Selenium/Playwright) untuk JS-heavy sites
- [ ] Machine learning untuk extract harga dari free text
- [ ] API dari koperasi tembakau
- [ ] Crowdsourcing platform

### **Long Term:**
- [ ] Build internal marketplace (real trading data)
- [ ] Partnership dengan BAPPEBTI untuk official API
- [ ] Mobile app untuk petani submit harga
- [ ] Blockchain untuk data integrity

---

## 📞 Support

**Jika Scraper Fail:**
1. Check website accessibility: `curl -I https://infoharga.bappebti.go.id`
2. Check logs: `tail -f /var/log/tobaccotrack.log`
3. Verify HTML structure (might have changed)
4. Update selectors in `scraper.go`

**Contact:**
- GitHub Issues: [your-repo]/issues
- Email: [your-email]

---

## 📚 References

**Web Scraping Libraries:**
- goquery: https://github.com/PuerkitoBio/goquery
- colly: https://github.com/gocolly/colly

**Data Sources:**
- BAPPEBTI: https://infoharga.bappebti.go.id
- InfoPublik: https://infopublik.id
- ANTARA News: https://antaranews.com

**Best Practices:**
- Respectful Web Scraping: https://www.scrapehero.com/web-scraping-best-practices/
- robots.txt: https://developers.google.com/search/docs/crawling-indexing/robots/intro

---

## ✅ Checklist Deployment

- [ ] Install dependencies (`go get goquery`)
- [ ] Test scraping locally
- [ ] Set proper User-Agent
- [ ] Implement rate limiting
- [ ] Add error logging
- [ ] Setup monitoring
- [ ] Document data sources
- [ ] Get legal clearance (if needed)

---

**Last Updated**: November 30, 2024  
**Version**: 1.0.0  
**License**: Educational Use Only