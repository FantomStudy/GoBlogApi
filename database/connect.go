package database

import (
	"API/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	// Выгрузка строки подключения
	dsn := os.Getenv("DSN")

	// Подключение к БД
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} else {
		log.Println("Connected to database")
	}

	// Глобальное обращение к БД
	DB = database

	// Автомиграция при запуске
	database.AutoMigrate(
		&models.User{},
		&models.Photo{},
		&models.Post{},
		&models.Like{},
	)
}
