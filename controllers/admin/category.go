package controllers

import (
	"fmt"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

type showc struct {
	Id          uint
	Name        string
	Description string
	Status      string
}

func ShowCategory(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY SHOWING----------------------")

	var category []models.Category
	var Status string
	var show []showc

	database.Db.Order("Id asc").Find(&category)

	for i := 0; i < len(category); i++ {
		if category[i].Blocking {
			Status = "Active"
		} else {
			Status = "Blocked"
		}
		l := showc{
			Id:          category[i].Id,
			Name:        category[i].Name,
			Description: category[i].Dscptn,
			Status:      Status,
		}
		show = append(show, l)
	}
	c.JSON(200, gin.H{"categories": show})
}

func AddCtgry(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY ADDING----------------------")

	var ctgry models.Category

	c.BindJSON(&ctgry)
	ctgry.Blocking = true

	e := database.Db.Create(&ctgry)
	if e.Error != nil {
		c.JSON(422, gin.H{"error": "Fill Category details correctly"})
	} else {
		c.JSON(200, gin.H{"message": "Created Category successfully"})
	}
}

func EditCategory(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY EDITING----------------------")

	name := c.Param("Name")

	var ctgry models.Category
	var check models.Category

	c.BindJSON(&ctgry)

	database.Db.First(&check, "Name=?", name)

	if check.Id == 0 {
		c.JSON(404, gin.H{"error": "Category not found"})
	} else {
		database.Db.Model(&check).Update("Name", ctgry.Name)
		database.Db.Model(&check).Update("Dscptn", ctgry.Dscptn)
		c.JSON(200, gin.H{"message": "Category edited successfully"})
	}

}

func DeleteCategory(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY DELETING----------------------")

	name := c.Param("Name")

	var ctgry models.Category
	var product models.Products

	database.Db.First(&ctgry, "Name=?", name)
	database.Db.First(&product, "Ctgry_Id=?", ctgry.Id)

	if product.ID != 0 {
		c.JSON(409, "You can't delete, There are some products in this category.")
	} else {
		e := database.Db.Delete(&ctgry)
		if e.Error != nil {
			c.JSON(422, gin.H{"error": "Couldn't delete the category, Please try again."})
		} else {
			c.JSON(200, gin.H{"message": "Category deleted successfully"})
		}
	}
}

func BlockingCategory(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY BLOCKING----------------------")

	name := c.Param("Name")

	var ctgry models.Category
	var product []models.Products

	database.Db.First(&ctgry, "Name=?", name)
	database.Db.Find(&product, "Ctgry_Id=?", ctgry.Id)

	if !ctgry.Blocking {
		database.Db.Model(&ctgry).Update("Blocking", true)

		for i := 0; i < len(product); i++ {
			product[i].CtgryBlock = true
			database.Db.Save(&product[i])
		}

		c.JSON(200, gin.H{
			"message": "Category unblocked successfully",
		})
	} else {
		database.Db.Model(&ctgry).Update("Blocking", false)

		for i := 0; i < len(product); i++ {
			product[i].CtgryBlock = false
			database.Db.Save(&product[i])
		}

		c.JSON(200, gin.H{
			"message": "Category blocked successfully",
		})
	}

}
