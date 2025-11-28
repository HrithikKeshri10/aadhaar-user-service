package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Client() *gorm.DB {
	return DB
}

func Connect() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "aadhaar_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	var dsn string
	if password != "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s",
			host, port, user, dbname, sslmode)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Unable to open database, err: %v\n", err)
		panic(err)
	}

	sql, err := db.DB()
	if err != nil {
		fmt.Printf("Unable to get sql database from gorm, err: %v\n", err)
		panic(err)
	}

	if err := sql.Ping(); err != nil {
		fmt.Printf("Unable to connect to database, err: %v\n", err)
		panic(err)
	}

	fmt.Println("Successfully connected to the postgres db")

	DB = db
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
