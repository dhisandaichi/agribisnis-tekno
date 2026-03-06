package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IntegrationController struct {
	DB *gorm.DB
}

func NewIntegrationController(db *gorm.DB) *IntegrationController {
	return &IntegrationController{DB: db}
}

// GetWeather – Proxy ke API BMKG (Open Data)
// API BMKG tidak memerlukan API key untuk area-based forecast
func (ctrl *IntegrationController) GetWeather(c *gin.Context) {
	// Parameter: ?city=malang (default: jakarta)
	city := c.DefaultQuery("city", "jakarta")

	// BMKG Open Data - Prakiraan Cuaca per Kota
	bmkgURL := "https://api.bmkg.go.id/publik/prakiraan-cuaca?adm4=" + url.QueryEscape(city)

	client := &http.Client{}
	req, err := http.NewRequest("GET", bmkgURL, nil)
	if err != nil {
		// Fallback to mock data jika BMKG tidak tersedia
		c.JSON(http.StatusOK, gin.H{
			"source": "fallback",
			"data": gin.H{
				"city":           city,
				"temperature":    28,
				"condition":      "Cerah Berawan",
				"humidity":       72,
				"wind_speed_kmh": 14,
				"uv_index":       "Sedang",
				"rain_chance":    20,
			},
		})
		return
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		// Fallback jika BMKG tidak dapat dijangkau
		c.JSON(http.StatusOK, gin.H{
			"source": "fallback",
			"data": gin.H{
				"city":           city,
				"temperature":    28,
				"condition":      "Cerah Berawan",
				"humidity":       72,
				"wind_speed_kmh": 14,
				"uv_index":       "Sedang",
				"rain_chance":    20,
			},
		})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result interface{}
	json.Unmarshal(body, &result)

	c.JSON(http.StatusOK, gin.H{"source": "bmkg", "data": result})
}

// GetCommodityPrices – Data harga komoditas (mock/dapat disambungkan ke PIHPS/BPS API)
func (ctrl *IntegrationController) GetCommodityPrices(c *gin.Context) {
	// TODO: Connect to PIHPS Nasional API (https://hargapangan.id)
	// Untuk saat ini menggunakan data representatif
	c.JSON(http.StatusOK, gin.H{
		"source": "internal",
		"data": []gin.H{
			{"name": "Gabah Kering (GKP)", "price_per_kg": 6400, "unit": "Rp/kg", "trend": "up", "change_pct": 2.1},
			{"name": "Bawang Merah", "price_per_kg": 24500, "unit": "Rp/kg", "trend": "down", "change_pct": -1.4},
			{"name": "Cabai Merah Besar", "price_per_kg": 32000, "unit": "Rp/kg", "trend": "up", "change_pct": 5.2},
			{"name": "Jagung Pipilan", "price_per_kg": 5100, "unit": "Rp/kg", "trend": "up", "change_pct": 0.9},
			{"name": "Kedelai", "price_per_kg": 13200, "unit": "Rp/kg", "trend": "down", "change_pct": -0.7},
		},
	})
}
