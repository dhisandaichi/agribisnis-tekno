package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"agritec/backend/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type CurrentUser struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Active   bool   `json:"-"`
}

// Auth memverifikasi JWT lalu memuat user dari database (role selalu terbaru).
func Auth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Silakan login terlebih dahulu"})
			return
		}
		secret, err := config.JWTSecret()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Server belum dikonfigurasi"})
			return
		}
		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return secret, nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Sesi login habis, silakan login ulang"})
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		uid, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			return
		}
		var u CurrentUser
		res := db.Raw(`select id, email, coalesce(full_name,'') as full_name, coalesce(phone,'') as phone, role, active
			from users where id = ?`, int64(uid)).Scan(&u)
		if res.Error != nil || res.RowsAffected == 0 || !u.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Akun tidak ditemukan atau dinonaktifkan"})
			return
		}
		c.Set("user", &u)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := User(c)
		for _, r := range roles {
			if u != nil && u.Role == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Anda tidak punya izin untuk tindakan ini"})
	}
}

func User(c *gin.Context) *CurrentUser {
	v, ok := c.Get("user")
	if !ok {
		return nil
	}
	return v.(*CurrentUser)
}
