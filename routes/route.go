package routes

import (
	"API/controller"
	"API/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Не защищены
	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)

	// Использование мидлвейр
	app.Use(middleware.IsAuthenticate)

	// Требуется авторизация
	app.Post("/api/post", controller.CreatePost)

	app.Post("/api/all-post/:id/like", controller.AddLike)           // Добавление лайка
	app.Delete("api/all-post/:id/like", controller.RemoveLike)       // Удаление лайка
	app.Get("api/all-post/:id/like/count", controller.GetLikesCount) // Количество лайков
	app.Get("api/like/:postId", controller.GetUserLikeStatus)        // Проверить, поставил ли пользователь лайк

	app.Get("/api/all-post/search", controller.FilterPosts)
	app.Get("/api/all-post", controller.AllPost) // /api/all-post?page=X где X - номер страницы
	app.Get("/api/all-post/:id", controller.DetailPost)
	app.Patch("/api/update-post/:id", controller.UpdatePost)
	app.Get("/api/unique-post", controller.UniquePost)
	app.Delete("/api/delete-post/:id", controller.DeletePost)
}
