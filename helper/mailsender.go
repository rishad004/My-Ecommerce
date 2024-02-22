package helper

import (
	"fmt"
	"net/smtp"

	"github.com/gin-gonic/gin"
)

func SendMail(c *gin.Context, mail string, sub string, body string) {
	fmt.Println("")
	fmt.Println("---------------------SENDING MAIL-----------------------")
	auth := smtp.PlainAuth(
		"Rishad's Project",
		"rishadmuthu004@gmail.com",
		"pcojpuhrbblrbhgv",
		"smtp.gmail.com",
	)
	msg := "Subject: " + sub + "\n" + body
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"rishadmuthu004@gmail.com",
		[]string{mail},
		[]byte(msg),
	)
	if err != nil {
		c.JSON(503, "We couldn't send the mail, Please check email address")
		return
	}
}
