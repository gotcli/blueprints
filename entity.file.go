package blueprints

var EntityFile = `
	package {{ .ApiName }}

	import "time"

	type {{ .TypeName }} struct {
		ID        int
		CreatedAt time.Time
	}
`
