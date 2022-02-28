package version

import "github.com/gofiber/fiber/v2"

var (
	Commit    = "unset"
	BuildTime = "unset"
	Release   = "unset"
)

func VersionHandler(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("content-type", "application/json")

	return c.Status(fiber.StatusOK).JSON(map[string]string{
		"buildTime": BuildTime,
		"commit":    Commit,
		"release":   Release,
	})
}
