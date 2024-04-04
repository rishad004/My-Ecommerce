package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

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
	marginX := 10.0
	marginY := 20.0

	lineHt := 10.0
	const colNumber = 7

	pdf := fpdf.New("P", "mm", "A4", "")

	pdf.SetMargins(marginX, marginY, marginX)
	pdf.AddPage()

	pdf.ImageOptions("assets/logo.gif", 5, 0, 75, 25, false, fpdf.ImageOptions{ImageType: "GIF", ReadDpi: true}, 0, "")
	pdf.Ln(5)

	header := [colNumber]string{"Order Id", "User", "Product Name", "Quantity", "Amount", "Date and Time", "Status"}
	colWidth := [colNumber]float64{20.0, 20.0, 35.0, 20.0, 20.0, 50.0, 20.0}

	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	for colJ := 0; colJ < colNumber; colJ++ {
		pdf.CellFormat(colWidth[colJ], lineHt, header[colJ], "1", 0, "C", true, 0, "")
	}

	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 7)

	for _, v := range orders {

		pdf.CellFormat(colWidth[0], lineHt, fmt.Sprintf("%d", v.OrderId), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth[1], lineHt, v.Order.User.Name, "1", 0, "LM", false, 0, "")
		pdf.CellFormat(colWidth[2], lineHt, v.Prdct.Name, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth[3], lineHt, strconv.Itoa(v.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth[4], lineHt, strconv.Itoa(int(v.Order.Amount)), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth[5], lineHt, v.Order.CreatedAt.String(), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth[6], lineHt, v.Status, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)

	}

	path := "./ReportPDF/sales_report.pdf"
	if err := pdf.OutputFileAndClose(path); err != nil {
		c.JSON(500, gin.H{"Error": "Couldn't download the report!", "err": err.Error()})
		return
	}

	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(path)

	c.JSON(200, gin.H{"Message": "Pdf downloaded successfully!"})
}
