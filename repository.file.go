package blueprints

var RepositoryInterface = `
	package {{ .ApiName }}

	import (
		"gorm.io/gorm"
		entity "{{ .ModuleName }}/entities/{{ .ApiName }}"
	)

	type {{ .TypeName }}Repository interface {
		Get{{ .TypeName }}(id int) ([]entity.{{ .TypeName }}, error)
		Create{{ .TypeName }}(item *entity.{{ .TypeName }}) error
		Update{{ .TypeName }}(item entity.{{ .TypeName }}) error
		Delete{{ .TypeName }}(id int) error
	}

	type {{ .ApiName }}Repo struct {
		conn *gorm.DB
	}

	func New{{ .TypeName }}Repository(conn *gorm.DB) {{ .TypeName }}Repository {
		return &{{ .ApiName }}Repo{conn: conn}
	}
`

var RepositoryMethod = `
	package {{ .ApiName }}

	import (
		entity "{{ .ModuleName }}/entities/{{ .ApiName }}"
	)

	func (r *{{ .ApiName }}Repo) Get{{ .TypeName }}(id int) ([]entity.{{ .TypeName }}, error) {
		panic("unimplemented")
	}
	
	func (r *{{ .ApiName }}Repo) Create{{ .TypeName }}(item *entity.{{ .TypeName }}) error {
		panic("unimplemented")
	}

	func (r *{{ .ApiName }}Repo) Update{{ .TypeName }}(item entity.{{ .TypeName }}) error {
		panic("unimplemented")
	}

	func (r *{{ .ApiName }}Repo) Delete{{ .TypeName }}(id int) error {
		panic("unimplemented")
	}
`
