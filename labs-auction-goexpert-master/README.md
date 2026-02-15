# Auction Service (Go Expert)

## Dev setup (Docker)

1) From the repository root, build and start the stack:

```bash
docker compose up --build
```

2) The API will be available at `http://localhost:8080`.

The application reads environment variables from [cmd/auction/.env](cmd/auction/.env) via the container. Adjust `AUCTION_INTERVAL`, `BATCH_INSERT_INTERVAL`, and `MAX_BATCH_SIZE` there if needed.

## Running tests

The auto-close test is an integration test and requires a running MongoDB instance.

1) Start MongoDB (you can reuse the `docker compose` stack).
2) Export `MONGODB_URL` to point at your MongoDB instance, for example:

```bash
export MONGODB_URL=mongodb://admin:admin@localhost:27017/auctions?authSource=admin
```

PowerShell:

```powershell
$env:MONGODB_URL = "mongodb://admin:admin@localhost:27017/auctions?authSource=admin"
```

3) Run tests:

```bash
go test ./...
```
