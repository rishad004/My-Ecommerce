package controllers

import (
	"fmt"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

func ShowUser(c *gin.Context) {
	fmt.Println("")
	fmt.Println("-------------------------SHOWING USERS-------------------------")

	var users []models.Users
	var status string
	var l []gin.H

	database.Db.Order("ID asc").Find(&users)

	for i := 0; i < len(users); i++ {
		if users[i].Blocking {
			status = "Active"
		} else {
			status = "Blocked"
		}
		l = append(l, gin.H{
			"Id":     users[i].ID,
			"Name":   users[i].Name,
			"Email":  users[i].Email,
			"Phone":  users[i].Phone,
			"Gender": users[i].Gender,
			"Status": status,
		})
	}
	c.JSON(200, gin.H{"Users": l})
}

func BlockingUser(c *gin.Context) {
	fmt.Println("")
	fmt.Println("-------------------------BLOCKING USER-------------------------")

	Id := c.Param("Id")

	var u models.Users

	database.Db.First(&u, "ID=?", Id)
	if u.ID != 0 {
		if u.Blocking {
			database.Db.Model(&u).Update("Blocking", false)
			c.JSON(200, gin.H{"Message": "Blocked user successfully"})
		} else {
			database.Db.Model(&u).Update("Blocking", true)
			c.JSON(200, gin.H{"Message": "Unblocked user successfully"})
		}
	} else {
		c.JSON(404, gin.H{"Error": "User not found"})
	}
}
