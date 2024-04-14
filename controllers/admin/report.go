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

// ShowReport godoc
// @Summary Report Show
// @Description Showing sales details in admin side
// @Tags Admin Report
// @Produce  json
// @Router /admin/report [get]
func GetReportData(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SHOW SALES REPORT------------------------")

	var orders []models.Orderitem
	var today time.Time

	Filter := c.Query("filter")
	if err := database.Db.Preload("Order").Preload("Order.User").Preload("Prdct").Find(&orders).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "No orders found!",
			"Data":    gin.H{},
		})
		return
	}

	switch Filter {
	case "Today":
		today = time.Now().Truncate(24 * time.Hour)
	case "This week":
		today = time.Now().Truncate(168 * time.Hour)
	case "This month":
		today = time.Now().Truncate(730 * time.Hour)
	default:
		today = time.Now()
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
	pdf.SetFont("Arial", "B", 25)
	pdf.CellFormat(0, 0, "SALES REPORT", "1", 0, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 0, today.String()[:10]+" to "+time.Now().String()[:10], "1", 0, "R", false, 0, "")
	pdf.Ln(5)

	header := [colNumber]string{"Order Id", "Product Name", "Price / Unit", "Quantity", "Amount", "Date", "Status"}
	colWidth := [colNumber]float64{20.0, 35.0, 35.0, 20.0, 20.0, 30.0, 20.0}

	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	for colJ := 0; colJ < colNumber; colJ++ {
		pdf.CellFormat(colWidth[colJ], lineHt, header[colJ], "1", 0, "C", true, 0, "")
	}

	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)

	for _, v := range orders {
		if today == time.Now() {

			pdf.CellFormat(colWidth[0], lineHt, fmt.Sprintf("%d", v.OrderId), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[1], lineHt, v.Prdct.Name, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[2], lineHt, strconv.Itoa(v.Prdct.Price), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[3], lineHt, strconv.Itoa(v.Quantity), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[4], lineHt, strconv.Itoa(int(v.Order.Amount)), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[5], lineHt, v.Order.CreatedAt.String()[:10], "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[6], lineHt, v.Status, "1", 0, "C", false, 0, "")
			pdf.Ln(-1)

		} else {
			if time.Now().After(today) {

				pdf.CellFormat(colWidth[0], lineHt, fmt.Sprintf("%d", v.OrderId), "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[1], lineHt, v.Prdct.Name, "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[2], lineHt, strconv.Itoa(v.Prdct.Price), "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[3], lineHt, strconv.Itoa(v.Quantity), "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[4], lineHt, strconv.Itoa(int(v.Order.Amount)), "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[5], lineHt, v.Order.CreatedAt.String()[:10], "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[6], lineHt, v.Status, "1", 0, "C", false, 0, "")
				pdf.Ln(-1)
			}
		}
	}

	path := "./ReportPDF/salesReport_" + time.Now().String()[:10] + "_" + Filter + ".pdf"
	if err := pdf.OutputFileAndClose(path); err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   err.Error(),
			"Message": "Couldn't download the report!",
			"Data":    gin.H{},
		})
		return
	}

	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "salesReport_" + time.Now().String()[:10] + "_" + Filter + ".pdf"))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(path)

	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Pdf downloaded successfully!",
		"Data":    gin.H{},
	})
}
