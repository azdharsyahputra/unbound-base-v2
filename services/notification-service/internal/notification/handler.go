package notification

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetNotificationsHandler(c *fiber.Ctx, db *gorm.DB) error {
	userIDStr := c.Params("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	var notifs []Notification

	if err := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifs).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch notifications")
	}

	return c.JSON(fiber.Map{
		"count": len(notifs),
		"data":  notifs,
	})
}
func MarkAsReadHandler(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")

	var notif Notification
	if err := db.First(&notif, id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "notification not found")
	}

	if notif.Read {
		return c.JSON(fiber.Map{"message": "already read"})
	}

	notif.Read = true
	if err := db.Save(&notif).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update notification")
	}

	return c.JSON(fiber.Map{
		"message": "notification marked as read",
	})
}
func DeleteNotificationHandler(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")

	if err := db.Delete(&Notification{}, id).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete notification")
	}

	return c.JSON(fiber.Map{
		"deleted": true,
	})
}
