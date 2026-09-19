package blueprints

var RouteDynamic = `
	package routes

	import (
		"github.com/gofiber/fiber/v2"
		"gorm.io/gorm"

		repository "{{ .ModuleName }}/repositories/{{ .ApiName }}"
		service "{{ .ModuleName }}/services/{{ .ApiName }}"
		handler "{{ .ModuleName }}/handlers/{{ .ApiName }}"
	)

	func Register{{ .TypeName }}Routes(db *gorm.DB, route fiber.Router) {

		repo := repository.New{{ .TypeName }}Repository(db)
		srv := service.New{{ .TypeName }}Service(repo)
		hdr := handler.New{{ .TypeName }}Handler(srv)

		routeGroup := route.Group("{{ .ApiName }}")
		routeGroup.Get("/", hdr.Get{{ .TypeName }})
		routeGroup.Post("/", hdr.Create{{ .TypeName }})
		routeGroup.Put("/:id", hdr.Update{{ .TypeName }})
		routeGroup.Delete("/:id", hdr.Delete{{ .TypeName }})
	}
`
