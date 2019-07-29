# Esender

HTTP (REST API based on [echo](https://github.com/labstack/echo)) and GRPC microservice (based on [go-kit](https://github.com/go-kit/kit)) for sending emails.

## Migrations

migration tool - [goose](https://github.com/pressly/goose)

```bash
# up
cd ./boiler/migrations/mysql
goose mysql "user:password@/database" up
# 2019/03/01 23:22:36 OK    00001_init.sql

# down
goose mysql "user:password@/database" down
# 2019/03/01 23:22:36 OK    00001_init.sql
```