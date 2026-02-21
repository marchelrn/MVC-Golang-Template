package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	IsProd bool
	DBUrl string
}

var config *Config

func GetConfig() *Config {
	return config
}

func Load(){
	log.Println("Loading configuration...")
	
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	isProd := os.Getenv("ENV") == "production"
	
	config = &Config{
		Port: strconv.Itoa(port),
		IsProd: isProd,
		DBUrl: Production(),
	}

}

func Production() string {
	if config == nil || !config.IsProd{
		log.Println("Using local database")
		return LocalDb()
	}
	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	return db_url
}

func LocalDb() string {
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")

	return "postgres://" + db_user + ":" + db_pass + "@" + db_host + ":" + db_port + "/" + db_name
}
