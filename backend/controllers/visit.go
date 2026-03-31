package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VisitController struct {
	DB *gorm.DB
}

func NewVisitController(db *gorm.DB) *VisitController {
	return &VisitController{DB: db}
}

// CreateVisit simulates adding a visit + image
func (h *VisitController) CreateVisit(c *gin.Context) {
	// Usually needs parsing form-data + image compression
	c.JSON(http.StatusCreated, gin.H{
		"message": "Visit created successfully. Image processing to be implemented.",
	})
}
