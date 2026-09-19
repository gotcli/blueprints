package blueprints

var RouteRoot = `
	package routes

	import (
			"github.com/gofiber/fiber/v2"
			"gorm.io/gorm"
			"{{ .ModuleName }}/pkg/response"
	)

	func RegisterRootRoutes(db *gorm.DB, route *fiber.App) {

			routeGroup := route.Group("api")
			routeGroup.Get("", func(c *fiber.Ctx) error {
					return response.Send(c, fiber.StatusOK, "version 1.0.0")
			})

			// All api v1 goes inside this function.
			RegisterV1Routes(db, routeGroup)

			// got:api-routes
	}									
`

var RouteV1 = `
	package routes

	import (
			"github.com/gofiber/fiber/v2"
			"gorm.io/gorm"
			"{{ .ModuleName }}/pkg/response"
	)

	func RegisterV1Routes(db *gorm.DB, route fiber.Router) {

			routeGroup := route.Group("v1")
			routeGroup.Get("", func(c *fiber.Ctx) error {
					return response.Send(c, fiber.StatusOK, "version 1.0.0")
			})
	}
`
