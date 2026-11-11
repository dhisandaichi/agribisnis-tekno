// Package backend menyediakan http.Handler yang dipakai oleh server lokal (backend/cmd/server)
// dan Vercel Function (api/index.go).
package backend

import (
	"log"
	"net/http"
	"sync"

	"agritec/backend/db"
	"agritec/backend/routes"

	"github.com/gin-gonic/gin"
)

var (
	once    sync.Once
	handler http.Handler
)

func Handler() http.Handler {
	once.Do(func() {
		gin.SetMode(gin.ReleaseMode)
		g, err := db.Get()
		if err != nil {
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte(`{"error":"Database tidak dapat dihubungi"}`))
			})
			return
		}
		// Buat tabel otomatis bila database masih kosong (deploy pertama)
		var ready bool
		if err := g.Raw(`select to_regclass('public.education_contents') is not null`).Scan(&ready).Error; err == nil && !ready {
			if err := db.Migrate(g); err != nil {
				log.Println("auto-migrate gagal:", err)
			} else {
				log.Println("auto-migrate: skema database dibuat")
			}
		}
		handler = routes.SetupRouter(g)
	})
	return handler
}
