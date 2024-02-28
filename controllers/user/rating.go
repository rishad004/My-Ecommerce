package controllers

import (
	"fmt"
	"project/database"
	"project/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddRating(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------PRODUCT RATING------------------------")

	var rating models.Rating

	c.BindJSON(&rating)

	id, _ := strconv.Atoi(c.Param("Id"))

	if Logged != 0 {
		rating.Prdct_Id = uint(id)
		rating.User_Id = Logged

		database.Db.Create(&rating)
		c.JSON(201, "Rating added successfully!")
	} else {
		c.JSON(401, "Please login first.")
	}

}
