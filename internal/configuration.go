package internal

import "time"

type AppConfig struct {
	LogLevel   string `envconfig:"LOG_LEVEL" required:"true"`
	GrpcAdress string `envconfig:"GRPC_ADDRESS" required:"true"`
	Auth       AuthConfig
	PostgreSQL PostgreSqlConfig
}

type AuthConfig struct {
	Key                  string        `envconfig:"AUTH_KEY" required:"true"`
	AccessTokenLifetime  time.Duration `envconfig:"AUTH_ACCESS_TOKEN_LIFETIME" required:"true"`
	RefreshTokenLifetime time.Duration `envconfig:"AUTH_REFRESH_TOKEN_LIFETIME" required:"true"`
}

type PostgreSqlConfig struct {
	Host                string        `envconfig:"DB_HOST" required:"true"`
	Port                int           `envconfig:"DB_PORT" required:"true"`
	Name                string        `envconfig:"DB_NAME" required:"true"`
	User                string        `envconfig:"DB_USER" required:"true"`
	Password            string        `envconfig:"DB_PASSWORD" required:"true"`
	SSLMode             string        `envconfig:"DB_SSL_MODE" default:"disable"`
	PoolMaxConns        int           `envconfig:"DB_POOL_MAX_CONNS" default:"5"`
	PoolMaxConnLifetime time.Duration `envconfig:"DB_POOL_MAX_CONN_LIFETIME" default:"180s"`
	PoolMaxConnIdleTime time.Duration `envconfig:"DB_POOL_MAX_CONN_IDLE_TIME" default:"100s"`
}
