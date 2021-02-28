PROJECTNAME := $(shell basename "$(PWD)")
all : clean fmt test build run

build:
	@echo " > Building csv import service..."
	@go build $(LDFLAGS) -o csv-importer

test-geolocation:
	go test github.com/goku321/geolocation/geolocation -v

test-e2e:
	@docker-compose up -d postgres
	go test -build=integration -v
	docker-compose down

test-store:
	docker run --name postgres -e POSTGRES_PASSWORD=postgres -d postgres
	go test github.com/goku321/geolocation/store -v
	docker stop postgres && docker rm postgres

test: test-geolocation test-store

fmt:
	go fmt ./...