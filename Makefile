PROJECT_DIR := ${CURDIR}

## build: build docker image
build:
	docker build --no-cache -t park .

## rm: rm docker image
rm:
	docker image rm $$(docker images park -q)


## run: run docker container
run:
	docker run -d --rm --memory 2G --log-opt max-size=5M --log-opt max-file=3 --name park_perf -p 5000:5000 park

## stop: stop docker container
stop:
	docker stop $$(docker ps -a -q --filter name=park_perf)


## run-postgres: run postgres docker container
run-postgres:
	docker run -d --rm -p 5432:5432 --name postgres park

## stop-postgres: stop postgres docker container
stop-postgres:
	docker stop $$(docker ps -a -q --filter name=postgres)

## rm-postgres: rm postgres docker container
rm-postgres:
	docker rm postgres


## run-go: run go docker container
run-go:
	docker run -d --rm -p 5000:5000 --name golang park

## stop-go: stop go docker container
stop-go:
	docker stop $$(docker ps -a -q --filter name=golang)

## rm-go: rm go docker container
rm-go:
	docker rm golang


## test-func: run functional test
test-func:
	./tech-db-forum func -u http://localhost:5000/api -r report.html -k

## test-fill: run fill db test
test-fill:
	./tech-db-forum fill --url=http://localhost:5000/api --timeout=900

## test-perf: run performance test
test-perf:
	./tech-db-forum perf --url=http://localhost:5000/api --duration=600 --step=60

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command to run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo