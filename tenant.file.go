package blueprints

// Enterprise tier: tenant isolation scaffolding for multi-tenant deployments.
var TenantMiddleware = `package middlewares

import "github.com/gofiber/fiber/v2"

const TenantIDLocal = "tenant_id"

func RequireTenant() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Get("X-Tenant-ID")
		if tenantID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing X-Tenant-ID header")
		}
		c.Locals(TenantIDLocal, tenantID)
		return c.Next()
	}
}

func TenantID(c *fiber.Ctx) (string, bool) {
	tenantID, ok := c.Locals(TenantIDLocal).(string)
	return tenantID, ok
}
`

var TenantScope = `package databases

import "gorm.io/gorm"

func ForTenant(tenantID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ?", tenantID)
	}
}
`

var TenantEntity = `
	package {{ .ApiName }}

	import "time"

	type {{ .TypeName }} struct {
		ID        int
		TenantID  string
		CreatedAt time.Time
	}
`
