 # Go Auth Service

This project is a Go-based authentication service that uses PostgreSQL for database management and Goose for database migrations. The application is containerized using Docker.

## Features

- Authentication service
- PostgreSQL database
- Database migrations using Goose
- Dockerized for easy deployment

## Prerequisites

Before you begin, ensure you have met the following requirements:

- [Go](https://golang.org/dl/) installed
- [Docker](https://www.docker.com/get-started) installed
- [Docker Compose](https://docs.docker.com/compose/install/) installed

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/yourusername/go-auth-service.git
   cd go-auth-service

2. Install dependencies:
   ```sh
   go mod download

3. Rename file `.env.example` to `.env`:
   ```sh
   mv .env.sample .env

4. Generate JWT secret
5. Open .env and fill it with your values:
   ```sh
   DB_URL = "postgres://$DBUSER:$DBPASSWORD@$DBHOST/$DBNAME?sslmode=$DBSSL"
   JWT_SECRET = "your_jwt_token"
   DBHOST=host.docker.internal
   DBUSER=your_db_username
   DBPASSWORD=your_db_password
   DBNAME=your_db_name
   DBSSL=disable

## Database Migrations
This project uses Goose for database migrations. Migrations are stored in the sql/migrations directory.
1. Apply migrations:

   ```sh
   goose -dir ./sql/migrations postgres "user=db_username password=db_password dbname=db_name sslmode=disable" up
## Usage
1. Build the application:

   ```sh
   go build -o main cmd/main.go
2. Run the application:
   ```sh
   ./main

## Docker
The project includes a Dockerfile and docker-compose.yml for containerizing the application.

1. Build the Docker image:
    ```sh
    docker-compose build
2. Run the Docker container:
   ```sh
   docker-compose up
