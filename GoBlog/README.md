## GoBlog

A minimal Go-based personal blog system skeleton using **gin**, **PostgreSQL**, **gorm**, and **slog**.

### Quick start

1. Set the required environment variable for the database:

```bash
export DATABASE_DSN="postgres://user:pass@localhost:5432/goblog?sslmode=disable"
```

Optional environment variables:

- `APP_ENV`: `production` | `development` | `test` (default: `development`)
- `HTTP_PORT`: HTTP listen port (default: `8080`)

2. Run the service:

```bash
cd GoBlog
go run ./cmd/blog
```

3. Test endpoints:

- `GET http://localhost:8080/health` – health check
- `GET http://localhost:8080/` – home placeholder

