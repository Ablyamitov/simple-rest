# Simple REST API

Simple REST API — demo RESTfull service, writing via Golang by using [Chi](https://github.com/go-chi/chi) 

## Project structure
***[cmd/app/](https://github.com/Ablyamitov/simple-rest/cmd/app/)*** - main package for starting the app


***[config](https://github.com/Ablyamitov/simple-rest/config/)*** - configuration files


***[internal/app](https://github.com/Ablyamitov/simple-rest/internal/app/)*** - app functionality


***[internal/db](https://github.com/Ablyamitov/simple-rest/internal/store/)*** - database, redis and web storage functionality


***[migration](https://github.com/Ablyamitov/simple-rest/migrations/)*** - migration source


***[gen](https://github.com/Ablyamitov/simple-rest/gen/)*** - swagger generated files

## Technologies used
- Language - [go](https://go.dev/)
- Database - [postgresql](https://www.postgresql.org/)
- Caching - [redis](https://redis.io/?ref=kubedexcom)
- Migration - [migrate](https://github.com/golang-migrate/migrate)

## Installation
1. Clone :
    ```bash
    git clone https://github.com/Ablyamitov/simple-rest.git
    cd simple-rest
    ```

2. Download dependencies:
    ```bash
    go mod download
    ```

3. Launch Postgresql, Redis server, create database

4. Set up the configuration file `config.yaml`:
    ```yaml
    server:
      port: 8080

    database:
      url: "postgresql://postgres:12345678@localhost:5432/library?sslmode=disable"

    cache:
      redis_url: "redis://localhost:6379"
    ```
   
## Using via Makefile
- **Build project**:
    ```bash
    make build
    ```

- **Run server**:
    ```bash
    make run
    ```

- **Apply migrations**:
    ```bash
    make migrate
    ```

- **Create new migration**:
    ```bash
    make migrate-create name=<migration_name>
    ```
  Replace `<migration_name>` with the desired migration name.


- **Clean up**:
    ```bash
    make clean
    ```
  
## Deployment
You can use Docker to containerize your application. An example Dockerfile and deployment instructions are included in the project for your convenience.


