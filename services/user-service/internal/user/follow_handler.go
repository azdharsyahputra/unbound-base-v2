package user

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"unbound-v2/services/user-service/internal/common/middleware"
	"unbound-v2/services/user-service/internal/events"
	"unbound-v2/services/user-service/internal/grpcclient"
)

// RegisterFollowRoutes menangani semua route follow/unfollow, followers, dan following
func RegisterFollowRoutes(app *fiber.App, db *gorm.DB, authClient *grpcclient.AuthClient) {
	r := app.Group("/users")
	producer := events.NewKafkaProducer()

	// ✅ FOLLOW / UNFOLLOW USER (toggle)
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

		follower := uint(ParseUint(userID))
		following := uint(target.Id)

		var existing Follow
		err = db.Where("follower_id = ? AND following_id = ?", follower, following).First(&existing).Error

		switch {
		case err == nil:
			// ✅ Sudah follow → unfollow
			if err := db.Delete(&existing).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to unfollow user")
			}

			// Kirim event Kafka untuk "unfollow"
			if perr := producer.PublishFollowEvent("unfollow", uint64(follower), uint64(following)); perr != nil {
				fmt.Printf("⚠️ Failed to publish unfollow event: %v\n", perr)
			}

			return c.JSON(fiber.Map{"following": false})

		case err == gorm.ErrRecordNotFound:
			// ✅ Belum follow → buat relasi baru
			newFollow := Follow{
				FollowerID:  follower,
				FollowingID: following,
			}

			if err := db.Create(&newFollow).Error; err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "unique") {
					return c.JSON(fiber.Map{"following": true})
				}
				return fiber.NewError(fiber.StatusInternalServerError, "failed to follow user")
			}

			// Kirim event Kafka untuk "follow"
			if perr := producer.PublishFollowEvent("follow", uint64(follower), uint64(following)); perr != nil {
				fmt.Printf("⚠️ Failed to publish follow event: %v\n", perr)
			}

			return c.JSON(fiber.Map{"following": true})

		default:
			// ✅ Error selain record not found
			return fiber.NewError(fiber.StatusInternalServerError, "database error")
		}
	})

	// ✅ LIST FOLLOWERS
	r.Get("/:username/followers", func(c *fiber.Ctx) error {
		username := c.Params("username")

		target, err := authClient.GetUserByUsername(username)
		if err != nil || target == nil {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		var followerIDs []uint64
		if err := db.Raw(`SELECT follower_id FROM follows WHERE following_id = ?`, target.Id).Scan(&followerIDs).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch followers")
		}

		var followers []fiber.Map
		for _, id := range followerIDs {
			userData, err := authClient.GetUserByID(id)
			if err == nil && userData != nil {
				followers = append(followers, fiber.Map{
					"id":       userData.Id,
					"username": userData.Username,
				})
			}
		}

		return c.JSON(followers)
	})

	// ✅ LIST FOLLOWING
	r.Get("/:username/following", func(c *fiber.Ctx) error {
		username := c.Params("username")

		target, err := authClient.GetUserByUsername(username)
		if err != nil || target == nil {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}

		var followingIDs []uint64
		if err := db.Raw(`SELECT following_id FROM follows WHERE follower_id = ?`, target.Id).Scan(&followingIDs).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch following list")
		}

		var following []fiber.Map
		for _, id := range followingIDs {
			userData, err := authClient.GetUserByID(id)
			if err == nil && userData != nil {
				following = append(following, fiber.Map{
					"id":       userData.Id,
					"username": userData.Username,
				})
			}
		}

		return c.JSON(following)
	})
}

// Helper: konversi string ke uint64 dengan aman
func ParseUint(s string) uint64 {
	var v uint64
	fmt.Sscan(s, &v)
	return v
}
