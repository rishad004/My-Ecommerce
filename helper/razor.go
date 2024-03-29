package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/razorpay/razorpay-go"
)

func Executerazorpay(order string, amount int) (string, error) {

	client := razorpay.NewClient(os.Getenv("RAZOR_KEY"), os.Getenv("RAZOR_SECRET"))
	fmt.Println("")
	fmt.Println(order)

	data := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  order,
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return "", errors.New("payment not initiated")
	}
	razorId, _ := body["id"].(string)
	return razorId, nil
}

func RazorPaymentVerification(sign, orderId, paymentId string) error {
	signature := sign
	secret := os.Getenv("RAZOR_SECRET")
	data := orderId + "|" + paymentId

	h := hmac.New(sha256.New, []byte(secret))

	_, err := h.Write([]byte(data))
	if err != nil {
		panic(err)
	}

	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) != 1 {
		return errors.New("payment failed")
	} else {
		return nil
	}
}
