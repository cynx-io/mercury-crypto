tidy:
	go mod tidy
	go fmt ./...
	fieldalignment -fix ./...
	go vet ./...
	golangci-lint run --fix ./...

run:
	make tidy
	go run main.go

build:
	make tidy
	go build

install_deps:
	# These needs sudo
	# apt install build-essential -y
    # curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	go install github.com/google/wire/cmd/wire@latest
	go get -u gorm.io/gorm
	go get -u gorm.io/driver/sqlite

.PHONY: proto
proto:
	@echo "Generating proto files..."
	protoc --experimental_allow_proto3_optional \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		api/proto/crypto/crypto.proto
	@echo "Proto files generated successfully!"

.PHONY: build
build: proto
	@echo "Building application..."
	go build -o bin/mercury main.go

.PHONY: run
run:clean proto
	@echo "Running application..."
	go run main.go

.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	rm -f api/proto/user/*.pb.go
	rm -f bin/mercury

.PHONY: all
all: clean proto build

build_docker_dev:
	docker build -t mercury-crypto-dev:latest .
	docker tag mercury-crypto-dev:latest derwin334/mercury-crypto-dev:latest
	docker push derwin334/mercury-crypto-dev:latest

