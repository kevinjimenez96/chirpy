# chirpy

## Getting Started

### Prerequisites

- Go 1.18+
- PostgreSQL

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/chirpy.git
    cd chirpy
    ```

2. Install dependencies:
    ```sh
    go mod download
    ```

3. Set up your environment variables (see `.env` for examples).

4. Run database migrations (if needed):
    ```sh
    # Example using sqlc or your preferred migration tool
    ```

5. Start the server:
    ```sh
    go run main.go
    ```

## API Endpoints

- `POST /users` - Register a new user
- `POST /login` - Authenticate and receive a JWT
- `POST /chirps` - Post a new chirp (requires authentication)
- `GET /chirps` - Retrieve chirps
- `GET /healthz` - Health check endpoint
- `GET /metrics` - Metrics endpoint

## Testing

Run tests with:
```sh
go test ./internal/...