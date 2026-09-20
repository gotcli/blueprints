package blueprints

var MicroserviceMain = `package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

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
	if err := configs.InitViper(); err != nil { log.Fatalf("failed to initialize config: %v", err) }
	loggers, err := logger.Setup(viper.GetString("LOG_DIR"))
	if err != nil { log.Fatalf("failed to initialize logging: %v", err) }
	defer loggers.Close()
	conn := databases.ConnectDB()

	app := fiber.New(fiber.Config{
		ReadBufferSize: 8192, BodyLimit:viper.GetInt("HTTP_BODY_LIMIT"), ErrorHandler: response.ErrorHandler,
		ReadTimeout:viper.GetDuration("HTTP_READ_TIMEOUT"), WriteTimeout:viper.GetDuration("HTTP_WRITE_TIMEOUT"),
		IdleTimeout:viper.GetDuration("HTTP_IDLE_TIMEOUT"),
	})
	app.Use(middlewares.RequestLogger(loggers))
	routes.RegisterHealthRoutes(conn,app)
	routes.RegisterRootRoutes(conn,app)

	serverError:=make(chan error,1)
	go func(){ serverError<-app.Listen(fmt.Sprintf(":%d",viper.GetInt("APP_PORT"))) }()
	shutdown,stop:=signal.NotifyContext(context.Background(),syscall.SIGINT,syscall.SIGTERM)
	defer stop()
	select { case <-shutdown.Done(): case err:=<-serverError: if err!=nil { log.Printf("server stopped: %v",err) } }
	ctx,cancel:=context.WithTimeout(context.Background(),viper.GetDuration("SHUTDOWN_TIMEOUT")); defer cancel()
	if err:=app.ShutdownWithContext(ctx); err!=nil { log.Printf("graceful shutdown failed: %v",err) }
	if sqlDB,err:=conn.DB(); err==nil { _=sqlDB.Close() }
}
`

var MicroserviceHealth = `package routes

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"{{ .ModuleName }}/pkg/response"
)

func RegisterHealthRoutes(db *gorm.DB, app *fiber.App) {
	app.Get("/health/live",func(c *fiber.Ctx) error { return response.Send(c,fiber.StatusOK,map[string]string{"state":"alive"}) })
	app.Get("/health/ready",func(c *fiber.Ctx) error {
		sqlDB,err:=db.DB(); if err!=nil { return fiber.NewError(fiber.StatusServiceUnavailable,"database unavailable") }
		ctx,cancel:=context.WithTimeout(c.UserContext(),2*time.Second); defer cancel()
		if err:=sqlDB.PingContext(ctx); err!=nil { return fiber.NewError(fiber.StatusServiceUnavailable,"database unavailable") }
		return response.Send(c,fiber.StatusOK,map[string]string{"state":"ready"})
	})
}
`

var MicroserviceDockerfile = `FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /service .

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /service /app/service
COPY config.yml /app/config.yml
USER nonroot:nonroot
EXPOSE 3000
ENTRYPOINT ["/app/service"]
`

var MicroserviceDockerIgnore = `.git
.github
bin
logs
uploads
*.log
.env
README.md
`
