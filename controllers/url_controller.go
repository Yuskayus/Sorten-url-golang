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
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("Invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Cek apakah OriginalURL kosong
	if request.OriginalURL == "" {
		log.Println("Received empty URL!")
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL tidak boleh kosong"})
		return
	}

	// Generate short code
	shortCode := uuid.New().String()[:6]
	log.Println("Generated short code:", shortCode)

	// Insert ke DB
	_, err := config.DB.Exec(context.Background(), "INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", shortCode, request.OriginalURL)
	if err != nil {
		log.Println("Error inserting into DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}

	// Return hasilnya
	shortURL := "http://localhost:9080/" + shortCode
	log.Println("Short URL generated:", shortURL)

	c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
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
