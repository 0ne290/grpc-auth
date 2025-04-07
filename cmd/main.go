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
	core "grpc-auth/internal/core/services/auth"
	"grpc-auth/internal/infrastructure"
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

	unitOfWorkStarter := infrastructure.NewPostgresUnitOfWorkStarter(pool)
	timeProvider := infrastructure.NewRealTimeProvider()
	uuidProvider := infrastructure.NewRealUuidProvider()
	hasher := infrastructure.NewSha512Hasher()
	salter := infrastructure.NewRealSalter()
	jwtManager := infrastructure.NewRealJwtManager([]byte(cfg.Auth.Key))

	service := core.NewRealService(cfg.Auth.AccessTokenLifetime, cfg.Auth.RefreshTokenLifetime, unitOfWorkStarter, timeProvider, uuidProvider, hasher, salter, jwtManager)

	controller := web.NewController(service)

	grpcServer := BuildGrpc(controller, logger)

	lis, err := net.Listen("tcp", cfg.GrpcAdress)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		logger.Info("Starting server on ", cfg.GrpcAdress)
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

func NewPostgresConnectionPool(ctx context.Context, cfg internal.PostgreSqlConfig) (*pgxpool.Pool, error) {
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

	web.RegisterController(grpcServer, controller)

	return grpcServer
}
