build:
	@go mod download
	@go mod tidy
	@CGO_ENABLED=0 go build -ldflags="-s -w" -o tfunlock ./cmd/main.go