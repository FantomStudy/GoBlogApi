package controller

import (
	"API/database"
	"API/models"
	"API/util"
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CreatePost(c *fiber.Ctx) error {
	// Извлечение айди пользователя из куки
	cookie := c.Cookies("jwt")
	userId, err := util.ParseJwt(cookie)
	if err != nil {
		c.Status(401)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Тело запроса
	var post models.Post
	if err := c.BodyParser(&post); err != nil {
		fmt.Println("Unable to parse body!")
	}

	validExtensions := []string{".jpg", ".jpeg", ".png", ".svg", ".webp", ".gif"}
	form, err := c.MultipartForm()
	var files []*multipart.FileHeader
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid form data!",
		})
	}
	if form != nil {
		files = form.File["photos"]
	}

	if len(files) > 0 {
		for _, file := range files {
			ext := strings.ToLower(filepath.Ext(file.Filename))

			isValidExt := false
			for _, validExtension := range validExtensions {
				if ext == validExtension {
					isValidExt = true
					break
				}
			}
			if !isValidExt {
				c.Status(400)
				return c.JSON(fiber.Map{
					"message": fmt.Sprintf("Invalid file type: %s! Only images are allowed", ext),
				})
			}
		}
	}
	// Пост создается именно авторизованным юзером
	post.UserId = userId
	if err := database.DB.Create(&post).Error; err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Failed to create post!",
		})
	}

	for _, file := range files {
		filePath := fmt.Sprintf("uploads/%s", file.Filename)

		if err := c.SaveFile(file, filePath); err != nil {
			database.DB.Delete(&post)
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to save photo!",
			})
		}

		photo := models.Photo{
			Path:   filePath,
			PostId: post.Id,
		}
		if err := database.DB.Create(&photo).Error; err != nil {
			// Если не удалось сохранить фотографию, удаляем пост и файл
			database.DB.Delete(&post)
			os.Remove(filePath)
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to save photo metadata!",
			})
		}
	}
	// Ответ о успешном создании
	return c.JSON(fiber.Map{
		"message": "Congratulate! Your post is alive!",
	})
}

func AllPost(c *fiber.Ctx) error {
	// Пагинация
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// 5 постов в одной странице
	limit := 5

	// Начало следующей страницы
	offset := (page - 1) * limit

	// Итого постов
	var total int64

	//подсчет всех постов, предварительная загрузка пользователей, настройка пагинации и запись в массив
	var posts []models.Post
	database.DB.Model(&models.Post{}).
		Preload("User").Preload("Photos").
		Count(&total).Offset(offset).Limit(limit).
		Find(&posts)

	return c.JSON(fiber.Map{
		"data": posts,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": math.Ceil(float64(total) / float64(limit)),
		},
	})
}

func DetailPost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var post models.Post
	database.DB.Where("id=?", id).
		Preload("User").Preload("Photos").
		First(&post)
	if post.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Post not found!",
		})
	}
	return c.JSON(fiber.Map{
		"data": post,
	})
}
func UpdatePost(c *fiber.Ctx) error {
	// Получиение айди из куки
	cookie := c.Cookies("jwt")
	userId, err := util.ParseJwt(cookie)
	if err != nil {
		c.Status(401)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Поиск поста с айди из роутинга
	id, _ := strconv.Atoi(c.Params("id"))
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Post not Found",
		})
	}

	// Проверка является ли создателем
	if userId != post.UserId {
		c.Status(403)
		return c.JSON(fiber.Map{
			"message": "You are not the owner of this post!",
		})
	}

	// Парсинг ответа
	if err := c.BodyParser(&post); err != nil {
		fmt.Println("Unable to parse body!")
	}

	// Обновление
	if err := database.DB.Model(&post).Updates(post).Error; err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Error updating post!",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Post updated successfully!",
	})
}
func UniquePost(c *fiber.Ctx) error {
	// Получение айди из куки
	cookie := c.Cookies("jwt")
	id, _ := util.ParseJwt(cookie)

	// Получение всех постов текущего пользователя
	var post []models.Post
	database.DB.Model(&post).Where("user_id=?", id).
		Preload("User").Preload("Photos").
		Find(&post)
	return c.JSON(post)
}

func DeletePost(c *fiber.Ctx) error {
	// Получиение айди из куки
	cookie := c.Cookies("jwt")
	userId, err := util.ParseJwt(cookie)
	if err != nil {
		c.Status(401)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Поиск поста с айди из роутинга
	id, _ := strconv.Atoi(c.Params("id"))
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Post not Found",
		})
	}

	// Проверка является ли создателем
	if userId != post.UserId {
		c.Status(403)
		return c.JSON(fiber.Map{
			"message": "You are not the owner of this post!",
		})
	}

	var photos []models.Photo
	database.DB.Where("post_id = ?", post.Id).Find(&photos)
	for _, photo := range photos {
		// Удаляем файл изображения с сервера
		if err := os.Remove(photo.Path); err != nil {
			c.Status(500)
			return c.JSON(fiber.Map{
				"message": "Failed to delete photo file",
			})
		}
	}

	// Удаление
	if err := database.DB.Delete(&post).Error; err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Error deleting post!",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Post deleted successfully!",
	})
}
func FilterPosts(c *fiber.Ctx) error {
	searchTerm := c.Query("search", "") // Получаем поисковую строку из запроса
	var posts []models.Post

	// Строим запрос с фильтрацией по названию или описанию
	database.DB.Where("posts.title LIKE ? OR posts.desc LIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%").
		Preload("User").Preload("Photos").Find(&posts)
	if len(posts) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not found such posts",
		})
	}

	return c.JSON(posts)
}
