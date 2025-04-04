package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title Homework 4 API
// @version 1.0
// @description Task CRUD
// @host localhost:8080
// @BasePath /
func main() {
	var cfg internal.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	logger, err := NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := NewPostgresConnectionPool(context.Background(), cfg.PostgreSQL)
	if err != nil {
		log.Fatal(err)
	}

	taskUnitOfWork := taskInfrastructure.NewPostgresUnitOfWork(pool)
	userUnitOfWork := userInfrastructure.NewPostgresUnitOfWork(pool)

	timeProvider := core.NewRealTimeProvider()
	uuidProvider := core.NewRealUuidProvider()

	taskService := taskCore.NewService(taskUnitOfWork, timeProvider, uuidProvider)
	userService := userCore.NewService(userUnitOfWork, timeProvider, uuidProvider)

	taskController := taskWeb.NewController(taskService)
	userController := userWeb.NewController(userService)

	app := BuildRouting("http://"+cfg.Rest.ListenAddress, taskController, userController, logger)

	// Запуск HTTP-сервера в отдельной горутине
	go func() {
		logger.Infof("Starting server on %s", cfg.Rest.ListenAddress)
		if err = app.Listen(cfg.Rest.ListenAddress); err != nil {
			log.Fatal(err)
		}
	}()

	// Ожидание системных сигналов для корректного завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down gracefully...")
}

func NewLogger(level string) (*zap.SugaredLogger, error) {
	logLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	logger, err := zap.Config{
		Level:       logLevel,
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			TimeKey:    "timestamp",
			EncodeTime: zapcore.RFC3339NanoTimeEncoder,
		},
		DisableStacktrace: true,
	}.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func NewPostgresConnectionPool(ctx context.Context, cfg internal.PostgreSQL) (*pgxpool.Pool, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	// Парсим конфигурацию подключения
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func BuildRouting(allowOrigins string, taskController *taskWeb.Controller, userController *userWeb.Controller, logger *zap.SugaredLogger) *fiber.App {
	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Use(cors.New(cors.Config{
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-SomeID",
		ExposeHeaders:    "Link",
		AllowCredentials: true,
		AllowOrigins:     allowOrigins,
		MaxAge:           300,
	}))

	apiGroup := app.Group("/v1")
	taskApiGroup := apiGroup.Group("/tasks")
	userApiGroup := apiGroup.Group("/users")

	taskApiGroup.Post("", middlewares.ErrorHandlingAndLogging(logger), middlewares.Authentication(), taskController.Create)
	taskApiGroup.Get("", middlewares.ErrorHandlingAndLogging(logger), middlewares.Authentication(), taskController.GetAll)
	taskApiGroup.Get("/:uuid<guid>", middlewares.ErrorHandlingAndLogging(logger), middlewares.Authentication(), taskController.GetByUuid)
	taskApiGroup.Put("/:uuid<guid>", middlewares.ErrorHandlingAndLogging(logger), middlewares.Authentication(), taskController.Update)
	taskApiGroup.Delete("/:uuid<guid>", middlewares.ErrorHandlingAndLogging(logger), middlewares.Authentication(), taskController.Delete)

	userApiGroup.Post("", middlewares.ErrorHandlingAndLogging(logger), userController.Create)

	return app
}
