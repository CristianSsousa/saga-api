# Saga API

A media tracker backend built with Go (Clean Architecture). Track movies, TV shows, games, books, and manga in one place.

## Prerequisites

- Go 1.22+
- PostgreSQL (or a [Neon](https://neon.tech) serverless Postgres instance)
- `swag` CLI for doc generation: `go install github.com/swaggo/swag/cmd/swag@latest`

## Local Setup

```bash
# 1. Clone
git clone git@github.com:CristianSsousa/saga-api.git
cd saga-api

# 2. Configure environment
cp .env.example .env
# Fill in DATABASE_URL, JWT_SECRET, TMDB_API_KEY, RAWG_API_KEY

# 3. Run migrations
make migrate
# Then execute each psql command shown against your database

# 4. (Optional) Generate Swagger docs
make swagger

# 5. Start the server
make run
```

## Environment Variables

| Variable       | Required | Description                              |
|----------------|----------|------------------------------------------|
| `DATABASE_URL` | Yes      | PostgreSQL connection string (pgx DSN)   |
| `JWT_SECRET`   | Yes      | Random secret for signing JWT tokens     |
| `TMDB_API_KEY` | No       | TMDB API key (movies & TV search)        |
| `RAWG_API_KEY` | No       | RAWG API key (game search)               |
| `PORT`         | No       | HTTP listen port (default: `8080`)       |

## API Endpoints

| Method | Path                      | Auth | Description                    |
|--------|---------------------------|------|--------------------------------|
| POST   | `/api/v1/auth/register`   | No   | Register a new user            |
| POST   | `/api/v1/auth/login`      | No   | Login and receive a JWT        |
| GET    | `/api/v1/auth/me`         | Yes  | Get current user profile       |
| GET    | `/api/v1/search`          | Yes  | Search media (`?q=&type=`)     |
| GET    | `/api/v1/library`         | Yes  | List library (cursor paginated)|
| POST   | `/api/v1/library`         | Yes  | Add media to library           |
| PUT    | `/api/v1/library/:id`     | Yes  | Update a library entry         |
| DELETE | `/api/v1/library/:id`     | Yes  | Remove a library entry         |

Media types accepted: `movie`, `tv`, `game`, `book`, `manga`

## Swagger UI

After running `make swagger` and starting the server:

```
http://localhost:8080/swagger/index.html
```

## Deploy to Cloud Run

```bash
# Build and push image
docker build -t gcr.io/YOUR_PROJECT/saga-api .
docker push gcr.io/YOUR_PROJECT/saga-api

# Deploy
gcloud run deploy saga-api \
  --image gcr.io/YOUR_PROJECT/saga-api \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars DATABASE_URL=...,JWT_SECRET=...,TMDB_API_KEY=...,RAWG_API_KEY=...
```
