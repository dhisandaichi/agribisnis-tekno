package handlers

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"agritec/backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const maxPhotoBytes = 3 * 1024 * 1024

type visitInput struct {
	ClientRef         string     `json:"client_ref" binding:"required,max=64"`
	FarmerName        string     `json:"farmer_name" binding:"required,max=255"`
	FarmerPhone       string     `json:"farmer_phone" binding:"max=50"`
	Village           string     `json:"village" binding:"max=255"`
	CropType          string     `json:"crop_type" binding:"required,max=100"`
	LandArea          float64    `json:"land_area" binding:"gt=0"`
	FertilizerUsage   *float64   `json:"fertilizer_usage" binding:"omitempty,gte=0"`
	CompetitorProduct string     `json:"competitor_product" binding:"max=255"`
	Notes             string     `json:"notes"`
	Latitude          *float64   `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude         *float64   `json:"longitude" binding:"required,gte=-180,lte=180"`
	VisitTime         *time.Time `json:"visit_time"`
	Photo             string     `json:"photo" binding:"required"` // data URL base64 (image/jpeg)
}

func decodePhoto(dataURL string) ([]byte, string, error) {
	comma := strings.Index(dataURL, ",")
	if !strings.HasPrefix(dataURL, "data:image/") || comma < 0 {
		return nil, "", errors.New("format foto tidak valid")
	}
	meta := dataURL[5:comma] // image/jpeg;base64
	ct := strings.SplitN(meta, ";", 2)[0]
	if ct != "image/jpeg" && ct != "image/png" && ct != "image/webp" {
		return nil, "", errors.New("tipe foto harus JPEG/PNG/WEBP")
	}
	data, err := base64.StdEncoding.DecodeString(dataURL[comma+1:])
	if err != nil {
		return nil, "", errors.New("foto rusak")
	}
	if len(data) > maxPhotoBytes {
		return nil, "", errors.New("ukuran foto maksimal 3 MB")
	}
	return data, ct, nil
}

func (h *H) CreateVisit(c *gin.Context) {
	var in visitInput
	if !bindJSON(c, &in) {
		return
	}
	photo, ct, err := decodePhoto(in.Photo)
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	u := middleware.User(c)
	visitTime := time.Now()
	if in.VisitTime != nil && in.VisitTime.Before(time.Now().Add(5*time.Minute)) {
		visitTime = *in.VisitTime
	}

	// Idempoten: kunjungan offline yang dikirim ulang tidak menjadi duplikat
	var existing int64
	if h.DB.Raw(`select id from visits where client_ref = ?`, in.ClientRef).Scan(&existing); existing > 0 {
		c.JSON(http.StatusOK, gin.H{"id": existing, "duplicate": true})
		return
	}

	var visitID int64
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		var farmerID int64
		if err := tx.Raw(`insert into farmers (name, phone, village, crop_type, land_area, fertilizer_usage, competitor_product, latitude, longitude, created_by)
			values (?, nullif(?, ''), nullif(?, ''), ?, ?, ?, nullif(?, ''), ?, ?, ?) returning id`,
			strings.TrimSpace(in.FarmerName), in.FarmerPhone, in.Village, in.CropType, in.LandArea, in.FertilizerUsage,
			in.CompetitorProduct, in.Latitude, in.Longitude, u.ID).Scan(&farmerID).Error; err != nil {
			return err
		}
		if err := tx.Raw(`insert into visits (sales_id, farmer_id, has_photo, notes, latitude, longitude, visit_time, client_ref)
			values (?, ?, true, nullif(?, ''), ?, ?, ?, ?) returning id`,
			u.ID, farmerID, in.Notes, in.Latitude, in.Longitude, visitTime, in.ClientRef).Scan(&visitID).Error; err != nil {
			return err
		}
		return tx.Exec(`insert into visit_photos (visit_id, content_type, data) values (?, ?, ?)`, visitID, ct, photo).Error
	})
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": visitID})
}

type VisitRow struct {
	ID                int64     `json:"id"`
	VisitTime         time.Time `json:"visit_time"`
	Notes             *string   `json:"notes"`
	Latitude          *float64  `json:"latitude"`
	Longitude         *float64  `json:"longitude"`
	HasPhoto          bool      `json:"has_photo"`
	SalesName         string    `json:"sales_name"`
	FarmerName        string    `json:"farmer_name"`
	Village           *string   `json:"village"`
	CropType          *string   `json:"crop_type"`
	LandArea          *float64  `json:"land_area"`
	CompetitorProduct *string   `json:"competitor_product"`
}

const visitSelect = `select v.id, v.visit_time, v.notes, v.latitude, v.longitude, v.has_photo,
	coalesce(nullif(u.full_name, ''), u.email) as sales_name,
	f.name as farmer_name, f.village, f.crop_type, f.land_area::float8 as land_area, f.competitor_product
	from visits v join users u on u.id = v.sales_id join farmers f on f.id = v.farmer_id `

func (h *H) MyVisits(c *gin.Context) {
	var out []VisitRow
	if err := h.DB.Raw(visitSelect+`where v.sales_id = ? order by v.visit_time desc limit 50`, middleware.User(c).ID).Scan(&out).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, rows(out))
}

// VisitPhoto: sales pemilik kunjungan atau supervisor
func (h *H) VisitPhoto(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	u := middleware.User(c)
	var p struct {
		SalesID     int64
		ContentType string
		Data        []byte
	}
	res := h.DB.Raw(`select v.sales_id, p.content_type, p.data from visit_photos p join visits v on v.id = p.visit_id where p.visit_id = ?`, id).Scan(&p)
	if res.Error != nil {
		serverError(c, res.Error)
		return
	}
	if res.RowsAffected == 0 || (u.Role != "supervisor" && p.SalesID != u.ID) {
		fail(c, http.StatusNotFound, "Foto tidak ditemukan")
		return
	}
	c.Header("Cache-Control", "private, max-age=86400")
	c.Data(http.StatusOK, p.ContentType, p.Data)
}
