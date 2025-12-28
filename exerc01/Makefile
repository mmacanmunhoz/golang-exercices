APP_NAME=configparser
BIN_DIR=bin
SRC=./main.go

CONFIG_FILE=config.yaml


.PHONY: all build test run clean docker-build docker-run
all: build
build:
mkdir -p $(BIN_DIR)
go build -o $(BIN_DIR)/$(APP_NAME) $(SRC)
test:
go test ./...
run: build
./$(BIN_DIR)/$(APP_NAME) parse --file $(CONFIG_FILE)
clean:
rm -rf $(BIN_DIR)
docker-build:
docker build -t $(APP_NAME) .
docker-run:
docker run --rm -v $(PWD):/app $(APP_NAME) parse --file
/app/$(CONFIG_FILE)