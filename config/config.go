package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

var (
    RazorpayKey    string
    RazorpaySecret string
)

func Init() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found, using system env vars")
    }

    RazorpayKey = os.Getenv("RAZORPAY_KEY")
    RazorpaySecret = os.Getenv("RAZORPAY_SECRET")

    log.Println("Razorpay config loaded")
}
