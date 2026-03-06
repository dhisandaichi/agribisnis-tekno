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

	// Root Route (Mencegah 404 ketika domain utama dibuka)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Welcome to Agritec API. System is up and running.",
		})
	})

	// Controllers
	authController := controllers.NewAuthController(db)
	visitController := controllers.NewVisitController(db)
	distributorController := controllers.NewDistributorController(db)
	educationController := controllers.NewEducationController(db)
	notificationController := controllers.NewNotificationController(db)
	integrationController := controllers.NewIntegrationController(db)

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
			sales.GET("/visits", visitController.GetVisits)    // list own visits
		}

		// Supervisor group route (Protected by JWT)
		supervisor := v1.Group("/supervisor")
		supervisor.Use(middleware.JWTAuthMiddleware())
		supervisor.Use(middleware.RoleMiddleware("supervisor"))
		{
			supervisor.GET("/visits", visitController.GetVisits) // list all visits
		}

		// Distributor Service (JWT required)
		distributor := v1.Group("/distributors")
		distributor.Use(middleware.JWTAuthMiddleware())
		{
			distributor.GET("/", distributorController.GetDistributors)
			distributor.GET("/:id", distributorController.GetDistributorByID)
			distributor.GET("/stock", distributorController.CheckStock) // ?distributor_id=X
			distributor.GET("/debts", distributorController.GetDebts)   // ?distributor_id=X
			distributor.POST("/debts", distributorController.CreateDebt)
		}

		// Education Service (GET: public, POST: supervisor only)
		education := v1.Group("/education")
		{
			education.GET("/content", educationController.GetContents) // ?category=X
			education.GET("/content/:id", educationController.GetContentByID)
		}
		educationAuth := v1.Group("/education")
		educationAuth.Use(middleware.JWTAuthMiddleware())
		educationAuth.Use(middleware.RoleMiddleware("supervisor"))
		{
			educationAuth.POST("/content", educationController.CreateContent)
		}

		// Notification Service (JWT required)
		notification := v1.Group("/notifications")
		notification.Use(middleware.JWTAuthMiddleware())
		{
			notification.GET("/", notificationController.GetNotifications)
			notification.POST("/", notificationController.CreateNotification)
			notification.PUT("/:id/read", notificationController.MarkAsRead)
		}

		// External Integration Service (public)
		integration := v1.Group("/integration")
		{
			integration.GET("/weather", integrationController.GetWeather) // ?city=malang
			integration.GET("/commodity-prices", integrationController.GetCommodityPrices)
		}
	}

	return r
}
