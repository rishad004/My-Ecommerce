package controllers

import (
	"fmt"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

// ShowHome godoc
// @Summary Home Show
// @Description Showing all Products details in user side
// @Tags User Home&Product
// @Produce  json
// @Router /user/home [get]
func UserHome(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------HOME SHOWING------------------------")

	var product []models.Products
	var show []gin.H
	var img string
	var rate float32

	database.Db.Find(&product)
	for i := 0; i < len(product); i++ {
		if product[i].CtgryBlock {
			var r []models.Rating
			database.Db.Find(&r, "Prdct_Id=?", product[i].ID)
			for _, k := range r {
				rate += k.Rating
			}
			if len(r) == 0 {
				rate = 0
			} else {
				rate = rate / float32(len(r))
			}
			if product[i].ImageURLs == nil {
				img = ""
			} else {
				img = product[i].ImageURLs[0]
			}
			show = append(show, gin.H{
				"Image":  img,
				"Name":   product[i].Name,
				"Price":  product[i].Price,
				"Rating": rate,
			})
			rate = 0
		}
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Showing products in home page!",
		"Data":    show,
	})
}

// ShowSortedHome godoc
// @Summary Sorted Home Show
// @Description Showing all Products details by sorting in user side
// @Tags User Home&Product
// @Param type query string true "sort type"
// @Produce  json
// @Router /user/sort [get]
func SortProduct(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SORTING PRODUCT------------------------")

	var p []models.Products
	var show []gin.H
	var img string

	sortType := c.Query("type")
	fmt.Println(sortType)

	switch sortType {
	case "aA-zZ":
		database.Db.Order("name asc").Find(&p)
	case "zZ-aA":
		database.Db.Order("name desc").Find(&p)
	case "low to high":
		database.Db.Order("price asc").Find(&p)
	case "high to low":
		database.Db.Order("price desc").Find(&p)
	case "recent":
		database.Db.Order("id desc").Find(&p)
	case "most rated":
		database.Db.Order("Avrg_Rating desc").Find(&p)
	default:
		c.JSON(404, gin.H{"Error": "Products not found"})
		return
	}
	for _, product := range p {
		var rate float32
		var r []models.Rating
		database.Db.Find(&r, "Prdct_Id=?", product.ID)
		for _, k := range r {
			rate += k.Rating
		}
		if len(r) == 0 {
			rate = 0
		} else {
			rate = rate / float32(len(r))
		}
		if product.CtgryBlock {
			if product.ImageURLs == nil {
				img = ""
			} else {
				img = product.ImageURLs[0]
			}
			show = append(show, gin.H{
				"Image":  img,
				"Name":   product.Name,
				"Price":  product.Price,
				"Rating": rate,
			})
		}
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Showing product by sorting!",
		"Data":    show,
	})
}

// SearchingProduct godoc
// @Summary Searching Products
// @Description Searching Products in user side
// @Tags User Home&Product
// @Param search query string true "product search"
// @Produce  json
// @Router /user/search [get]
func UserSearchP(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SEARCHING PRODUCT------------------------")

	var products []models.Products
	var show []gin.H
	var img string

	searchQuery := c.Query("search")
	fmt.Println(searchQuery)

	database.Db.Where("name ILIKE ?", "%"+searchQuery+"%").Find(&products)
	if len(products) == 0 {
		c.JSON(404, gin.H{
			"Status":  "Fail!",
			"Code":    404,
			"Message": "Products not found!",
			"Data":    gin.H{},
		})
		return
	}
	for _, product := range products {
		var rate float32
		var r []models.Rating
		database.Db.Find(&r, "Prdct_Id=?", product.ID)
		for _, k := range r {
			rate += k.Rating
		}
		if len(r) == 0 {
			rate = 0
		} else {
			rate = rate / float32(len(r))
		}
		if product.CtgryBlock {
			if product.ImageURLs == nil {
				img = ""
			} else {
				img = product.ImageURLs[0]
			}
			show = append(show, gin.H{
				"Image":  img,
				"Name":   product.Name,
				"Price":  product.Price,
				"Rating": rate,
			})
		}
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Showing searched products!",
		"Data":    show,
	})
}

// FilterProduct godoc
// @Summary Filtered Products
// @Description Filtering Products in user side
// @Tags User Home&Product
// @Param category query string true "filter search"
// @Produce  json
// @Router /user/filter [get]
func FilterProduct(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SEARCHING PRODUCT------------------------")

	var Product []models.Products
	var show []gin.H
	var img string
	var rate float32

	Category := c.Query("category")
	fmt.Println(Category)

	if err := database.Db.Preload("Ctgry").Find(&Product).Error; err != nil {
		c.JSON(404, gin.H{"Error": "Couldn't find any product!"})
		return
	}

	for _, v := range Product {
		if v.Ctgry.Name == Category {
			if v.ImageURLs == nil {
				img = ""
			} else {
				img = v.ImageURLs[0]
			}
			if v.AvrgRating != 0 {
				rate = v.AvrgRating
			} else {
				rate = 0
			}
			show = append(show, gin.H{
				"Image":  img,
				"Name":   v.Name,
				"Price":  v.Price,
				"Rating": rate,
			})
		}
	}
	if show == nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "No products found in this category!",
			"Data":    show,
		})
		return
	}
	c.JSON(200, gin.H{"Products": show})
}
