package blueprints

var ServiceInterface = `
	package {{ .ApiName }}

	import (
		repository "{{ .ModuleName }}/repositories/{{ .ApiName }}"
		schema "{{ .ModuleName }}/schemas/{{ .ApiName }}"
	)

	type {{ .TypeName }}Service interface {
		Get{{ .TypeName }}(req schema.Get{{ .TypeName }}Request) (schema.Get{{ .TypeName }}Response, error)
		Create{{ .TypeName }}(req schema.Create{{ .TypeName }}Request) error
		Update{{ .TypeName }}(req schema.Update{{ .TypeName }}Request) error
		Delete{{ .TypeName }}(req schema.Delete{{ .TypeName }}Request) error
	}

	type {{ .ApiName }}Service struct {
		repo repository.{{ .TypeName }}Repository
	}

	func New{{ .TypeName }}Service(repo repository.{{ .TypeName }}Repository) {{ .TypeName }}Service {
		return &{{ .ApiName }}Service{repo: repo}
	}
`

var ServiceMethod = `
	package {{ .ApiName }}

	import (
		schema "{{ .ModuleName }}/schemas/{{ .ApiName }}"
	)

	func (s *{{ .ApiName }}Service) Get{{ .TypeName }}(req schema.Get{{ .TypeName }}Request) (schema.Get{{ .TypeName }}Response, error) {
		panic("unimplemented")
	}
	
	func (s *{{ .ApiName }}Service) Create{{ .TypeName }}(req schema.Create{{ .TypeName }}Request) error {
		panic("unimplemented")
	}

	func (s *{{ .ApiName }}Service) Update{{ .TypeName }}(req schema.Update{{ .TypeName }}Request) error {
		panic("unimplemented")
	}

	func (s *{{ .ApiName }}Service) Delete{{ .TypeName }}(req schema.Delete{{ .TypeName }}Request) error {
		panic("unimplemented")
	}
`
