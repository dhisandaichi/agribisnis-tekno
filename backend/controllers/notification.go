package controllers

import (
	"net/http"

	"github.com/agritec/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificationController struct {
	DB *gorm.DB
}

func NewNotificationController(db *gorm.DB) *NotificationController {
	return &NotificationController{DB: db}
}

// GetNotifications – Mengambil notifikasi milik user yang sedang login
func (ctrl *NotificationController) GetNotifications(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDVal.(float64))

	var notifications []models.Notification
	if err := ctrl.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil notifikasi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": notifications, "total": len(notifications)})
}

// CreateNotification – Membuat notifikasi baru (internal / dipanggil oleh service lain)
func (ctrl *NotificationController) CreateNotification(c *gin.Context) {
	var input struct {
		UserID  uint   `json:"user_id" binding:"required"`
		Type    string `json:"type" binding:"required"` // "SMS" | "IN_APP"
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notif := models.Notification{
		UserID:  input.UserID,
		Type:    input.Type,
		Message: input.Message,
		Status:  "pending",
	}
	if err := ctrl.DB.Create(&notif).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat notifikasi"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Notifikasi berhasil dibuat", "data": notif})
}

// MarkAsRead – Menandai notifikasi sebagai sudah dibaca
func (ctrl *NotificationController) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Model(&models.Notification{}).Where("id = ?", id).Update("status", "sent").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status notifikasi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi ditandai sudah dibaca"})
}
