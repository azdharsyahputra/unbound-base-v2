package user

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"unbound-v2/services/user-service/internal/common/middleware"
	"unbound-v2/services/user-service/internal/grpcclient"
)

func RegisterFollowRoutes(app *fiber.App, db *gorm.DB, authClient *grpcclient.AuthClient) {
	r := app.Group("/users")

	// ✅ FOLLOW / UNFOLLOW USER
	r.Post("/:username/follow", middleware.JWTProtected(authClient), func(c *fiber.Ctx) error {
		targetUsername := c.Params("username")

		// Ambil user ID dari JWT (hasil validasi gRPC)
		userID, ok := c.Locals("userID").(string)
		if !ok || userID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
		}

		// Ambil target user dari auth-service via gRPC
		target, err := authClient.GetUserByUsername(targetUsername)
		if err != nil || target == nil {
			return fiber.NewError(fiber.StatusNotFound, "target user not found")
		}

		// Tidak boleh follow diri sendiri
		if fmt.Sprintf("%d", target.Id) == userID {
			return fiber.NewError(fiber.StatusBadRequest, "you can't follow yourself")
		}

		// Cek apakah sudah follow
		var existing Follow
		if err := db.Where("follower_id = ? AND following_id = ?", ParseUint(userID), target.Id).First(&existing).Error; err == nil && existing.ID != 0 {
			// Sudah follow → unfollow
			if err := db.Delete(&existing).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to unfollow")
			}
			return c.JSON(fiber.Map{"following": false})
		}

		// Belum follow → buat relasi baru
		newFollow := Follow{
			FollowerID:  uint(ParseUint(userID)),
			FollowingID: uint(target.Id),
		}

		if err := db.Create(&newFollow).Error; err != nil {
			if strings.Contains(err.Error(), "unique") {
				return c.JSON(fiber.Map{"following": true})
			}
			return fiber.NewError(fiber.StatusInternalServerError, "failed to follow")
		}

		return c.JSON(fiber.Map{"following": true})
	})

	// ✅ LIST FOLLOWERS
	r.Get("/:username/followers", func(c *fiber.Ctx) error {
		username := c.Params("username")

		// Ambil target user dari auth-service via gRPC
		target, err := authClient.GetUserByUsername(username)
		if err != nil || target == nil {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		var followers []struct {
			Username string `json:"username"`
		}

		query := `
			SELECT DISTINCT u.username
			FROM follows f
			JOIN users u ON u.id = f.follower_id
			WHERE f.following_id = ?
		`

		if err := db.Raw(query, target.Id).Scan(&followers).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch followers")
		}

		return c.JSON(followers)
	})

	// ✅ LIST FOLLOWING
	r.Get("/:username/following", func(c *fiber.Ctx) error {
		username := c.Params("username")

		// Ambil target user dari auth-service via gRPC
		target, err := authClient.GetUserByUsername(username)
		if err != nil || target == nil {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		var following []struct {
			Username string `json:"username"`
		}

		query := `
			SELECT DISTINCT u.username
			FROM follows f
			JOIN users u ON u.id = f.following_id
			WHERE f.follower_id = ?
		`

		if err := db.Raw(query, target.Id).Scan(&following).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch following")
		}

		return c.JSON(following)
	})
}

// Helper: konversi string ke uint64 aman
func ParseUint(s string) uint64 {
	var v uint64
	fmt.Sscan(s, &v)
	return v
}
