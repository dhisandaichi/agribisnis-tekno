package api

import (
	"net/http"

	"github.com/agritec/backend/config"
	"github.com/agritec/backend/routes"
	"github.com/gin-gonic/gin"
)

var app *gin.Engine

func init() {
	config.InitConfig()
	db := config.ConnectDB()
	app = routes.SetupRouter(db)
}

// Handler adalah entry point untuk Vercel Serverless Function
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
