go-lint:
	golangci-lint run -E gofmt -E nakedret -E gochecknoglobals -E unconvert -E gocritic -E maligned -E prealloc -E gosec -E bodyclose -E exhaustive -E golint -E gocyclo
generate-doc:
	swag init -g cmd/app/main.go
download-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
