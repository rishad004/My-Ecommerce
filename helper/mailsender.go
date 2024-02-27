package helper

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

func SendMail(c *gin.Context, mail string, sub string, body string) {
	fmt.Println("")
	fmt.Println("---------------------SENDING MAIL-----------------------")
	auth := smtp.PlainAuth(
		"Rishad's Project",
		os.Getenv("SENDING_MAIL"),
		os.Getenv("APP_PASS"),
		"smtp.gmail.com",
	)
	msg := "Subject: " + sub + "\n" + body
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("SENDING_MAIL"),
		[]string{mail},
		[]byte(msg),
	)
	if err != nil {
		c.JSON(503, "We couldn't send the mail, Please check email address")
		return
	}
}
