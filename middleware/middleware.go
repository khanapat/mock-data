package middleware

import (
	"mock-data/common"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type middleware struct {
	ZapLogger *zap.Logger
}

func NewMiddleware(zapLogger *zap.Logger) *middleware {
	return &middleware{
		ZapLogger: zapLogger,
	}
}

func (m *middleware) BasicAuthenication() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			viper.GetString("swagger.user"): viper.GetString("swagger.password"),
		},
		Realm: "Restricted",
		Authorizer: func(user, pass string) bool {
			if user == viper.GetString("swagger.user") && pass == viper.GetString("swagger.password") {
				return true
			}
			return false
		},
		Unauthorized:    nil,
		ContextUsername: "_user",
		ContextPassword: "_password",
	})
}

func (m *middleware) CorsMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH",
		AllowHeaders: "Content-Type, Origin, Authorization, Accept, OTP, ReferenceNo",
	})
}

func (m *middleware) JSONMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts(common.ApplicationJSON)
		return c.Next()
	}
}

func (m *middleware) ContextLocaleMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(common.LocaleKey, c.Query(common.LocaleKey))
		return c.Next()
	}
}

func (m *middleware) LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		if c.Request().Header.Peek(common.XRequestID) == nil {
			c.Request().Header.Add(common.XRequestID, uuid.New().String())
		}

		logger := m.ZapLogger.With(zap.String(common.XRequestID, string(c.Request().Header.Peek(common.XRequestID))))

		logger.Debug(common.RequestInfoMsg,
			zap.String("method", string(c.Request().Header.Method())),
			zap.String("host", string(c.Request().Header.Host())),
			zap.String("path_uri", c.Request().URI().String()),
			zap.String("remote_addr", c.Context().RemoteAddr().String()),
			zap.String("body", string(c.Request().Body())),
		)

		if err := c.Next(); err != nil {
			return err
		}
		logger.Debug(common.ResponseInfoMsg,
			zap.String("body", string(c.Response().Body())),
		)
		logger.Info("Summary Information",
			zap.String("method", string(c.Request().Header.Method())),
			zap.String("path_uri", c.Request().URI().String()),
			zap.Duration("duration", time.Since(start)),
			zap.Int("status_code", c.Response().StatusCode()),
		)
		return nil
	}
}
