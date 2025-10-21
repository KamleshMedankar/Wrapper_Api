package routes

import (
	"payment_wrapper/controllers"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(r *gin.Engine) {
	r.POST("/create-payment", controllers.CreatePayment)
	r.POST("/verify-payment", controllers.VerifyPayment)
	r.POST("/webhook", controllers.PaymentWebhook)
}
