package blueprints

var SchemaFile = `
	package {{ .ApiName }}

	type (
		Get{{ .TypeName }}Request struct {
			ID int
		}

		Create{{ .TypeName }}Request struct {
			ID int
			Name string
		}

		Update{{ .TypeName }}Request struct {
			ID int
			Name string
		}

		Delete{{ .TypeName }}Request struct {
			ID int
		}
	)

	type (
		Get{{ .TypeName }}Response struct {
			ID int
			Name string
		}
	)
`
