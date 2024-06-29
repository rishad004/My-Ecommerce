package helper

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail(mail string, sub string, body string) error {
	fmt.Println("")
	fmt.Println("---------------------SENDING MAIL-----------------------")
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SENDING_MAIL"))
	m.SetHeader("To", mail)
	m.SetHeader("Subject", sub)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("SENDING_MAIL"), os.Getenv("APP_PASS"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
