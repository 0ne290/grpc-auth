package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"grpc-auth/internal"
	core "grpc-auth/internal/core/auth"
	infrastructure2 "grpc-auth/internal/infrastructure"
	infrastructure "grpc-auth/internal/infrastructure/auth"
	web "grpc-auth/internal/web/auth"
	"grpc-auth/internal/web/interceptors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

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

	unitOfWork := infrastructure.NewPostgresUnitOfWork(pool)
	timeProvider := infrastructure2.NewRealTimeProvider()
	uuidProvider := infrastructure2.NewRealUuidProvider()
	hasher := infrastructure2.NewSha512Hasher()
	salter := infrastructure2.NewRealSalter()

	service := core.NewRealService(unitOfWork, timeProvider, uuidProvider, hasher, salter)

	controller := web.NewController(service)

	grpcServer := BuildGrpc(controller, logger)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 1337))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		logger.Infof("Starting server on %d", 1337)
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down gracefully...")
	grpcServer.GracefulStop()
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

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func BuildGrpc(controller *web.Controller, logger *zap.SugaredLogger) *grpc.Server {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.ErrorHandlingAndLogging(logger)))

	//grpcServer := grpc.NewServer()

	web.RegisterController(grpcServer, controller)

	return grpcServer
}
