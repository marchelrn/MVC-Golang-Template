package config

import (
	"log"
	"os"
	"strconv"

	"mini_jira/utils"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	PORT                 int
	IsProd               bool
	DB_URL               string
	JWTPrivateKeyPath    string
	JWTPublicKeyPath     string
	JWTPrivateKey        string
	JWTPublicKey         string
	JWTIssuer            string
	JWTAudience          string
	JWTAccessTTLMinutes  int
	JWTRefreshTTLDays    int
	JWTRefreshHashSecret string
	SMTPHost             string
	SMTPPort             string
	SMTPEmail            string
	SMTPPassword         string
}

var config *AppConfig

func GetConfig() *AppConfig {
	return config
}

func Load() {
	log.SetOutput(os.Stdout)
	log.Println("Loading configuration...")
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	isProd := utils.SafeCompareStrings(os.Getenv("IS_PRODUCTION"), "TRUE")

	config = &AppConfig{
		PORT:                 port,
		IsProd:               isProd,
		DB_URL:               Production(),
		JWTPrivateKeyPath:    os.Getenv("JWT_PRIVATE_KEY_PATH"),
		JWTPublicKeyPath:     os.Getenv("JWT_PUBLIC_KEY_PATH"),
		JWTPrivateKey:        os.Getenv("JWT_PRIVATE_KEY"),
		JWTPublicKey:         os.Getenv("JWT_PUBLIC_KEY"),
		JWTIssuer:            utils.DefaultString(os.Getenv("JWT_ISSUER"), "mini_jira"),
		JWTAudience:          utils.DefaultString(os.Getenv("JWT_AUDIENCE"), "mini_jira_api"),
		JWTAccessTTLMinutes:  utils.DefaultInt(os.Getenv("JWT_ACCESS_TTL_MINUTES"), 15),
		JWTRefreshTTLDays:    utils.DefaultInt(os.Getenv("JWT_REFRESH_TTL_DAYS"), 7),
		JWTRefreshHashSecret: os.Getenv("JWT_REFRESH_HASH_SECRET"),
		SMTPHost:             utils.DefaultString(os.Getenv("SMTP_HOST"), "smtp.gmail.com"),
		SMTPPort:             utils.DefaultString(os.Getenv("SMTP_PORT"), "587"),
		SMTPEmail:            os.Getenv("SMTP_EMAIL"),
		SMTPPassword:         os.Getenv("SMTP_PASSWORD"),
	}
}

func Production() string {
	if config == nil || !config.IsProd {
		log.Println("Using local database")
		return LocalDb()
	}
	log.Println("Using production database")
	return os.Getenv("DATABASE_URL")
}

func LocalDb() string {
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")

	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	return "postgres://" + db_user + ":" + db_password + "@" + db_host + ":" + db_port + "/" + db_name + "?sslmode=disable"
}
