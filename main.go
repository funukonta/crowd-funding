package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/funukonta/crowd-funding/handler"
	"github.com/funukonta/crowd-funding/user"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	host := os.Getenv("DB_HOST")
	userdb := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslm := os.Getenv("DB_SSLM")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, userdb, pass, name, port, sslm)

	migration(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Connection to db established!")

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)

	router.Run()

}

func migration(connStr string) {
	query := `CREATE DATABASE crowdFundingDB;`

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err.Error())
	}
	db.Exec(query)
	db.Close()

	query = `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR,
		occupation VARCHAR,
		email VARCHAR,
		password_hash VARCHAR,
		avatar_file_name VARCHAR,
		role VARCHAR,
		token VARCHAR,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS campaigns (
		id SERIAL PRIMARY KEY,
		user_id INT,
		name VARCHAR,
		short_description VARCHAR,
		description TEXT,
		goal_amount INT,
		current_amount INT,
		perks TEXT,
		slug VARCHAR,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	
	CREATE TABLE IF NOT EXISTS campaign_images (
		id SERIAL PRIMARY KEY,
		campaign_id INT,
		file_name VARCHAR,
		is_primary BOOLEAN,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
	);
	
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		campaign_id INT,
		user_id INT,
		amount INT,
		status VARCHAR,
		code VARCHAR,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (campaign_id) REFERENCES campaigns(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	connStr = strings.Replace(connStr, "dbname=postgres", "dbname=crowdfundingdb", -1)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// docker run --name crowd-fund -e POSTGRES_PASSWORD=securePass -p 5432:5432 -d postgres
