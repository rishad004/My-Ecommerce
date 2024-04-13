package controllers

import (
	"fmt"
	"strconv"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

// ShowCategory godoc
// @Summary Category Show
// @Description Showing category details in admin side
// @Tags Admin Category
// @Produce  json
// @Router /admin/category [get]
func ShowCategory(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY SHOWING----------------------")

	var category []models.Category
	var Status string
	var l []gin.H

	database.Db.Order("Id asc").Find(&category)

	for i := 0; i < len(category); i++ {
		if category[i].Blocking {
			Status = "Active"
		} else {
			Status = "Blocked"
		}
		l = append(l, gin.H{
			"Id":          category[i].Id,
			"Name":        category[i].Name,
			"Description": category[i].Dscptn,
			"Status":      Status,
		})
	}
	c.JSON(200, gin.H{"categories": l})
}

// AddCategory godoc
// @Summary Category Add
// @Description Adding category with it's details
// @Tags Admin Category
// @Accept  json
// @Produce  json
// @Param cat body models.AddCat true "Add Category"
// @Router /admin/category [post]
func AddCtgry(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY ADDING----------------------")

	var cat models.AddCat
	var ctgry models.Category

	c.BindJSON(&cat)
	ctgry.Blocking = true

	e := database.Db.Create(&ctgry)
	if e.Error != nil {
		c.JSON(422, gin.H{"Error": "Fill Category details correctly"})
	} else {
		c.JSON(200, gin.H{"Message": "Created Category successfully"})
	}
}

// EditCategory godoc
// @Summary Category Edit
// @Description Editing category with it's details
// @Tags Admin Category
// @Accept  json
// @Produce  json
// @Param cat body models.AddCat true "Add Category"
// @Router /admin/category [put]
func EditCategory(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY EDITING----------------------")

	name := c.Param("Name")

	var cat models.AddCat
	var check models.Category

	c.BindJSON(&cat)

	database.Db.First(&check, "Name=?", name)

	if check.Id == 0 {
		c.JSON(404, gin.H{"Error": "Category not found"})
	} else {
		database.Db.Model(&check).Update("Name", cat.Name)
		database.Db.Model(&check).Update("Dscptn", cat.Description)
		c.JSON(200, gin.H{"Message": "Category edited successfully"})
	}

}

// DeleteCategory godoc
// @Summary Category Delete
// @Description Deleting category completely
// @Tags Admin Category
// @Produce  json
// @Param id query string false "name search by id"
// @Router /admin/category [delete]
func DeleteCategory(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY DELETING----------------------")

	Id, _ := strconv.Atoi(c.Query("id"))

	var ctgry models.Category
	var product models.Products

	database.Db.First(&ctgry, "Id=?", uint(Id))
	database.Db.First(&product, "Ctgry_Id=?", ctgry.Id)

	if product.ID != 0 {
		c.JSON(409, gin.H{"Message": "You can't delete, There are some products in this category."})
	} else {
		e := database.Db.Delete(&ctgry)
		if e.Error != nil {
			c.JSON(422, gin.H{"Error": "Couldn't delete the category, Please try again."})
		} else {
			c.JSON(200, gin.H{"Message": "Category deleted successfully"})
		}
	}
}

// BlockingCategory godoc
// @Summary Category Blocking/Unblocking
// @Description Blocking or unblocking category with products
// @Tags Admin Category
// @Produce  json
// @Param id query string false "name search by id"
// @Router /admin/category [patch]
func BlockingCategory(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY BLOCKING----------------------")

	Id, _ := strconv.Atoi(c.Query("id"))

	var ctgry models.Category
	var product []models.Products

	database.Db.First(&ctgry, "Id=?", uint(Id))
	database.Db.Find(&product, "Ctgry_Id=?", ctgry.Id)

	if !ctgry.Blocking {
		database.Db.Model(&ctgry).Update("Blocking", true)

		for i := 0; i < len(product); i++ {
			product[i].CtgryBlock = true
			database.Db.Save(&product[i])
		}

		c.JSON(200, gin.H{
			"Message": "Category unblocked successfully",
		})
	} else {
		database.Db.Model(&ctgry).Update("Blocking", false)

		for i := 0; i < len(product); i++ {
			product[i].CtgryBlock = false
			database.Db.Save(&product[i])
		}

		c.JSON(200, gin.H{
			"Message": "Category blocked successfully",
		})
	}

}
