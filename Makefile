.PHONY: build run clean start


APPLICATION_NAME := unnamed

all: build run

build: clean
	@echo "Building..."
	@mkdir -p bin
	@go build -o bin/${APPLICATION_NAME} ./cmd/${APPLICATION_NAME}/main.go

run:
	@echo "Running..."
	@./bin/${APPLICATION_NAME}

clean:
	@echo "Cleaning up..."
	@mkdir -p bin
	@rm -rf bin/*


