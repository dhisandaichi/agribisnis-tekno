package controllers

import (
	"net/http"
	"strconv"

	"github.com/agritec/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DistributorController struct {
	DB *gorm.DB
}

func NewDistributorController(db *gorm.DB) *DistributorController {
	return &DistributorController{DB: db}
}

// GetDistributors – Mengembalikan semua distributor dari database
func (ctrl *DistributorController) GetDistributors(c *gin.Context) {
	var distributors []models.Distributor
	if err := ctrl.DB.Find(&distributors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data distributor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": distributors, "total": len(distributors)})
}

// GetDistributorByID – Detail satu distributor
func (ctrl *DistributorController) GetDistributorByID(c *gin.Context) {
	id := c.Param("id")
	var distributor models.Distributor
	if err := ctrl.DB.First(&distributor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Distributor tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": distributor})
}

// CheckStock – Mengembalikan stok milik distributor tertentu
func (ctrl *DistributorController) CheckStock(c *gin.Context) {
	distIDStr := c.Query("distributor_id")
	var stocks []models.Stock
	query := ctrl.DB
	if distIDStr != "" {
		distID, err := strconv.Atoi(distIDStr)
		if err == nil {
			query = query.Where("distributor_id = ?", distID)
		}
	}
	if err := query.Find(&stocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data stok"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stocks, "total": len(stocks)})
}

// CreateDebt – Mencatat hutang baru untuk distributor
func (ctrl *DistributorController) CreateDebt(c *gin.Context) {
	var input struct {
		DistributorID uint    `json:"distributor_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
		DueDate       string  `json:"due_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	debt := models.Debt{
		DistributorID: input.DistributorID,
		Amount:        input.Amount,
		Status:        "unpaid",
	}
	if err := ctrl.DB.Create(&debt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat data hutang"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Hutang berhasil dicatat", "data": debt})
}

// GetDebts – Mengembalikan semua hutang, bisa filter by distributo_id
func (ctrl *DistributorController) GetDebts(c *gin.Context) {
	var debts []models.Debt
	distIDStr := c.Query("distributor_id")
	query := ctrl.DB
	if distIDStr != "" {
		distID, err := strconv.Atoi(distIDStr)
		if err == nil {
			query = query.Where("distributor_id = ?", distID)
		}
	}
	if err := query.Find(&debts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data hutang"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": debts})
}
