package config

import (
	"errors"
	"os"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

// JWTSecret wajib diisi minimal 32 karakter. Tidak ada nilai cadangan.
func JWTSecret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if len(s) < 32 {
		return nil, errors.New("JWT_SECRET belum diatur atau kurang dari 32 karakter")
	}
	return []byte(s), nil
}
