package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Didar1505/project_test.git/internal/auth"
	"github.com/Didar1505/project_test.git/internal/user"
	"github.com/Didar1505/project_test.git/pkg/config"
	"github.com/Didar1505/project_test.git/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"

	// "github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Application struct {
	r *gin.Engine
}

func NewApplication(r *gin.Engine) *Application {
	return &Application{
		r: r,
	}
}

func (a *Application) InitApp() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	dbPool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect")
	}

	sqlDB := stdlib.OpenDBFromPool(dbPool)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect gorm: ", err)
	}
	// logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	jwtMgr := auth.NewJWTManager(cfg.JWTSecret)

	userRepo := user.NewGormRepository(gormDB)
	otpRepo := auth.NewGormOTPRepository(gormDB)
	sessRepo := auth.NewGormSessionRepository(gormDB)
	mailer := auth.NewDevMailer()

	authService := auth.NewService(userRepo, otpRepo, sessRepo, mailer, jwtMgr)
	authHandler := auth.NewHandler(authService)
	
	api:= a.r.Group("/api")
	authHandler.RegisterRoutes(api)

	// protected:
	protected := api.Group("/")
	protected.Use(auth.Middleware(jwtMgr))

	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(protected)

	api.Use(user.FakeAuth())

	// static files
	// a.r.Static("/media", "./media")
}
