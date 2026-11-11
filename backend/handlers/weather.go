package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"agritec/backend/config"

	"github.com/gin-gonic/gin"
)

var adm4Pattern = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{2}\.\d{4}$`)
var weatherClient = &http.Client{Timeout: 8 * time.Second}

// Weather: proxy prakiraan cuaca BMKG berdasarkan kode wilayah kelurahan (adm4).
func (h *H) Weather(c *gin.Context) {
	adm4 := c.DefaultQuery("adm4", config.GetEnv("BMKG_ADM4", "31.71.03.1001"))
	if !adm4Pattern.MatchString(adm4) {
		fail(c, http.StatusBadRequest, "Kode wilayah tidak valid, contoh: 31.71.03.1001")
		return
	}
	resp, err := weatherClient.Get("https://api.bmkg.go.id/publik/prakiraan-cuaca?adm4=" + adm4)
	if err != nil {
		fail(c, http.StatusBadGateway, "Gagal menghubungi BMKG")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fail(c, http.StatusBadGateway, fmt.Sprintf("BMKG merespons %d", resp.StatusCode))
		return
	}
	type lokasi struct {
		Desa, Kecamatan, Kotkab, Provinsi string
	}
	var raw struct {
		Lokasi lokasi `json:"lokasi"`
		Data   []struct {
			Lokasi lokasi `json:"lokasi"`
			Cuaca  [][]struct {
				LocalDatetime string  `json:"local_datetime"`
				T             float64 `json:"t"`
				Hu            float64 `json:"hu"`
				WeatherDesc   string  `json:"weather_desc"`
				Ws            float64 `json:"ws"`
				Image         string  `json:"image"`
			} `json:"cuaca"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		fail(c, http.StatusBadGateway, "Format data BMKG berubah")
		return
	}
	loc := raw.Lokasi
	type fc struct {
		Time     string  `json:"time"`
		Temp     float64 `json:"temp"`
		Humidity float64 `json:"humidity"`
		Desc     string  `json:"desc"`
		Wind     float64 `json:"wind"`
		Icon     string  `json:"icon"`
	}
	out := []fc{}
	if len(raw.Data) > 0 {
		if loc.Desa == "" {
			loc = raw.Data[0].Lokasi
		}
		for _, day := range raw.Data[0].Cuaca {
			for _, x := range day {
				out = append(out, fc{x.LocalDatetime, x.T, x.Hu, x.WeatherDesc, x.Ws, x.Image})
			}
		}
	}
	parts := []string{}
	for _, p := range []string{loc.Desa, loc.Kecamatan, loc.Kotkab, loc.Provinsi} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	c.Header("Cache-Control", "public, s-maxage=1800, stale-while-revalidate=3600")
	c.JSON(http.StatusOK, gin.H{"location": strings.Join(parts, ", "), "forecasts": out, "source": "BMKG"})
}
