.PHONY: all build-server build-relay build-client build-frontend build-docker build-all test clean lint

VERSION ?= 1.0.0

all: build-server build-relay build-frontend

build-server:
	cd backend && go build -ldflags="-s -w -X main.Version=$(VERSION)" -o ../bin/cinaseek-server ./cmd/server/

build-relay:
	cd websocket/backend && go build -ldflags="-s -w -X main.Version=$(VERSION)" -o ../../bin/cinaseek-relay .

build-client:
	cd backend && go build -ldflags="-s -w -X main.Version=$(VERSION)" -o ../bin/cinaseek-client ./cmd/client/

build-frontend:
	cd frontend && pnpm install && pnpm build

build-docker:
	docker-compose -f build/docker/docker-compose.yml build

build-all:
	./build/build-all.sh $(VERSION)

test:
	cd backend && go test ./...
	cd websocket/backend && go test ./...

clean:
	rm -rf bin/ build/output/ frontend/dist/

lint:
	cd backend && go vet ./...
	cd websocket/backend && go vet ./...
