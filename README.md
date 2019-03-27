# Esender

HTTP (REST) and GRPC microservices for sending emails written on Go.

# Migrations

```bash
# up
cd ./boiler/migrations/mysql
goose mysql "user:password@/database" up
# 2019/03/01 23:22:36 OK    00001_init.sql

# down
goose mysql "user:password@/database" down
# 2019/03/01 23:22:36 OK    00001_init.sql
```