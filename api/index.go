// Vercel Function: semua request /api/* diarahkan ke sini (lihat vercel.json).
package handler

import (
	"net/http"

	"agritec/backend"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	backend.Handler().ServeHTTP(w, r)
}
