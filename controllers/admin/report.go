package controllers

import (
	"fmt"
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pdf/fpdf"
)

func GetReportData(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SHOW SALES REPORT------------------------")

	var orders []models.Orderitem
	var today time.Time

	Filter := c.Query("filter")

	switch Filter {
	case "Today":
		today = time.Now().Truncate(24 * time.Hour)
		if err := database.Db.Preload("Order").Preload("Order.User").Preload("Prdct").Where("Order.Created_At>=?", today).Find(&orders).Error; err != nil {
			c.JSON(404, gin.H{"Message": "No orders found"})
			return
		}
	case "This week":
		today = time.Now().Truncate(168 * time.Hour)
		if err := database.Db.Preload("Order").Preload("Order.User").Preload("Prdct").Where("Order.Created_At>=?", today).Find(&orders).Error; err != nil {
			c.JSON(404, gin.H{"Message": "No orders found"})
			return
		}
	case "This month":
		today = time.Now().Truncate(730 * time.Hour)
		if err := database.Db.Preload("Order").Preload("Order.User").Preload("Prdct").Where("Order.Created_At>=?", today).Find(&orders).Error; err != nil {
			c.JSON(404, gin.H{"Message": "No orders found"})
			return
		}
	default:
		if err := database.Db.Preload("Order").Preload("Order.User").Preload("Prdct").Find(&orders).Error; err != nil {
			c.JSON(404, gin.H{"Message": "No orders found"})
			return
		}
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	header := []string{"Order Id", "User", "Product Name", "Quantity", "Amount", "Date and Time", "Status"}
	pdf.SetFont("Arial", "B", 8)
	pdf.Cell(20, 5, header[0])
	pdf.Cell(20, 5, header[1])
	pdf.Cell(25, 5, header[2])
	pdf.Cell(20, 5, header[3])
	pdf.Cell(20, 5, header[4])
	pdf.Cell(47, 5, header[5])
	pdf.Cell(20, 5, header[6])
	pdf.Ln(-1)

	for _, v := range orders {
		pdf.SetFont("Arial", "", 7)
		pdf.Cell(20, 5, strconv.Itoa(v.OrderId))
		pdf.Cell(20, 5, v.Order.User.Name)
		pdf.Cell(25, 5, v.Prdct.Name)
		pdf.Cell(20, 5, strconv.Itoa(v.Quantity))
		pdf.Cell(20, 5, strconv.Itoa(int(v.Order.Amount)))
		pdf.Cell(47, 5, v.Order.CreatedAt.String())
		pdf.Cell(20, 5, v.Status)
		pdf.Ln(-1)
	}

	path := "./ReportPDF/sales_report.pdf"
	if err := pdf.OutputFileAndClose(path); err != nil {
		c.JSON(500, gin.H{"Error": "Couldn't download the report!"})
		return
	}

	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(path)

	c.JSON(200, gin.H{"Message": "Pdf downloaded successfully!"})
}
