package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	storage_go "github.com/supabase-community/storage-go"
)

var (
	Port                         string
	AccessSecret                 string
	RefreshSecret                string
	PasswordSecret               string
	SmtpPasswd                   string
	SupabaseKey                  string
	SUPABASE_PROJECT_REF         string
	BASE_URL_FRONTEND            string
	BASE_URL_BACKEND             string
	DB_DNS                       string
	StorageClient                *storage_go.Client
	BASE_URL_BACKEND_FOR_SWAGGER string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, continuing with environment variables")
	}
	Port = ":" + os.Getenv("PORT")
	AccessSecret = os.Getenv("AccessSecret")
	RefreshSecret = os.Getenv("refreshSecret")
	PasswordSecret = os.Getenv("PasswordSecret")
	SmtpPasswd = os.Getenv("SMTP_PASSWORD")
	SupabaseKey = os.Getenv("SupabaseKey")
	SUPABASE_PROJECT_REF = os.Getenv("SUPABASE_PROJECT_REF")
	BASE_URL_BACKEND_FOR_SWAGGER = os.Getenv("BASE_URL_BACKEND_FOR_SWAGGER")
	BASE_URL_FRONTEND = os.Getenv("BASE_URL_FRONTEND")
	BASE_URL_BACKEND = os.Getenv("BASE_URL_BACKEND")
	DB_DNS = os.Getenv("DATABASE_URL")
	StorageClient = storage_go.NewClient("https://mvfitrngobsxryjosznw.supabase.co/storage/v1", SupabaseKey, nil)
}
