.PHONY: build run clean test tidy

APP_NAME=simple_im
BUILD_DIR=./build

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd

run:
	go run ./cmd -c config -cPath "./,./configs/"

clean:
	rm -rf $(BUILD_DIR)
	rm -rf ./logs
	rm -rf ./uploads/*

test:
	go test -v ./...

tidy:
	go mod tidy

migrate:
	@echo "Run migrations manually with psql"
	@echo "psql -U postgres -d simple_im -f migrations/001_init.sql"
