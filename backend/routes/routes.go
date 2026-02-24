package routes

import (
	"github.com/agritec/backend/controllers"
	"github.com/agritec/backend/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // allow all origins
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	// Controllers
	authController := controllers.NewAuthController(db)
	visitController := controllers.NewVisitController(db)

	// API V1 Group
	v1 := r.Group("/api/v1")
	{
		// Authentication group route
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Sales group route (Protected by JWT)
		sales := v1.Group("/sales")
		sales.Use(middleware.JWTAuthMiddleware())
		sales.Use(middleware.RoleMiddleware("sales"))
		{
			sales.POST("/visits", visitController.CreateVisit) // create a visit
		}

		// Tambahkan rute untuk supervisor/farmer lainnya
	}

	return r
}
