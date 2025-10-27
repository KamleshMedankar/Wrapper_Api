package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"payment_wrapper/db"
	"payment_wrapper/models"
	"payment_wrapper/utils"
	"strings"
	"time"
)

func CreatePayment(c *gin.Context) {
	var req models.Payment
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.OrderID = fmt.Sprintf("order_%d", time.Now().UnixNano()/1e6)

	// 1️⃣ Create Razorpay order
	paymentURL, err := utils.CreateRazorpayPayment(req.Amount, req.Currency, req.OrderID, req.CustomerEmail, req.CustomerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2️⃣ Extract Razorpay order ID from response URL
	req.RazorpayOrderID = extractOrderIDFromURL(paymentURL)
	req.Status = "created"

	// 3️⃣ Insert payment into DB immediately
	if err := db.InsertPayment(&req); err != nil {
		fmt.Println("DB Insert Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
        "error": fmt.Sprintf("DB insert failed: %v", err),
    })
		return
	}

	fmt.Println("Payment inserted successfully:", req.OrderID, req.RazorpayOrderID)

	// 4️⃣ Return payment URL to frontend
	c.JSON(http.StatusOK, gin.H{"payment_url": paymentURL})
}


func extractOrderIDFromURL(url string) string {
	parts := strings.Split(url, "order_id=")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}


// verifying payment
func VerifyPayment(c *gin.Context) {
	var req struct {
		PaymentID string `json:"payment_id"`
		OrderID   string `json:"order_id"`
		Signature string `json:"signature"`
		Gateway   string `json:"gateway"`
	}

	// Step 1: Bind JSON and log it
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("BindJSON Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("VerifyPayment Request Received:", req) // ✅ Log the incoming request

	status := "failed"

	// Step 2: Verify Razorpay payment
	if req.Gateway == "razorpay" {
		s, err := utils.VerifyRazorpayPayment(req.PaymentID, req.OrderID, req.Signature)
		if err != nil {
			fmt.Println("Razorpay Verification Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		status = s
		fmt.Println("Razorpay Verification Status:", status)
	} else {
		fmt.Println("Invalid Gateway:", req.Gateway)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gateway"})
		return
	}

	// Step 3: Update payment in DB and log results
	err := db.UpdatePaymentStatus(req.OrderID, status, req.PaymentID, req.Signature)
	if err != nil {
		fmt.Println("DB Update Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update payment status in db"})
		return
	}

	fmt.Println("Payment status updated successfully for order_id:", req.OrderID)
	c.JSON(http.StatusOK, gin.H{"status": status})
}

// Webhook handler
func PaymentWebhook(c *gin.Context) {
	// You can validate webhook signatures here
	c.JSON(http.StatusOK, gin.H{"message": "Webhook received"})
}
