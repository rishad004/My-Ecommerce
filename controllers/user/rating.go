package controllers

import (
	"fmt"
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddRating(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ADDING RATING------------------------")

	var rating, alrdyRate models.Rating
	var rates []models.Rating
	var product models.Products
	var sum float32

	Logged := c.MustGet("Id").(uint)
	id, _ := strconv.Atoi(c.Param("Id"))
	er := database.Db.Where("Prdct_Id=? AND User_Id=?", id, Logged).First(&alrdyRate).Error

	c.BindJSON(&rating)
	if rating.Rating > 5 {
		c.JSON(401, gin.H{"Error": "Rating should be less than or equal to 5"})
		return
	}

	err := database.Db.First(&product, id).Error
	database.Db.Find(&rates, "Prdct_Id=?", uint(id))

	if er != nil {
		if err != nil {
			c.JSON(404, gin.H{"Message": "Product not found."})
		} else {
			rating.PrdctId = uint(id)
			rating.User_Id = Logged

			database.Db.Create(&rating)
			for _, v := range rates {
				sum += v.Rating
			}
			product.AvrgRating = sum / float32(len(rates))
			database.Db.Save(&product)
			c.JSON(201, gin.H{"Message": "Rating and review added successfully!"})
		}
	} else {
		c.JSON(401, gin.H{"Error": "Rating or review  already exists, Try to update it instead of adding again."})
	}

}

func EditRating(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------EDIT RATING------------------------")

	id, _ := strconv.Atoi(c.Param("Id"))
	Logged := c.MustGet("Id").(uint)

	var rate, rating models.Rating
	var rates []models.Rating
	var sum float32

	if err := database.Db.Where("Prdct_Id=?", uint(id)).Find(&rates).Error; err != nil {
		c.JSON(404, gin.H{"Error": "Couldn't find any rating for review for this product!"})
		return
	}

	if err := database.Db.Preload("Prdct").Where("User_Id=? AND Prdct_Id=?", Logged, uint(id)).First(&rate).Error; err != nil {
		c.JSON(404, gin.H{"Error": "Rating and review not found! Please add both first."})
		return
	}

	c.BindJSON(&rating)
	if rating.Rating > 5 {
		c.JSON(401, gin.H{"Error": "Rating should be less than or equal to 5"})
		return
	}
	if err := database.Db.Model(&rate).Updates(&rating).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Couldn't update the rating or review."})
		return
	}
	for _, v := range rates {
		sum += v.Rating
	}
	rate.Prdct.AvrgRating = sum / float32(len(rates))
	if err := database.Db.Save(&rate.Prdct).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Couldn't update the rating or review."})
		return
	}

	c.JSON(200, gin.H{"Message": "The rating has been updated."})
}
