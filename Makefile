include .env

OAPIGEN=$(shell go env GOPATH)/bin/oapi-codegen
GOOSE=$(shell go env GOPATH)/bin/goose
SQLC=$(shell go env GOPATH)/bin/sqlc

APISPECLOCATION=internal/api/spec
apigen:
	${OAPIGEN} --config=configs/.oapi.server.yaml ${APISPECLOCATION}/openapi.yaml
	${OAPIGEN} --config=configs/.oapi.types.yaml ${APISPECLOCATION}/openapi.yaml
	${OAPIGEN} --config=configs/.oapi.spec.yaml ${APISPECLOCATION}/openapi.yaml

dbgen:
	${SQLC} generate -f configs/.sqlc.yaml

dbmigration-up:
	${GOOSE} -dir sql/migrations ${DATABASE_DRIVER} ${DATABASE_URL} up

build:
	/usr/local/go/bin/go mod tidy
	/usr/local/go/bin/go build

postgres-mock-up:
	docker run --rm --name postgres15 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=school -p 5432:5432 -d postgres:15-alpine

dockerc-upb:
	docker compose -f deployments/docker-compose.yaml up --build