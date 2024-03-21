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

	var rating, alrdyRate models.Rating
	var rates []models.Rating
	var product models.Products
	var sum float32

	Logged := c.MustGet("Id").(uint)
	id, _ := strconv.Atoi(c.Param("Id"))
	er := database.Db.First(&alrdyRate, "User_Id=?", Logged).Where("Prdct_Id = ?", id)

	c.BindJSON(&rating)

	database.Db.First(&product, id)
	database.Db.Find(&rates, "Prdct_Id=?", uint(id))

	if er.Error != nil {
		if product.ID == 0 {
			c.JSON(404, gin.H{"message": "Product not found."})
		} else {
			rating.Prdct_Id = uint(id)
			rating.User_Id = Logged

			database.Db.Create(&rating)
			for _, v := range rates {
				sum += v.Rating
			}
			product.AvrgRating = sum / float32(len(rates))
			database.Db.Save(&product)
			c.JSON(201, gin.H{"message": "Rating added successfully!"})
		}
	} else {
		c.JSON(401, gin.H{"error": "Rating  already exists, Try to update it instead of adding again."})
	}

}
