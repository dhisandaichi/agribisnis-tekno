package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type H struct{ DB *gorm.DB }

func fail(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{"error": msg})
}

// serverError mencatat detail di log, tapi tidak membocorkannya ke user.
func serverError(c *gin.Context, err error) {
	log.Printf("%s %s: %v", c.Request.Method, c.Request.URL.Path, err)
	fail(c, http.StatusInternalServerError, "Terjadi kesalahan di server. Coba lagi.")
}

func idParam(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		fail(c, http.StatusBadRequest, "ID tidak valid")
		return 0, false
	}
	return id, true
}

func bindJSON(c *gin.Context, v interface{}) bool {
	if err := c.ShouldBindJSON(v); err != nil {
		fail(c, http.StatusBadRequest, "Data tidak lengkap atau tidak valid: "+err.Error())
		return false
	}
	return true
}

// rows memastikan JSON mengembalikan [] bukan null
func rows[T any](v []T) []T {
	if v == nil {
		return []T{}
	}
	return v
}
