PROJECTNAME := $(shell basename "$(PWD)")

all : clean fmt test build run

build:
	@echo " > Building csv import service..."
	@go build $(LDFLAGS) -o csv-importer

test-geolocation:
	go test github.com/goku321/geolocation/geolocation -v

test-e2e:
	docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 6432:5432 postgres
	go test github.com/goku321/geolocation/e2e_test -v
	docker stop postgres && docker rm postgres

test-store:
	docker run --name postgres -e POSTGRES_PASSWORD=postgres -d -p 6432:5432 postgres
	go test github.com/goku321/geolocation/store -v
	docker stop postgres && docker rm postgres

test: test-geolocation test-store

fmt:
	go fmt ./...

start-postgres-container:
	docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 6432:5432 postgres

stop-postgres-container:
	docker stop postgres && docker rm postgres