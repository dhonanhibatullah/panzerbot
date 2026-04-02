.PHONY: swagger build

swagger:
	@swag init -g cmd/core/main.go -o docs/swagger

build: swagger
	@rm -rf build
	@mkdir -p build
	@go build -o build/panzerbot cmd/core/main.go