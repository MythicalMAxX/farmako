# Farmako Backend

Backend service for the Farmako Pharmacy Coupon System built with Go.

## Environment Variables

Create a `.env` file in the backend directory with the following format:

```
DB_DRIVER=sqlite
DB_NAME=coupon_system.db
SERVER_PORT=0000
```

## API Documentation

The API is documented using OpenAPI/Swagger. Access the Swagger UI at:

- http://localhost:8080/api/v1/swagger/

## Running the Backend

### With Docker

```bash
docker-compose up -d backend
```

### Without Docker

```bash
cd backend
go mod download
go run cmd/api/main.go
``` 