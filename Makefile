BUILD_DIR = ./bin

build:
	mkdir -p $(BUILD_DIR) #
	go mod tidy
	go build -o ./bin -v ./cmd/app

gen:
	protoc -I ./proto ./proto/object-storage.proto --go_out=. --go-grpc_out=.
	wire ./internal/app
