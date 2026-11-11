package handlers

import (
	"net/http"
	"time"

	"agritec/backend/middleware"

	"github.com/gin-gonic/gin"
)

type Price struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" binding:"required,max=255"`
	Price     float64   `json:"price" binding:"gte=0"`
	Unit      string    `json:"unit" binding:"required,max=20"`
	Trend     string    `json:"trend" binding:"required,oneof=up down flat"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Education struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title" binding:"required,max=255"`
	Description     *string   `json:"description"`
	VideoURL        *string   `json:"video_url" binding:"omitempty,url"`
	ThumbnailURL    *string   `json:"thumbnail_url" binding:"omitempty,url"`
	DurationMinutes *int      `json:"duration_minutes" binding:"omitempty,gte=0"`
	Published       bool      `json:"published"`
	CreatedAt       time.Time `json:"created_at"`
}

func (h *H) ListPrices(c *gin.Context) {
	var out []Price
	if err := h.DB.Raw(`select id, name, price::float8 as price, unit, trend, updated_at from commodity_prices order by name`).Scan(&out).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, rows(out))
}

func (h *H) ListEducation(c *gin.Context) {
	q := `select * from education_contents `
	if middleware.User(c).Role != "supervisor" {
		q += `where published `
	}
	var out []Education
	if err := h.DB.Raw(q + `order by created_at desc`).Scan(&out).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, rows(out))
}

func (h *H) SavePrice(c *gin.Context) {
	var in Price
	if !bindJSON(c, &in) {
		return
	}
	var err error
	if c.Param("id") == "" {
		err = h.DB.Exec(`insert into commodity_prices (name, price, unit, trend) values (?, ?, ?, ?)`, in.Name, in.Price, in.Unit, in.Trend).Error
	} else {
		id, ok := idParam(c)
		if !ok {
			return
		}
		err = h.DB.Exec(`update commodity_prices set name = ?, price = ?, unit = ?, trend = ?, updated_at = now() where id = ?`,
			in.Name, in.Price, in.Unit, in.Trend, id).Error
	}
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tersimpan"})
}

func (h *H) SaveEducation(c *gin.Context) {
	var in Education
	if !bindJSON(c, &in) {
		return
	}
	var err error
	if c.Param("id") == "" {
		err = h.DB.Exec(`insert into education_contents (title, description, video_url, thumbnail_url, duration_minutes, published) values (?, ?, ?, ?, ?, ?)`,
			in.Title, in.Description, in.VideoURL, in.ThumbnailURL, in.DurationMinutes, in.Published).Error
	} else {
		id, ok := idParam(c)
		if !ok {
			return
		}
		err = h.DB.Exec(`update education_contents set title = ?, description = ?, video_url = ?, thumbnail_url = ?, duration_minutes = ?, published = ? where id = ?`,
			in.Title, in.Description, in.VideoURL, in.ThumbnailURL, in.DurationMinutes, in.Published, id).Error
	}
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tersimpan"})
}

func (h *H) deleteFrom(table string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := idParam(c)
		if !ok {
			return
		}
		if err := h.DB.Exec(`delete from `+table+` where id = ?`, id).Error; err != nil {
			serverError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Dihapus"})
	}
}

func (h *H) DeletePrice() gin.HandlerFunc     { return h.deleteFrom("commodity_prices") }
func (h *H) DeleteEducation() gin.HandlerFunc { return h.deleteFrom("education_contents") }
