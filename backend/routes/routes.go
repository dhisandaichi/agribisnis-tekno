package routes

import (
	"net/http"
	"strings"

	"agritec/backend/config"
	"agritec/backend/handlers"
	mw "agritec/backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CORS hanya dibutuhkan saat frontend & API beda domain (mis. dev lokal).
// Atur ALLOWED_ORIGINS="http://localhost:5173,https://domain-lain.com"
func cors() gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, o := range strings.Split(config.GetEnv("ALLOWED_ORIGINS", "http://localhost:5173"), ",") {
		if o = strings.TrimSpace(o); o != "" {
			allowed[o] = true
		}
	}
	return func(c *gin.Context) {
		if origin := c.GetHeader("Origin"); allowed[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), cors())
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4<<20) // batas 4 MB per request
		c.Header("Cache-Control", "no-store")
		c.Next()
	})

	h := &handlers.H{DB: db}
	v1 := r.Group("/api/v1")

	v1.GET("/health", func(c *gin.Context) {
		if err := db.Exec("select 1").Error; err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db_error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	v1.GET("/weather", h.Weather)

	auth := v1.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	// Semua route di bawah ini butuh login
	p := v1.Group("", mw.Auth(db))
	p.GET("/me", h.Me)
	p.PUT("/me/password", h.ChangePassword)
	p.GET("/content/prices", h.ListPrices)
	p.GET("/content/education", h.ListEducation)
	p.GET("/visits/:id/photo", mw.RequireRole("sales", "supervisor"), h.VisitPhoto)

	sales := p.Group("/sales", mw.RequireRole("sales"))
	sales.POST("/visits", h.CreateVisit)
	sales.GET("/visits", h.MyVisits)

	dist := p.Group("/distributor", mw.RequireRole("distributor"))
	dist.GET("", h.DistributorOverview)
	dist.POST("/orders", h.DistributorCreateOrder)

	sup := p.Group("/supervisor", mw.RequireRole("supervisor"))
	sup.GET("/summary", h.Summary)
	sup.GET("/visits", h.AllVisits)
	sup.GET("/distributors", h.ListDistributors)
	sup.POST("/distributors", h.SaveDistributor)
	sup.PUT("/distributors/:id", h.SaveDistributor)
	sup.GET("/orders", h.ListOrders)
	sup.POST("/orders", h.SaveOrder)
	sup.PUT("/orders/:id", h.SaveOrder)
	sup.PATCH("/orders/:id/status", h.SetOrderStatus)
	sup.DELETE("/orders/:id", h.DeleteOrder)
	sup.GET("/users", h.ListUsers)
	sup.PATCH("/users/:id", h.UpdateUser)
	sup.POST("/users/:id/reset-password", h.ResetUserPassword)
	sup.POST("/prices", h.SavePrice)
	sup.PUT("/prices/:id", h.SavePrice)
	sup.DELETE("/prices/:id", h.DeletePrice())
	sup.POST("/education", h.SaveEducation)
	sup.PUT("/education/:id", h.SaveEducation)
	sup.DELETE("/education/:id", h.DeleteEducation())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint tidak ditemukan"})
	})
	return r
}
