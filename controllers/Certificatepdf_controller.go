package controllers

import (
	"fmt"
	"image/png"
	"net/http"
	"bytes"

	"github.com/ayushwar/major/database"
	"github.com/ayushwar/major/models"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// DownloadCertificate → GET /certificates/download/:id
func DownloadCertificate(ctx *gin.Context) {
	certID := ctx.Param("id")

	// Fetch certificate with user & course
	var cert models.Certificate
	if err := database.DB.Preload("User").Preload("Course").First(&cert, certID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "certificate not found"})
		return
	}

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(190, 20, "Certificate of Completion", "", 1, "C", false, 0, "")
	pdf.Ln(20)

	// Body
	pdf.SetFont("Arial", "", 14)
	bodyText := fmt.Sprintf(
		"This is to certify that %s has successfully completed the course \"%s\" on %s.\n\nCertificate Code: %s",
		cert.User.Name,
		cert.Course.Title,
		cert.IssuedAt.Format("02 Jan 2006"),
		cert.CertCode,
	)
	pdf.MultiCell(0, 10, bodyText, "", "C", false)

	// Signature line
	pdf.Ln(20)
	pdf.SetFont("Arial", "I", 12)
	pdf.CellFormat(190, 10, "______________________", "", 1, "R", false, 0, "")
	pdf.CellFormat(190, 10, "Instructor / Admin Signature", "", 1, "R", false, 0, "")

	// ----------------- QR Code -----------------
	qrData := fmt.Sprintf("https://yourdomain.com/verify/%s", cert.CertCode)
	qrCode, _ := qr.Encode(qrData, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 50, 50) // 50x50 pixels

	var buf bytes.Buffer
	png.Encode(&buf, qrCode)

	pdf.RegisterImageOptionsReader("qr", gofpdf.ImageOptions{ImageType: "PNG"}, &buf)
	pdf.ImageOptions("qr", 150, 250, 40, 40, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")

	// Send PDF as downloadable
	filename := fmt.Sprintf("certificate_%s.pdf", cert.CertCode)
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)

	if err := pdf.Output(ctx.Writer); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate certificate"})
		return
	}
}
