package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"agritec/backend/config"
	"agritec/backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type registerInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Role     string `json:"role" binding:"required,oneof=sales farmer distributor"`
	FullName string `json:"full_name" binding:"required,max=255"`
	Phone    string `json:"phone" binding:"max=50"`
}

func (h *H) Register(c *gin.Context) {
	var in registerInput
	if !bindJSON(c, &in) {
		return
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		serverError(c, err)
		return
	}
	var id int64
	res := h.DB.Raw(`insert into users (email, password_hash, full_name, phone, role)
		values (?, ?, ?, nullif(?, ''), ?) on conflict (email) do nothing returning id`,
		in.Email, string(hash), strings.TrimSpace(in.FullName), strings.TrimSpace(in.Phone), in.Role).Scan(&id)
	if res.Error != nil {
		serverError(c, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		fail(c, http.StatusConflict, "Email sudah terdaftar. Silakan login.")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Pendaftaran berhasil", "user_id": id})
}

type loginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *H) Login(c *gin.Context) {
	var in loginInput
	if !bindJSON(c, &in) {
		return
	}
	var u struct {
		middleware.CurrentUser
		PasswordHash string
	}
	res := h.DB.Raw(`select id, email, coalesce(full_name,'') as full_name, coalesce(phone,'') as phone, role, active, password_hash
		from users where email = ?`, strings.ToLower(strings.TrimSpace(in.Email))).Scan(&u)
	if res.Error != nil {
		serverError(c, res.Error)
		return
	}
	if res.RowsAffected == 0 || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
		fail(c, http.StatusUnauthorized, "Email atau password salah.")
		return
	}
	if !u.Active {
		fail(c, http.StatusForbidden, "Akun Anda dinonaktifkan. Hubungi supervisor.")
		return
	}
	secret, err := config.JWTSecret()
	if err != nil {
		serverError(c, err)
		return
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  u.ID,
		"role": u.Role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 hari, agar sales di lapangan tidak sering login ulang
	}).SignedString(secret)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "user": u.CurrentUser})
}

func (h *H) Me(c *gin.Context) {
	c.JSON(http.StatusOK, middleware.User(c))
}

type passwordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72"`
}

func (h *H) ChangePassword(c *gin.Context) {
	var in passwordInput
	if !bindJSON(c, &in) {
		return
	}
	u := middleware.User(c)
	var hash string
	h.DB.Raw(`select password_hash from users where id = ?`, u.ID).Scan(&hash)
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.OldPassword)) != nil {
		fail(c, http.StatusBadRequest, "Password lama salah.")
		return
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		serverError(c, err)
		return
	}
	if err := h.DB.Exec(`update users set password_hash = ? where id = ?`, string(newHash), u.ID).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}

func randomPassword() string {
	b := make([]byte, 9)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
