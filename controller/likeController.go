package controller

import (
	"API/database"
	"API/models"
	"API/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func AddLike(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	stringUserId, err := util.ParseJwt(cookie)
	if err != nil {
		c.Status(401)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	userId, err := strconv.ParseUint(stringUserId, 10, 32)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Incorrect User Id",
		})
	}

	postId, _ := strconv.Atoi(c.Params("id"))

	var like models.Like
	database.DB.Where("user_id=? AND post_id=?", userId, postId).First(&like)
	if like.Id != 0 {
		return c.JSON(fiber.Map{
			"message": "You already liked this post!",
		})
	}
	newLike := models.Like{
		UserId: uint(userId),
		PostId: uint(postId),
	}
	if err := database.DB.Create(&newLike).Error; err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Failed to add like",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Liked successfully!",
	})
}
func RemoveLike(c *fiber.Ctx) error {
	userId, _ := util.ParseJwt(c.Cookies("jwt"))
	postId, _ := strconv.Atoi(c.Params("id"))

	var like models.Like
	database.DB.Where("user_id = ? AND post_id = ?", userId, postId).First(&like)
	if like.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Like not found",
		})
	}

	// Удаляем лайк
	if err := database.DB.Delete(&like).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to remove like",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Like removed successfully!",
	})
}
func GetLikesCount(c *fiber.Ctx) error {
	postId, _ := strconv.Atoi(c.Params("id"))

	var likeCount int64
	if err := database.DB.Model(&models.Like{}).Where("post_id = ?", postId).Count(&likeCount).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get like count",
		})
	}

	return c.JSON(fiber.Map{
		"like_count": likeCount,
	})
}
func GetUserLikeStatus(c *fiber.Ctx) error {
	userId, _ := util.ParseJwt(c.Cookies("jwt"))
	postId, _ := strconv.Atoi(c.Params("postId"))

	var like models.Like
	database.DB.Where("user_id = ? AND post_id = ?", userId, postId).First(&like)
	if like.Id == 0 {
		return c.JSON(fiber.Map{
			"liked": false,
		})
	}

	return c.JSON(fiber.Map{
		"liked": true,
	})
}
