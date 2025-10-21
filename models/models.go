package models
import(
	"time"
)

type Payment struct {
	ID                int       `json:"id" db:"id"`
	OrderID           string    `json:"order_id" db:"order_id"`               // Your internal order reference
	Amount            int64   `json:"amount" db:"amount"`                   // Amount in rupees
	Currency          string    `json:"currency" db:"currency"`               // e.g. "INR"
	Gateway           string    `json:"gateway" db:"gateway"`                 // e.g. "razorpay"
	CustomerEmail     string    `json:"customer_email" db:"customer_email"`   // User email
	CustomerPhone     string    `json:"customer_phone" db:"customer_phone"`   // User phone
	CustomerName      string    `json:"customer_name" db:"customer_name"`     // Optional name
	RazorpayOrderID   string    `json:"razorpay_order_id" db:"razorpay_order_id"`
	RazorpayPaymentID string    `json:"razorpay_payment_id" db:"razorpay_payment_id"`
	RazorpaySignature string    `json:"razorpay_signature" db:"razorpay_signature"`
	Status            string    `json:"status" db:"status"`                   // e.g. "created", "success", "failed"
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}


type VerifyPayment struct {
	OrderID   string  `json:"order_id"`
	PaymentID string  `json:"payment_id"`
	Signature string  `json:"signature"`
	Name      string  `json:"name"`
	Phone     string  `json:"phone"`
	Amount    float64 `json:"amount"`
	Gateway   string  `json:"gateway"`
}
