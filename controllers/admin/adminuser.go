package controllers

import (
	"fmt"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

func BlockingUser(c *gin.Context) {
	fmt.Println("")
	fmt.Println("-------------------------BLOCKING USER-------------------------")

	Id := c.Param("Id")

	var u models.Users

	database.Db.First(&u, "ID=?", Id)
	if u.ID != 0 {
		if u.Blocking {
			database.Db.Model(&u).Update("Blocking", false)
			c.JSON(200, "Blocked user successfully")
		} else {
			database.Db.Model(&u).Update("Blocking", true)
			c.JSON(200, "Unblocked user successfully")
		}
	} else {
		c.JSON(404, "User not found")
	}
}
