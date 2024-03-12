package helper

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPass(pass string) string {
	fmt.Println("")
	fmt.Println("------------------------HASHING PASS-----------------------------")
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
	if err != nil {
		fmt.Println("")
		fmt.Println("")
		fmt.Println("ERROR ON HASHING>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		fmt.Println("")
		fmt.Println("")
	}
	return string(hashedpass)
}
