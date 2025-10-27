package utils

import (
	"crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
	"bytes"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "payment_wrapper/config"
	"log"
)

func CreateRazorpayPayment(amount int64, currency, orderID, email, name string) (string, error) {
	url := "https://api.razorpay.com/v1/orders"

	payload := map[string]interface{}{
		"amount": amount * 100, // Razorpay accepts paisa
		"currency": currency,
		"receipt": orderID,
		"payment_capture": 1,
		"notes": map[string]string{
			"customer_name":  name,
			"customer_email": email,
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.SetBasicAuth(config.RazorpayKey, config.RazorpaySecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	if id, ok := result["id"].(string); ok {
    return "https://checkout.razorpay.com/v1/checkout/embedded?order_id=" + id, nil
}

// Log Razorpay error response
log.Printf("Razorpay response: %s\n", string(respBody))
return "", fmt.Errorf("failed to create payment: %s", string(respBody))

}



// func VerifyRazorpayPayment(paymentID, orderID, signature string) (string, error) {
//     url := fmt.Sprintf("https://api.razorpay.com/v1/payments/%s", paymentID)
//     req, _ := http.NewRequest("GET", url, nil)
//     req.SetBasicAuth(config.RazorpayKey, config.RazorpaySecret)

//     client := &http.Client{}
//     resp, err := client.Do(req)
//     if err != nil {
//         fmt.Println("Error while fetching payment:", err)
//         return "", err
//     }
//     defer resp.Body.Close()

//     body, _ := ioutil.ReadAll(resp.Body)

//     fmt.Println("Razorpay Fetch Response:", string(body))
//     fmt.Println("HTTP Status:", resp.StatusCode)

//     if resp.StatusCode != http.StatusOK {
//         return "", fmt.Errorf("razorpay fetch error: %s", string(body))
//     }

//     var result map[string]interface{}
//     json.Unmarshal(body, &result)

//     fmt.Println("Parsed Payment Data:", result)

//     // Step 3: Verify signature
//     data := orderID + "|" + paymentID
//     h := hmac.New(sha256.New, []byte(config.RazorpaySecret))
//     h.Write([]byte(data))
//     expectedSignature := hex.EncodeToString(h.Sum(nil))

//     fmt.Println("🧾 Expected Signature:", expectedSignature)
//     fmt.Println("🧾 Received Signature:", signature)

//     if expectedSignature != signature {
//         fmt.Println("Signature verification failed")
//         return "failed", fmt.Errorf("signature verification failed")
//     }

//     fmt.Println("Signature verified successfully!")
//     return "success", nil
// }

func VerifyRazorpayPayment(paymentID, orderID, signature string) (string, error) {
	url := "https://api.razorpay.com/v1/payments/" + paymentID
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(config.RazorpayKey, config.RazorpaySecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "failed", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// Optional: read customer name from notes if needed
	if notes, ok := result["notes"].(map[string]interface{}); ok {
		if name, ok := notes["customer_name"].(string); ok {
			// You can save this to DB too
			fmt.Println("Customer Name:", name)
		}
	}

	// Signature verification
	data := orderID + "|" + paymentID
	h := hmac.New(sha256.New, []byte(config.RazorpaySecret))
	h.Write([]byte(data))
	expected := hex.EncodeToString(h.Sum(nil))

	if expected != signature {
		return "failed", fmt.Errorf("signature verification failed")
	}

	return "success", nil
}



