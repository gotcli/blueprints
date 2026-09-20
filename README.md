# GOT Blueprints

Public Community blueprints used by [GOT CLI](https://github.com/gotcli/got-community)
to generate Go/Fiber projects.

The package contains the standard, microservice, workspace-supporting, JWT,
upload, Swagger/OpenAPI, database, logging, and layered API source templates
used by GOT Community. Blueprints are compiled into the CLI so generation can
run without downloading templates at runtime.

## Requirements

- Go 1.25 or newer

## Validation

```sh
go fmt ./...
go vet ./...
go test ./...
go build ./...
```

## Versioning

GOT CLI pins a tagged version of this module. Changes that alter generated
output should be reviewed together with the corresponding CLI tests.

## License

Licensed under the [Apache License 2.0](LICENSE).
