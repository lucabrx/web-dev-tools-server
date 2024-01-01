# Web Dev Tools API 

## Project Environment Variables
```
DB_URL=postgres://postgres:postgres@localhost:5432/web_dev_tools
SERVER_ADDRESS=:8080
CLIENT_ADDRESS=http://localhost:3000
RESEND_API_KEY=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
AWS_ACCESS_KEY=
AWS_SECRET_KEY=
```

## Resources
- [Go](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)
- [AWS](https://aws.amazon.com/)
- [Resend](https://resend.com/)
- [migrate-tool](https://github.com/golang-migrate/migrate)

## Project outline
- users -> add tools to favorites, suggest tools
- tools -> paginated list of tools with search
- auth -> login with GitHub and magic link
- admin -> add tools to the database, approve suggested tools

## How to run
### Running PostgreSQL in Docker
```shell 
docker run --name web_dev_tools -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
### OR 
task postgres-container
```

### Migrating the database
```shell
migrate -path ./migrations -database postgres://postgres:postgres@localhost:5432/web_dev_tools || <pg db> up
### OR
task migrate-up ### DB_URL from .env or $DB_URL
```

### Running the server
```shell
go run ./cmd/api
### OR
task run-server
```