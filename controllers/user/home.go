package controllers

import (
	"fmt"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

func UserHome(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------HOME SHOWING------------------------")

	var product []models.Products
	var show []home
	var rate float32

	database.Db.Find(&product)
	for i := 0; i < len(product); i++ {
		if product[i].CtgryBlock {
			var r []models.Rating
			database.Db.Find(&r, "Prdct_Id=?", product[i].Id)
			for _, k := range r {
				rate += k.Rating
			}
			if len(r) == 0 {
				rate = 0
			} else {
				rate = rate / float32(len(r))
			}
			l := home{
				Image:  product[i].ImageURLs[0:1],
				Name:   product[i].Name,
				Price:  product[i].Price,
				Rating: rate,
			}
			show = append(show, l)
			rate = 0
		}
	}
	c.JSON(200, show)
}
