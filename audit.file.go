package blueprints

// Enterprise tier: request rate limiting and compliance audit logging.
var RateLimiterMiddleware = `package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/spf13/viper"
)

func RateLimiter() fiber.Handler {
	max := viper.GetInt("RATE_LIMIT_MAX")
	if max <= 0 {
		max = 100
	}
	window := viper.GetDuration("RATE_LIMIT_WINDOW")
	if window <= 0 {
		window = time.Minute
	}

	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: window,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "rate limit exceeded")
		},
	})
}
`

var AuditLogMiddleware = `package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"{{ .ModuleName }}/pkg/logger"
)

func AuditLog(loggers *logger.Loggers) fiber.Handler {
	return func(c *fiber.Ctx) error {
		started := time.Now()
		err := c.Next()

		subject := "anonymous"
		if claims, ok := JWTClaims(c); ok {
			subject = claims.Subject
		}

		loggers.Service.Info("audit event",
			"event", "audit",
			"subject", subject,
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", time.Since(started).Milliseconds(),
		)

		return err
	}
}
`
