package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Wenell09/lavendera-api/internal/config"
	"github.com/Wenell09/lavendera-api/internal/database"
	"github.com/Wenell09/lavendera-api/internal/middleware"
)

func main() {
	ctx := context.Background()
	// 1. Load Konfigurasi .env menggunakan Viper
	config.LoadConfig()
	// 2. Inisialisasi Database GORM (Koneksi ORM)
	db := database.DBConnection()
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}
	// 3. Inisialisasi PGX Pool (Koneksi Khusus RLS Transaction)
	poolConfig, err := pgxpool.ParseConfig(config.ENV.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to parse DB URL for PGX Pool: %v", err)
	}
	dbPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to PGX Pool: %v", err)
	}
	defer dbPool.Close()
	// 4. Inisialisasi Router Gin
	r := gin.Default()
	// Public Routes (Bisa diakses tanpa Login)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Lavendera API Gateway is running 🚀",
		})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Lavendera API Gateway is running 🚀",
		})
	})

	// Route Auth (Login / Register)
	// authGroup := r.Group("/api/v1/auth")
	{
		// authGroup.POST("/login", authHandler.Login)
	}

	// Protected Routes (Butuh Bearer Token JWT + RLS Active)
	protected := r.Group("/api/v1")
	// Rantai Middleware:
	// 1. Ekstrak & Verifikasi Token -> simpan tenant_id di context
	// 2. Ambil tenant_id -> jalankan SET LOCAL RLS di PostgreSQL
	protected.Use(middleware.JWTMiddleware(config.ENV.JWTSecret))
	protected.Use(middleware.RLSMiddleware(dbPool))
	{
		// Contoh Route Terproteksi:
		// protected.GET("/orders", orderHandler.GetOrders)
		// protected.POST("/orders", orderHandler.CreateOrder)
	}
	log.Printf("Server running on port %s 🚀\n", config.ENV.Port)
	if err := r.Run(":" + config.ENV.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
