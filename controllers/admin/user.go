package controllers

import (
	"fmt"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

// ShowUser godoc
// @Summary Users Show
// @Description Showing Users details in admin side
// @Tags Admin User
// @Produce  json
// @Router /admin/user [get]
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
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "User details retrieved!",
		"Data": gin.H{
			"Users": l,
		},
	})
}

// BlockingUser godoc
// @Summary User Blocking/Unblocking
// @Description Blocking or unblocking User
// @Tags Admin User
// @Produce  json
// @Param id query string true "name search by id"
// @Router /admin/user [patch]
func BlockingUser(c *gin.Context) {
	fmt.Println("")
	fmt.Println("-------------------------BLOCKING USER-------------------------")

	Id := c.Query("id")

	var u models.Users

	database.Db.First(&u, "ID=?", Id)
	if u.ID != 0 {
		if u.Blocking {
			database.Db.Model(&u).Update("Blocking", false)
			c.JSON(200, gin.H{
				"Status":  "Success!",
				"Code":    200,
				"Message": "Blocked user successfully!",
				"Data":    gin.H{},
			})
		} else {
			database.Db.Model(&u).Update("Blocking", true)
			c.JSON(200, gin.H{
				"Status":  "Success!",
				"Code":    200,
				"Message": "Unblocked user successfully!",
				"Data":    gin.H{},
			})
		}
	} else {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "User not found!",
			"Data":    gin.H{},
		})
	}
}
