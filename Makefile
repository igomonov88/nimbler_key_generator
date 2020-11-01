SHELL := /bin/bash

all:

proto:
	protoc --go_out=plugins=grpc:. nimbler_key_generator/proto/contract.proto

migrate:
	go run ./cmd/admin/main.go migrate

seed: migrate
	go run ./cmd/admin/main.go seed

up:
	docker-compose up

down:
	docker-compose down

test:
	go test ./... -cover

clean:
	docker system prune -f

stop-all:
	docker stop $(docker ps -aq)

remove-all:
	docker rm $(docker ps -aq)

tidy:
	go mod tidy

deps-reset:
	git checkout -- go.mod
	go mod tidy

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -t -d -v ./...

deps-cleancache:
	go clean -modcac