package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"agritec/backend"
	"agritec/backend/config"
	"agritec/backend/db"

	"github.com/joho/godotenv"
)

// Server mandiri (dev lokal & VPS/Docker).
//   go run ./backend/cmd/server            -> jalankan API di :8080
//   go run ./backend/cmd/server -migrate   -> jalankan skema database lalu keluar
// Bila STATIC_DIR diisi (mis. frontend/dist), server juga melayani frontend PWA.
func main() {
	migrate := flag.Bool("migrate", false, "jalankan backend/db/schema.sql lalu keluar")
	flag.Parse()
	_ = godotenv.Load("backend/.env")
	_ = godotenv.Load(".env")

	if *migrate {
		g, err := db.Get()
		if err != nil {
			log.Fatal(err)
		}
		if err := db.Migrate(g); err != nil {
			log.Fatal("migrasi gagal: ", err)
		}
		log.Println("migrasi selesai")
		return
	}
	if _, err := config.JWTSecret(); err != nil {
		log.Fatal(err)
	}

	api := backend.Handler()
	var handler http.Handler = api
	if dir := os.Getenv("STATIC_DIR"); dir != "" {
		handler = withStatic(api, dir)
		log.Printf("melayani frontend dari %s", dir)
	}

	port := config.GetEnv("PORT", "8080")
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           securityHeaders(handler),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	log.Printf("server berjalan di http://localhost:%s (API: /api/v1)", port)
	log.Fatal(srv.ListenAndServe())
}

// withStatic: /api/* ke backend, sisanya file statis dengan fallback SPA ke index.html.
func withStatic(api http.Handler, dir string) http.Handler {
	fs := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			api.ServeHTTP(w, r)
			return
		}
		p := filepath.Join(dir, filepath.Clean("/"+r.URL.Path))
		if info, err := os.Stat(p); err != nil || info.IsDir() {
			w.Header().Set("Cache-Control", "no-cache")
			http.ServeFile(w, r, filepath.Join(dir, "index.html"))
			return
		}
		switch {
		case strings.HasPrefix(r.URL.Path, "/assets/"):
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		case r.URL.Path == "/sw.js" || strings.HasSuffix(r.URL.Path, ".webmanifest") || r.URL.Path == "/index.html":
			w.Header().Set("Cache-Control", "no-cache")
		}
		fs.ServeHTTP(w, r)
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "camera=(self), geolocation=(self), microphone=()")
		next.ServeHTTP(w, r)
	})
}
