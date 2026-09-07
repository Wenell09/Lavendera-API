package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
}

var ENV Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Set nilai default jika variabel tidak ada di .env
	viper.SetDefault("PORT", "8080")

	// Otomatis membaca ENV dari OS (Sangat berguna untuk Deployment/Docker/Server)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, membaca dari OS Environment Variables...")
	}

	// Unmarshal nilai .env ke struct ENV
	if err := viper.Unmarshal(&ENV); err != nil {
		log.Fatalf("Gagal membaca konfigurasi environment: %v", err)
	}

	// Validasi variabel environment kritis
	if ENV.DatabaseURL == "" {
		log.Fatal("FATAL: DATABASE_URL tidak boleh kosong!")
	}

	if ENV.JWTSecret == "" {
		log.Fatal("FATAL: JWT_SECRET tidak boleh kosong!")
	}
}
