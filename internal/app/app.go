package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Didar1505/project_test.git/internal/auth"
	"github.com/Didar1505/project_test.git/internal/auth/providers/jwt"
	"github.com/Didar1505/project_test.git/internal/auth/providers/oauth"
	"github.com/Didar1505/project_test.git/internal/auth/providers/otp"
	"github.com/Didar1505/project_test.git/internal/auth/session"
	"github.com/Didar1505/project_test.git/internal/course/handler"
	"github.com/Didar1505/project_test.git/internal/course/repo"
	"github.com/Didar1505/project_test.git/internal/course/service"
	"github.com/Didar1505/project_test.git/internal/mailer"
	"github.com/Didar1505/project_test.git/internal/user"
	"github.com/Didar1505/project_test.git/pkg/config"
	"github.com/Didar1505/project_test.git/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"

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

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	jwtMgr := jwt.NewManager(cfg.JWTSecret)

	userRepo := user.NewGormRepository(gormDB)
	otpRepo := otp.NewGormRepository(gormDB)
	sessRepo := session.NewGormSessionRepository(gormDB)
	smtpMailer := mailer.New(cfg)

	authService := auth.NewService(userRepo, otpRepo, sessRepo, smtpMailer, jwtMgr)
	authHandler := auth.NewHandler(authService)

	api := a.r.Group("/api")

	if cfg.GoogleOAuthCredentials != "" {
		if cfg.GoogleOAuthCookieSecret == "" {
			log.Fatal("GOOGLE_OAUTH_COOKIE_SECRET is required for Google OAuth")
		}
		oauth.Setup(
			cfg.GoogleOAuthRedirectURL,
			cfg.GoogleOAuthCredentials,
			[]string{"email", "profile"},
			[]byte(cfg.GoogleOAuthCookieSecret),
		)
		api.Use(oauth.Session("ginoauth_google_session"))
	}

	authHandler.RegisterRoutes(api)

	if cfg.GoogleOAuthCredentials != "" {
		oauthGroup := api.Group("/auth/google")
		authHandler.RegisterOAuthRoutes(oauthGroup)
	}

	// protected:
	protected := api.Group("/")
	protected.Use(auth.Middleware(jwtMgr))

	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(protected)

	// COURSES
	courseRepo := repo.NewCourseRepository(gormDB, &logger)
	courseService := service.NewCourseService(*courseRepo)
	courseHandler := handler.NewCourseHandler(courseService)
	courseHandler.RegisterRoutes(api)

	// MODULES
	moduleRepo := repo.NewModuleRepository(gormDB, &logger)
	moduleService := service.NewModuleService(*moduleRepo)
	moduleHandler := handler.NewModuleHandler(moduleService)
	moduleHandler.RegisterRoutes(api)

	// LESSONS
	lessonRepo := repo.NewLessonRepository(gormDB, &logger)
	lessonService := service.NewLessonService(*lessonRepo)
	lessonHandler := handler.NewLessonHandler(lessonService)
	lessonHandler.RegisterRoutes(api)

	// static files
	// a.r.Static("/media", "./media")
}
