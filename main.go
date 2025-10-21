package main

import (
    "payment_wrapper/routes"
	"payment_wrapper/config"
    "github.com/gin-gonic/gin"
	"payment_wrapper/db"
)

func main() {
	config.Init()
	db.Connect()
    r := gin.Default()
	r.Static("/static", "./frontend")
	r.GET("/", func(c *gin.Context) {
    c.File("./frontend/index.html")
})
	r.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
    c.Next()
})


    // Routes
    routes.PaymentRoutes(r)

    // Run server
    r.Run(":8099") // localhost:8080
}
