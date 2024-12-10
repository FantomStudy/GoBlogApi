package main

import (
	"API/database"
	"API/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Подключение к базе данных
	database.Connect()

	//Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Выгрузка порта
	port := os.Getenv("PORT")

	// Создание приложения
	app := fiber.New()

	// Конфиг (Маршрутизация, порт)
	routes.Setup(app)
	app.Listen(":" + port)
}
