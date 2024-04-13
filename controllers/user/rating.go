package controllers

import (
	"fmt"
	"strconv"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

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
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "Rating should be less than or equal to 5!",
			"Data":    gin.H{},
		})
		return
	}

	err := database.Db.First(&product, id).Error
	database.Db.Find(&rates, "Prdct_Id=?", uint(id))

	if er != nil {
		if err != nil {
			c.JSON(404, gin.H{
				"Status":  "Error!",
				"Code":    404,
				"Message": "Product not found!",
				"Data":    gin.H{},
			})
		} else {
			rating.PrdctId = uint(id)
			rating.User_Id = Logged

			database.Db.Create(&rating)
			for _, v := range rates {
				sum += v.Rating
			}
			product.AvrgRating = sum / float32(len(rates))
			database.Db.Save(&product)
			c.JSON(201, gin.H{
				"Status":  "Success!",
				"Code":    201,
				"Message": "Rating and review added successfully!",
				"Data":    gin.H{},
			})
		}
	} else {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "Rating or review  already exists, Try to update it instead of adding again!",
			"Data":    gin.H{},
		})
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
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Couldn't find any rating for review for this product!",
			"Data":    gin.H{},
		})
		return
	}

	if err := database.Db.Preload("Prdct").Where("User_Id=? AND Prdct_Id=?", Logged, uint(id)).First(&rate).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Rating and review not found! Please add first!",
			"Data":    gin.H{},
		})
		return
	}

	c.BindJSON(&rating)
	if rating.Rating > 5 {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "Rating should be less than or equal to 5!",
			"Data":    gin.H{},
		})
		return
	}
	if err := database.Db.Model(&rate).Updates(&rating).Error; err != nil {
		c.JSON(500, gin.H{
			"Status":  "Error!",
			"Code":    500,
			"Message": "Couldn't update the rating or review!",
			"Data":    gin.H{},
		})
		return
	}
	for _, v := range rates {
		sum += v.Rating
	}
	rate.Prdct.AvrgRating = sum / float32(len(rates))
	if err := database.Db.Save(&rate.Prdct).Error; err != nil {
		c.JSON(500, gin.H{
			"Status":  "Error!",
			"Code":    500,
			"Error":   err.Error(),
			"Message": "Couldn't update the rating or review!",
			"Data":    gin.H{},
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":  "Error!",
		"Code":    404,
		"Message": "The rating has been updated!",
		"Data":    gin.H{},
	})
}
