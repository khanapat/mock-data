RELEASE?=$(shell git tag --points-at HEAD)
COMMIT?=$(shell git tag --points-at HEAD)
BUILD_TIME?=$(shell date '+%Y-%m-%d_%H:%M:%S')
PROJECT?=mock-data

GOOS?=linux
GOARCH?=amd64

APP?=goapp
PORT?=9090

CACHE_IMAGE?=$$(docker images --filter "dangling=true" -q --no-trunc)

run:
	go run main.go

clean:
	rm -f $(APP)

test: clean
	go test -v -cover ./...

create: test
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
    		-ldflags "-s -w -X $(PROJECT)/version.Release=$(RELEASE) \
    		-X $(PROJECT)/version.Commit=$(COMMIT) -X $(PROJECT)/version.BuildTime=$(BUILD_TIME)" \
    		-o $(APP)

container: create
	docker build . --no-cache -t $(PROJECT):$(RELEASE) -f build/app/Dockerfile

create-postgres:
	docker-compose -f ./build/postgres/docker-compose.yaml up -d

delete-postgres:
	docker-compose -f ./build/postgres/docker-compose.yaml down