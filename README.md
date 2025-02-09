# RSSagg

RSSagg is a Go-based RSS aggregator that scrapes RSS feeds and stores the data in a PostgreSQL database. It supports user authentication via API keys and allows users to follow feeds and retrieve posts.

## Features

- User management with API key authentication
- Feed management
- Follow/unfollow feeds
- Scrape RSS feeds and store posts
- Retrieve posts for a user

## Project Structure

```
RSSagg/
├── auth/
│   └── auth.go
├── internal/
│   └── database/
│       ├── db.go
│       ├── feed_follows.sql.go
│       ├── feeds.sql.go
│       ├── models.go
│       ├── posts.sql.go
│       └── users.sql.go
├── sql/
│   ├── queries/
│   │   ├── feed_follows.sql
│   │   ├── feeds.sql
│   │   ├── posts.sql
│   │   └── users.sql
│   └── schema/
│       ├── 001_users.sql
│       ├── 002_users_api.sql
│       ├── 003_feeds.sql
│       ├── 004_feed_follows.sql
│       ├── 005_feeds_lastfetchedat.sql
│       └── 006_posts.sql
├── .gitignore
├── go.mod
├── go.sum
├── handler_error.go
├── handler_feed.go
├── handler_feed_follows.go
├── handler_readiness.go
├── handler_user.go
├── json.go
├── LICENSE
├── main.go
├── middleware_auth.go
├── models.go
├── rss.go
└── scrapper.go
```

## Getting Started

### Prerequisites

- Go 1.23.6 or later
- PostgreSQL
- Git

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/CLinnSheng/RSSagg.git
    cd RSSagg
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    go mod download
    ```

3. Set up the environment variables:
    ```sh
    cp .env.example .env
    # Edit .env to set your database URL and other configurations
    ```

4. Run the database migrations:
    ```sh
    goose $DB_URL up
    ```

### Running the Application

1. Start the application:
    ```sh
    go run main.go
    ```

2. The server will start on the port specified in the `.env` file.

### API Endpoints

- **Health Check**
  - `GET /v1/healthz`

- **User Management**
  - `POST /v1/users`
  - `GET /v1/users`

- **Feed Management**
  - `POST /v1/feed`
  - `GET /v1/feed`

- **Feed Follows**
  - `POST /v1/feed_follows`
  - `GET /v1/feed_follows`
  - `DELETE /v1/feed_follows/{feedFollowID}`

- **Posts**
  - `GET /v1/posts`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.