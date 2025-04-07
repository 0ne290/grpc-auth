# Для себя:

## Команда для генерации файлов .go из .proto: protoc -I grpc/proto grpc/proto/auth.proto --go_out=grpc/gen --go_opt=paths=source_relative --go-grpc_out=grpc/gen --go-grpc_opt=paths=source_relative

# Для всех:

## Как запускать:
1. Создайте файл .env в корне проекта. Пример содержимого файла:
```console
LOG_LEVEL=info

# gRPC
GRPC_ADDRESS=:1337

# Auth
AUTH_KEY=cryptographically_random_string_(the_longer_the_better)
AUTH_ACCESS_TOKEN_LIFETIME=1h
AUTH_REFRESH_TOKEN_LIFETIME=720h

# PostgreSQL
DB_HOST=postgres
DB_PORT=5432
DB_NAME=homework
DB_USER=test_user
DB_PASSWORD=securepassword
DB_SSL_MODE=disable
DB_POOL_MAX_CONNS=10
DB_POOL_MAX_CONN_LIFETIME=300s
DB_POOL_MAX_CONN_IDLE_TIME=150s
```
2. Выполните команду docker compose up.