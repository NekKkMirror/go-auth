package app

import (
	"context"
	handler "github.com/NekKkMirror/go-auth/internal/handler/auth"
	repository "github.com/NekKkMirror/go-auth/internal/repository/refresh-token"
	jwtserv "github.com/NekKkMirror/go-auth/internal/service/jwt"
	"log"

	"github.com/NekKkMirror/go-auth/config"
	"github.com/NekKkMirror/go-auth/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App wraps all necessary components for running the application.
type App struct {
	Config *config.Config
	DB     *pgxpool.Pool
	Server *fiber.App
}

// NewApp initializes and configures the application.
func NewApp() *App {
	app := &App{}

	app.Config = config.LoadConfig()
	app.DB = app.initDB()
	app.Server = fiber.New(fiber.Config{
		Prefork: false,
	})

	app.registerMiddleware()
	app.registerRoutes()

	return app
}

// initDB establishes a connection to the PostgreSQL database using pgxpool.
func (a *App) initDB() *pgxpool.Pool {
	dbUrl := a.Config.DBUrl
	dbpool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}
	log.Println("Database connected successfully")
	return dbpool
}

// registerMiddleware sets up global middleware for Fiber.
func (a *App) registerMiddleware() {
	a.Server.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	a.Server.Use(recover.New())
}

// registerRoutes sets up the routing for the application.
func (a *App) registerRoutes() {
	refreshTokenRepo := repository.NewRefreshTokenRepository(a.DB)

	jwtService := jwtserv.NewJWTService(a.Config.AppJWTSecret)
	smtpCfg := config.LoadSMTPConfig()
	mailer := service.NewSMTPMailer(smtpCfg)

	authService := service.NewAuthService(refreshTokenRepo, jwtService, mailer)
	authHandler := handler.NewAuthHandler(authService)

	api := a.Server.Group(a.Config.AppAPIBasePath)

	handler.RegisterAuthRoutes(api, authHandler)
}

// Run starts the Fiber server.
func (a *App) Run() error {
	log.Printf("Server is listening on port %s", a.Config.AppPort)
	if err := a.Server.Listen(":" + a.Config.AppPort); err != nil {
		return err
	}
	return nil
}

// Shutdown disconnects from the database pool and gracefully shuts down the server.
func (a *App) Shutdown(ctx context.Context) {
	if a.DB != nil {
		a.DB.Close()
		log.Println("Database connection closed")
	}

	if err := a.Server.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}
	log.Println("Server shutdown successfully")
}
