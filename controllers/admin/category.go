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
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Retrieved category details!",
		"Data": gin.H{
			"categories": l,
		},
	})
}

// AddCategory godoc
// @Summary Category Add
// @Description Adding category with it's details
// @Tags Admin Category
// @Accept  multipart/form-data
// @Produce  json
// @Param name formData string true "cateory name"
// @Param description formData string true "cateory description"
// @Router /admin/category [post]
func AddCtgry(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY ADDING----------------------")

	var ctgry models.Category

	ctgry.Name = c.Request.FormValue("name")
	ctgry.Dscptn = c.Request.FormValue("description")
	ctgry.Blocking = true

	e := database.Db.Create(&ctgry)
	if e.Error != nil {
		c.JSON(422, gin.H{
			"Status":  "Error!",
			"Code":    422,
			"Error":   e.Error.Error(),
			"Message": "Fill Category details correctly!",
			"Data":    gin.H{},
		})
	} else {
		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Created Category successfully!",
			"Data":    gin.H{},
		})
	}
}

// EditCategory godoc
// @Summary Category Edit
// @Description Editing category with it's details
// @Tags Admin Category
// @Accept  multipart/form-data
// @Produce  json
// @Param id query string true "category id"
// @Param name formData string true "cateory name"
// @Param description formData string true "cateory description"
// @Router /admin/category [put]
func EditCategory(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------CATEGORY EDITING----------------------")

	Query := c.Query("id")

	var check models.Category

	database.Db.First(&check, "Id=?", Query)

	if check.Id == 0 {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Category not found!",
			"Data":    gin.H{},
		})
	} else {
		check.Name = c.Request.FormValue("name")
		check.Dscptn = c.Request.FormValue("description")
		database.Db.Save(&check)
		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Category edited successfully!",
			"Data":    gin.H{},
		})
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
	err := database.Db.First(&product, "Ctgry_Id=?", ctgry.Id).Error

	if err == nil {
		c.JSON(409, gin.H{
			"Status":  "Fail!",
			"Code":    409,
			"Error":   err,
			"Message": "You can't delete, There are some products in this category!",
			"Data":    gin.H{},
		})
	} else {
		e := database.Db.Delete(&ctgry)
		if e.Error != nil {
			c.JSON(422, gin.H{
				"Status":  "Error!",
				"Code":    422,
				"Error":   e.Error.Error(),
				"Message": "Couldn't delete the category, Please try again!",
				"Data":    gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"Status":  "Success!",
				"Code":    200,
				"Message": "Category deleted successfully!",
				"Data":    gin.H{},
			})
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
			"Status":  "Success!",
			"Code":    200,
			"Message": "Category unblocked successfully!",
			"Data":    gin.H{},
		})
	} else {
		database.Db.Model(&ctgry).Update("Blocking", false)

		for i := 0; i < len(product); i++ {
			product[i].CtgryBlock = false
			database.Db.Save(&product[i])
		}

		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Category blocked successfully!",
			"Data":    gin.H{},
		})
	}
}
