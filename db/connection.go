package db

import (
	
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	
)

var Conn *sql.DB

func Connect() {
	var err error
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL not found")
	}

dsn := strings.Replace(databaseURL, "mysql://", "", 1)
	// Now looks like: user:password@host:port/dbname

	dsn = fmt.Sprintf("%s?parseTime=true", dsn)
	
	Conn, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	err = Conn.Ping()
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	fmt.Println("Connected to MySQL successfully!")

}

// // InsertPayment adds a new payment record in DB
// func InsertPayment(payment *models.Payment) error {
//     query := `
//         INSERT INTO payments (
//             order_id, razorpay_order_id, amount, currency, gateway,
//             customer_email, customer_phone, customer_name, status, created_at
//         ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
//     `
//     _, err := Conn.Exec(query,
//         payment.OrderID,
//         payment.RazorpayOrderID,
//         payment.Amount,
//         payment.Currency,
//         payment.Gateway,
//         payment.CustomerEmail,
//         payment.CustomerPhone,
//         payment.CustomerName,
//         "created",
//         time.Now(),
//     )
//     return err
// }

// // UpdatePaymentStatus updates payment status after verification
// func UpdatePaymentStatus(orderID, status string) error {
//     query := `UPDATE payments SET status = ? WHERE order_id = ?`
//     _, err := Conn.Exec(query, status, orderID)
//     return err
// }
