package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port           string `mapstructure:"PORT"`
	DatabaseURL    string `mapstructure:"DATABASE_URL"`
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	SupabaseURL    string `mapstructure:"SUPABASE_URL"`
	SupabaseKey    string `mapstructure:"SUPABASE_KEY"`
	SupabaseBucket string `mapstructure:"SUPABASE_BUCKET"`
}

var ENV Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.SetDefault("PORT", "8080")
	viper.AutomaticEnv()

	if err := viper.BindEnv("PORT"); err != nil {
		log.Fatalf("Gagal bind PORT: %v", err)
	}
	if err := viper.BindEnv("DATABASE_URL"); err != nil {
		log.Fatalf("Gagal bind DATABASE_URL: %v", err)
	}
	if err := viper.BindEnv("JWT_SECRET"); err != nil {
		log.Fatalf("Gagal bind JWT_SECRET: %v", err)
	}
	if err := viper.BindEnv("SUPABASE_URL"); err != nil {
		log.Fatalf("Gagal bind SUPABASE_URL: %v", err)
	}
	if err := viper.BindEnv("SUPABASE_KEY"); err != nil {
		log.Fatalf("Gagal bind SUPABASE_KEY: %v", err)
	}
	if err := viper.BindEnv("SUPABASE_BUCKET"); err != nil {
		log.Fatalf("Gagal bind SUPABASE_BUCKET: %v", err)
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan OS Environment Variables")
	}
	if err := viper.Unmarshal(&ENV); err != nil {
		log.Fatalf("Gagal membaca konfigurasi environment: %v", err)
	}
	if ENV.DatabaseURL == "" {
		log.Fatal("FATAL: DATABASE_URL tidak boleh kosong!")
	}
	if ENV.JWTSecret == "" {
		log.Fatal("FATAL: JWT_SECRET tidak boleh kosong!")
	}
}
