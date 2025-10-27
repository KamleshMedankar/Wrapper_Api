package db

import (
	_ "github.com/go-sql-driver/mysql"
	"payment_wrapper/models"
	
	"fmt"
)

func InsertPayment(p *models.Payment) error {
	query := `
		INSERT INTO payment
		(order_id, razorpay_order_id, amount, currency, customer_name, customer_email, gateway, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`
	_, err := Conn.Exec(query, p.OrderID, p.RazorpayOrderID, p.Amount, p.Currency, p.CustomerName, p.CustomerEmail, p.Gateway, p.Status)
	return err
}


// UpdatePaymentStatus updates payment status after verification
func UpdatePaymentStatus(orderID, status, paymentID, signature string) error {
	query := `
		UPDATE payments
		SET status = ?, razorpay_payment_id = ?, razorpay_signature = ?, updated_at = NOW()
		WHERE razorpay_order_id = ?
	`
	res, err := Conn.Exec(query, status, paymentID, signature, orderID)
	if err != nil {
		fmt.Println("DB Update Error:", err)
		return err
	}
	rows, _ := res.RowsAffected()
	fmt.Println("Rows affected:", rows)
	return err
}
