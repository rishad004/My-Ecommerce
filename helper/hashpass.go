package helper

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPass(c *gin.Context, pass string) string {
	fmt.Println("")
	fmt.Println("------------------------HASHING PASS-----------------------------")
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
	if err != nil {
		c.JSON(303, "Can't hash the password")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("ERROR ON HASHING>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		fmt.Println("")
		fmt.Println("")
	}
	return string(hashedpass)
}
