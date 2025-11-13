package notification

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	authclient "unbound-v2/notification-service/internal/grpcclient"
)

func RegisterNotificationRoutes(app *fiber.App, db *gorm.DB) {
	// Inisialisasi grpc client DI DALAM routes
	authClient := authclient.NewAuthClient("auth-service:50051")

	r := app.Group("/notifications")

	log.Println("🔧 Registering /notifications routes")

	r.Get("/:userID", func(c *fiber.Ctx) error {
		userIDParam := c.Params("userID")
		userID, _ := strconv.ParseUint(userIDParam, 10, 64)

		var notifs []Notification
		if err := db.Where("user_id = ?", userID).
			Order("created_at desc").
			Find(&notifs).Error; err != nil {
			return fiber.NewError(500, "failed to fetch notifications")
		}

		var results []fiber.Map

		for _, n := range notifs {
			sender, err := authClient.GetUserByID(uint64(n.SenderID))

			senderUsername := "(unknown)"
			if err == nil && sender != nil {
				senderUsername = sender.Username
			}

			results = append(results, fiber.Map{
				"id":          n.ID,
				"type":        n.Type,
				"message":     n.Message,
				"sender_id":   n.SenderID,
				"sender_name": senderUsername,
				"read":        n.Read,
				"created_at":  n.CreatedAt,
			})

		}

		return c.JSON(fiber.Map{
			"count": len(results),
			"data":  results,
		})
	})

	r.Patch("/:id/read", func(c *fiber.Ctx) error {
		return MarkAsReadHandler(c, db)
	})

	r.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteNotificationHandler(c, db)
	})
}
