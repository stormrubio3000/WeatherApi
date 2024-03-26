package main

import (
	"weather-api/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	//Create gin router with single endpooint needed and then run on local host
	router := gin.Default()
	router.GET("/", controller.ShowWeather)

	router.Run("localhost:8080")
}
