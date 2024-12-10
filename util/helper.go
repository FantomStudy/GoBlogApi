package util

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GetSecretKey() string {
	// Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Выгрузка секретного ключа и его возврат
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		panic("Secret key is not set in .env file")
	}

	return secretKey
}

func GenerateJwt(userId string) (string, error) {
	// Создаем claims, которые содержатся в токене
	claims := &jwt.RegisteredClaims{
		Issuer:    userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Истекает через 24 часа
	}

	// Создаем новый токен с указанным методом подписи и claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем и возвращаем строку токена
	return token.SignedString([]byte(GetSecretKey()))
}

func ParseJwt(cookie string) (string, error) {
	// Парсим токен
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Проверка по ключу
		return []byte(GetSecretKey()), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	// Преобразуем claims в правильный тип
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	// Возвращаем отправителя из claims
	return claims.Issuer, nil
}
