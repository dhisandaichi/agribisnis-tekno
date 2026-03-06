package controllers

import (
	"net/http"

	"github.com/agritec/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EducationController struct {
	DB *gorm.DB
}

func NewEducationController(db *gorm.DB) *EducationController {
	return &EducationController{DB: db}
}

// GetContents – Mengembalikan semua konten edukasi, bisa filter by category
func (ctrl *EducationController) GetContents(c *gin.Context) {
	category := c.Query("category")
	var contents []models.Content
	query := ctrl.DB

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Order("created_at desc").Find(&contents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil konten"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": contents, "total": len(contents)})
}

// GetContentByID – Detail satu konten edukasi
func (ctrl *EducationController) GetContentByID(c *gin.Context) {
	id := c.Param("id")
	var content models.Content
	if err := ctrl.DB.First(&content, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Konten tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": content})
}

// CreateContent – Menambah konten edukasi baru (khusus supervisor)
func (ctrl *EducationController) CreateContent(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		VideoURL    string `json:"video_url"`
		Thumbnail   string `json:"thumbnail"`
		Category    string `json:"category"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content := models.Content{
		Title:       input.Title,
		Description: input.Description,
		VideoURL:    input.VideoURL,
		Thumbnail:   input.Thumbnail,
		Category:    input.Category,
	}
	if err := ctrl.DB.Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat konten"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Konten berhasil ditambahkan", "data": content})
}
