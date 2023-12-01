# Web Devs Tools API Server

## Overview

This project is an API server for developer tools  built in Go, utilizing the chi router for routing, S3 for image storage, and PostgreSQL with Neon as the database. The backend is written in Go, providing functionality for user authentication, email sending, and tools management.

## Tech Stack

- Go (Golang)
- Chi Router
- S3 for image storage
- Resend for email sending
- PostgreSQL with Neon

## Routes

### Authentication

#### Register User with Magic Link
- `POST /v1/auth/magic-link`
- Body: `{"email": " "}`

#### Authenticate User with Magic Link
- `GET /v1/auth/magic-link/{token}`


#### Logout
- `DELETE /v1/auth/logout`

#### GitHub Login
- `GET /v1/auth/github/login`

#### GitHub Callback
- `GET /v1/auth/github/callback`

### Users

#### Get User
- `GET /v1/users/`


### Favorites

#### Add Favorite
- `POST /v1/favorites/`
- Body: `{"tool_id": " ". "user_id": " "}`

#### Get Favorites
- `GET /v1/favorites/`

#### Remove Favorite
- `DELETE /v1/favorites/{id}`

### Tools

#### Create Tool
- `POST /v1/tools/`
- Body: `{"name": " ", "description": " ", "website": " ", "imageUrl": " ", "category": " "}`

#### Get Tool
- `GET /v1/tools/{id}`

#### Delete Tool
- `DELETE /v1/tools/{id}`

#### Update Tool
- `PATCH /v1/tools/{id}`

#### Get All Tools
- `GET /v1/tools/`

#### Get Admin Tools
- `GET /v1/tools/admin`

#### Toggle Tool Published
- `GET /v1/tools/toggle-published/{id}`

### Categories

#### Create Category
- `POST /v1/categories/`
- Body: `{"name": " ", }`

#### Get All Categories
- `GET /v1/categories/`

#### Get Admin Categories
- `GET /v1/categories/admin`

#### Delete Category
- `DELETE /v1/categories/{id}`

#### Toggle Category Published
- `GET /v1/categories/toggle-published/{id}`

### Upload

#### Upload Image
- `POST /v1/upload/image`
- Body: `{"file": " "}`

### Health Check

#### Health Check
- `GET /v1/healthcheck`

Returns a JSON response with status "ok" to indicate the health of the API.

## Usage

1. Install dependencies: `go mod download`
2. Set up the database and configure S3, Resend credentials.
3. Run the server: `go run main.go`
