package main

import (
	"fmt"
	"log"
	"mock-data/internal/database"
	"mock-data/logz"
	"mock-data/middleware"
	"mock-data/user"
	"mock-data/version"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	_ "mock-data/docs"
	_ "time/tzdata"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func init() {
	runtime.GOMAXPROCS(1)
	initTimezone()
	initViper()
}

var isReady bool

// @title Mock API
// @version 1.0
// @description Mock api with fiber framework.
// @termsOfService http://swagger.io/terms/
// @contact.name K.apiwattanawong
// @contact.url http://www.swagger.io/support
// @contact.email k.apiwattanawong@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9090
// @BasePath /mock
// @schemes http https
func main() {
	timeout := viper.GetDuration("app.timeout")

	app := fiber.New(fiber.Config{
		StrictRouting: true,
		CaseSensitive: true,
		Immutable:     true,
		ReadTimeout:   timeout,
		WriteTimeout:  timeout,
		IdleTimeout:   timeout,
	})

	logger, err := logz.NewLogConfig()
	if err != nil {
		log.Fatal(err)
	}

	postgresDB, err := database.NewPostgresConn()
	if err != nil {
		logger.Error(err.Error())
	}
	defer postgresDB.Close()

	if err := postgresDB.Ping(); err != nil {
		logger.Error(err.Error())
	}

	middle := middleware.NewMiddleware(logger)

	app.Use(middle.CorsMiddleware())

	swag := app.Group("/swagger")
	swag.Use(middle.BasicAuthenication())
	registerSwaggerRoute(swag)

	api := app.Group(viper.GetString("app.context"))

	api.Use(middle.JSONMiddleware())
	api.Use(middle.ContextLocaleMiddleware())
	api.Use(middle.LoggingMiddleware())

	userHandler := user.NewUserHandler(
		user.NewUserRepositoryDB(postgresDB),
	)

	api.Get("/user", userHandler.GetUser)

	app.Get("/version", version.VersionHandler)
	app.Get("/liveness", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Get("/readiness", func(c *fiber.Ctx) error {
		if isReady {
			return c.SendStatus(fiber.StatusOK)
		}
		return c.SendStatus(fiber.StatusServiceUnavailable)
	})

	logger.Info(fmt.Sprintf("⇨ http server started on [::]:%s", viper.GetString("app.port")))

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", viper.GetString("app.port"))); err != nil {
			logger.Info(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	isReady = true

	select {
	case <-c:
		logger.Info("terminating: by signal")
	}

	app.Shutdown()

	logger.Info("shutting down")
	os.Exit(0)
}

func initViper() {
	viper.SetDefault("app.name", "mock-data")
	viper.SetDefault("app.port", "9090")
	viper.SetDefault("app.timeout", "60s")
	viper.SetDefault("app.context", "/mock")

	viper.SetDefault("swagger.host", "localhost:9090")
	viper.SetDefault("swagger.user", "admin")
	viper.SetDefault("swagger.password", "password")

	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.env", "dev")

	viper.SetDefault("postgres.type", "postgres")
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", "5432")
	viper.SetDefault("postgres.username", "postgres")
	viper.SetDefault("postgres.password", "P@ssw0rd")
	viper.SetDefault("postgres.database", "test")
	viper.SetDefault("postgres.timeout", 100)
	viper.SetDefault("postgres.sslmode", "disable")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func initTimezone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("error loading location 'Asia/Bangkok': %v\n", err)
	}
	time.Local = ict
}

func registerSwaggerRoute(swag fiber.Router) {
	swag.Get("/*", swagger.New(swagger.Config{
		URL:         fmt.Sprintf("http://%s/swagger/doc.json", viper.GetString("swagger.host")),
		DeepLinking: false,
	}))
}
