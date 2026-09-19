package blueprints

var HandlerInterface = `
	package {{ .ApiName }}

	import (
		service "{{ .ModuleName }}/services/{{ .ApiName }}"
		"github.com/gofiber/fiber/v2"
	)

	type {{ .TypeName }}Handler interface {
		Get{{ .TypeName }}(c *fiber.Ctx) error
		Create{{ .TypeName }}(c *fiber.Ctx) error
		Update{{ .TypeName }}(c *fiber.Ctx) error
		Delete{{ .TypeName }}(c *fiber.Ctx) error
	}

	type {{ .ApiName }}Handler struct {
		srv service.{{ .TypeName }}Service
	}

	func New{{ .TypeName }}Handler(srv service.{{ .TypeName }}Service) {{ .TypeName }}Handler {
		return &{{ .ApiName }}Handler{srv: srv}
	}
`

var HandlerMethod = `
	package {{ .ApiName }}

	import (
		"github.com/gofiber/fiber/v2"
		"{{ .ModuleName }}/pkg/response"
	)

	func (s *{{ .ApiName }}Handler) Get{{ .TypeName }}(c *fiber.Ctx) error {
		return response.Send(c, fiber.StatusNotImplemented, nil)
	}
	
	func (s *{{ .ApiName }}Handler) Create{{ .TypeName }}(c *fiber.Ctx) error {
		return response.Send(c, fiber.StatusNotImplemented, nil)
	}

	func (s *{{ .ApiName }}Handler) Update{{ .TypeName }}(c *fiber.Ctx) error {
		return response.Send(c, fiber.StatusNotImplemented, nil)
	}

	func (s *{{ .ApiName }}Handler) Delete{{ .TypeName }}(c *fiber.Ctx) error {
		return response.Send(c, fiber.StatusNotImplemented, nil)
	}
`
