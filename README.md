## Как запускать:
1. Создайте файл .env в корне проекта. Пример содержимого файла:
```console
LOG_LEVEL=info

LISTEN_ADDRESS=:8080
WRITE_TIMEOUT=15s
SERVER_NAME=homework-4

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

## Приложение запустится по адресу http://localhost:8080. Спецификацию можно посмотреть по адресу http://localhost:8080/swagger.

## Все эндпоинты v1/tasks защищены Bearer-авторизацией. Формат хэдера следующий: "Authorization: Bearer <UUID пользователя>".
## Т. е. чтобы работать с тасками, сперва надо создать юзера и использовать его UUID в качестве токена Bearer-авторизации. С тасками может взаимодействовать только тот юзер, который их создал - иначе 401.