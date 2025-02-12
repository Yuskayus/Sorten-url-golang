package controllers

import (
	"context" // 🔹 Tambahkan package context
	"net/http"
	"url-shortener/config"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ShortenURL(c *gin.Context) {
	var request struct {
		OriginalURL string `json:"original_url"`
		Alias       string `json:"alias"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if request.OriginalURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL tidak boleh kosong"})
		return
	}

	shortCode := request.Alias
	if shortCode == "" {
		for {
			shortCode = uuid.New().String()[:6]

			// Cek apakah shortCode sudah ada di database
			var exists bool
			err := config.DB.QueryRow(c, "SELECT EXISTS(SELECT 1 FROM urls WHERE short_code=$1)", shortCode).Scan(&exists)
			if err != nil {
				log.Println("Error checking existing short code:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa kode unik"})
				return
			}

			if !exists {
				break // Keluar dari loop jika kode unik
			}
		}
	}

	// Insert ke database
	_, err := config.DB.Exec(c, "INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", shortCode, request.OriginalURL)
	if err != nil {
		log.Println("Database Insert Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:9080/" + shortCode})
}

func RedirectURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var originalURL string
	err := config.DB.QueryRow(context.Background(), "SELECT original_url FROM urls WHERE short_code=$1", shortCode).Scan(&originalURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Tambahkan log sebelum redirect
	log.Println("Redirecting to:", originalURL)

	c.Redirect(http.StatusFound, originalURL)
}
