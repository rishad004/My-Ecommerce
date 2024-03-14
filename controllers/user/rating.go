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
	var product models.Products

	Logged := c.MustGet("Id").(uint)

	c.BindJSON(&rating)

	id, _ := strconv.Atoi(c.Param("Id"))
	database.Db.First(&product, id)

	if Logged != 0 {
		if product.Id == 0 {
			c.JSON(404, gin.H{"message": "Product not found."})
		} else {
			rating.Prdct_Id = uint(id)
			rating.User_Id = Logged

			database.Db.Create(&rating)
			c.JSON(201, gin.H{"message": "Rating added successfully!"})
		}
	} else {
		c.JSON(401, gin.H{"error": "Please login first."})
	}

}