package main

import (
	"go-backend/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	router := gin.New()				// Create a gin router with default middleware
	router.Use(gin.Logger())	// Log all requests to the console
	router.Use(cors.Default()) 			// Enable CORS for all origins

	// User routes
	router.POST("/register", routes.CreateAccount)
	router.GET("/account/:id", routes.GetAccount)
	router.PUT("account/login", routes.LoginAccount)
	router.PUT("/account/update/:id", routes.UpdateAccount)
	router.DELETE("/account/delete/:id", routes.DeleteAccount)

	// Booking routes
	router.POST("/order", routes.CreateOrder)
	router.GET("/orders", routes.GetOrders)
	router.DELETE("/order/:id", routes.DeleteOrder)

	// Event routes
	router.POST("/event", routes.CreateEvent)
	router.GET("/events", routes.GetEvents)
	router.DELETE("/event/:id", routes.DeleteEvent)
	

	router.Run(":" + port)
}