package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProfileResponse struct {
	ID       uint       `json:"id"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Posts    []UserPost `json:"posts"`
}

type UserPost struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func RegisterProfileRoutes(app *fiber.App, db *gorm.DB) {
	r := app.Group("/users")

	r.Get("/:username", func(c *fiber.Ctx) error {
		username := c.Params("username")

		var user struct {
			ID       uint
			Username string
			Email    string
		}
		if err := db.Raw(`
			SELECT id, username, email
			FROM users
			WHERE username = ?
			LIMIT 1
		`, username).Scan(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch user")
		}

		if user.ID == 0 {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		var posts []UserPost
		if err := db.Raw(`
			SELECT id, content, created_at
			FROM posts
			WHERE user_id = ?
			ORDER BY created_at DESC
		`, user.ID).Scan(&posts).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch posts")
		}

		resp := ProfileResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Posts:    posts,
		}

		return c.JSON(resp)
	})
}
