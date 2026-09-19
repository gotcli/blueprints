package blueprints

var MainFile = `
package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"{{ .ModuleName }}/configs"
	"{{ .ModuleName }}/databases"
	"{{ .ModuleName }}/middlewares"
	"{{ .ModuleName }}/pkg/logger"
	"{{ .ModuleName }}/pkg/response"
	"{{ .ModuleName }}/routes"
)

func main() {
	if err := configs.InitViper(); err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}
	loggers, err := logger.Setup(viper.GetString("LOG_DIR"))
	if err != nil { log.Fatalf("failed to initialize logging: %v", err) }
	defer loggers.Close()

	conn := databases.ConnectDB()

	app := fiber.New(fiber.Config{
		ReadBufferSize: 8192,
		BodyLimit:      viper.GetInt("HTTP_BODY_LIMIT"),
		ErrorHandler:   response.ErrorHandler,
	})
	app.Use(middlewares.RequestLogger(loggers))

	routes.RegisterRootRoutes(conn, app)
	//routes.NotFoundRoute(app)

	if err := app.Listen(fmt.Sprintf(":%d", viper.GetInt("APP_PORT"))); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
`
