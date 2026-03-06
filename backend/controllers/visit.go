package controllers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/agritec/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VisitController struct {
	DB *gorm.DB
}

func NewVisitController(db *gorm.DB) *VisitController {
	return &VisitController{DB: db}
}

// CreateVisit – Handles multipart/form-data with optional photo upload and GPS coords
func (h *VisitController) CreateVisit(c *gin.Context) {
	// Get userID from JWT context (set by middleware)
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	salesID := uint(userIDVal.(float64))

	// Parse form fields
	farmerName := c.PostForm("farmer_name")
	cropType := c.PostForm("crop_type")
	notes := c.PostForm("notes")
	latStr := c.PostForm("latitude")
	lngStr := c.PostForm("longitude")

	if farmerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "farmer_name is required"})
		return
	}

	// Find or create farmer
	var farmer models.Farmer
	result := h.DB.Where("name = ?", farmerName).First(&farmer)
	if result.Error != nil {
		// Parse lat/lng
		var lat, lng float64
		if latStr != "" {
			lat, _ = strconv.ParseFloat(latStr, 64)
			lng, _ = strconv.ParseFloat(lngStr, 64)
		}
		farmer = models.Farmer{
			Name:      farmerName,
			Latitude:  lat,
			Longitude: lng,
			CropType:  cropType,
		}
		if err := h.DB.Create(&farmer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create farmer"})
			return
		}
	}

	// Handle optional photo upload
	photoPath := ""
	file, header, err := c.Request.FormFile("photo")
	if err == nil && header != nil {
		defer file.Close()
		ext := strings.ToLower(filepath.Ext(header.Filename))
		allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, PNG, WEBP allowed."})
			return
		}
		// Read file bytes (ready to upload to S3 / Supabase Storage)
		_, err = io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read photo"})
			return
		}
		// TODO: Upload to Supabase Storage / S3 and get URL
		// For now, store filename as placeholder
		photoPath = "uploads/" + header.Filename
	}

	// Create visit record
	visit := models.Visit{
		SalesID:   salesID,
		FarmerID:  farmer.ID,
		PhotoPath: photoPath,
		VisitTime: time.Now(),
		Notes:     notes,
	}

	if err := h.DB.Create(&visit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create visit"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Visit recorded successfully",
		"visit": gin.H{
			"id":         visit.ID,
			"farmer":     farmerName,
			"crop_type":  cropType,
			"photo_path": photoPath,
			"visit_time": visit.VisitTime,
		},
	})
}

// GetVisits – List all visits for the authenticated sales (or all for supervisor)
func (h *VisitController) GetVisits(c *gin.Context) {
	var visits []models.Visit
	userRole, _ := c.Get("userRole")
	userIDVal, _ := c.Get("userID")

	query := h.DB.Preload("Farmer")
	if userRole == "sales" {
		salesID := uint(userIDVal.(float64))
		query = query.Where("sales_id = ?", salesID)
	}

	if err := query.Order("visit_time desc").Find(&visits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch visits"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": visits, "total": len(visits)})
}
