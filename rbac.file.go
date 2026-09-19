package blueprints

// Pro tier: role-based access control, built on the Roles claim issued by JWTManager (auth.file.go)
// and read back through JWTClaims (observability.file.go's sibling, JWTMiddleware in auth.file.go).
var RBACMiddleware = `package middlewares

import "github.com/gofiber/fiber/v2"

func RequireRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := JWTClaims(c)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "missing or invalid access token")
		}

		for _, required := range roles {
			for _, held := range claims.Roles {
				if held == required {
					return c.Next()
				}
			}
		}

		return fiber.NewError(fiber.StatusForbidden, "insufficient role permissions")
	}
}
`
