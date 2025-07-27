# RSS Aggregator

This project is an RSS aggregator web application written in Go. It allows users to follow RSS feeds, fetch and store posts, and provides a RESTful API for interacting with feeds and posts.

## Features

- User registration and API key authentication
- Add and follow RSS feeds
- Periodic scraping of RSS feeds for new posts
- Store posts and feed data in PostgreSQL
- RESTful API endpoints for feeds, posts, and user management
- Health check endpoint

## Project Structure

- `cmd/main.go`: Application entry point
- `internal/`: Main application logic
  - `app/`: App struct and core logic
  - `auth/`: Authentication and API key management
  - `db/`: Database models and queries (uses sqlc)
  - `handler/`: HTTP handlers for API endpoints
  - `middleware/`: HTTP middleware (e.g., auth)
  - `response/`: Response formatting utilities
  - `scraper/`: RSS feed scraping logic
  - `server/`: HTTP server setup
- `sql/`: SQL schema and queries
- `vendor/`: Third-party dependencies

## Requirements

- Go 1.24+
- PostgreSQL

## Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/serhappy/rssagg.git
   cd rssagg
   ```

2. Copy `.env.example` to `.env` and set your environment variables:
   - `PORT`: Port to run the server (e.g., 8080)
   - `DB_URL`: PostgreSQL connection string
3. Run database migrations (see `sql/schema/` for SQL files):

    ```sh
    cd sql/schema
    goose postgres $DB_URL up
    ```

4. Build and run the application:

   ```sh
   make build
   ./build/server
   ```

   Or run directly:

   ```sh
   go run cmd/main.go
   ```

## API Endpoints

- `POST /users`: Register a new user
- `POST /feeds`: Add a new RSS feed
- `GET /feeds`: List all feeds
- `POST /feed_follows`: Follow a feed
- `GET /posts`: Get posts for followed feeds
- `GET /health`: Health check

Authentication is via API key in the `Authorization` header.
