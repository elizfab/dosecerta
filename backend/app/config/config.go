package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI   string
	DBName     string
	ServerPort string
	AppEnv     string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	return &Config{
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:   getEnv("DB_NAME", "dose-certa-backend"),
		// Plataformas serverless (Vercel, etc.) injetam a porta via PORT.
		// Cai para SERVER_PORT e, por fim, 8080 no desenvolvimento local.
		ServerPort: getEnv("PORT", getEnv("SERVER_PORT", "8080")),
		AppEnv:     getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
