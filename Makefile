SHELL := /bin/bash

PROJECT := $(shell basename $(PWD))

.PHONY: build
# build
build:
	CGO_ENABLED=0 go build -o $(PROJECT) .

.PHONY: dev
# build and run service foreground
dev:
	go run . run

.PHONY: run
# run service background
run:
	nohup ./${PROJECT} run > nohup.out 2>&1 &

.PHONY: stop
# stop
stop:
	@if [[ `ps -ef |grep -v 'grep' |grep ${PROJECT} |wc -l` != "0" ]]; then \
		kill `ps -ef |grep -v 'grep' |grep ${PROJECT} | awk -F ' ' '{print $$2}'`; \
		echo 'stopped'; \
	else \
		echo 'not running'; \
	fi

.PHONY: start
# start
start:
	make run
	@for (( i=0; i<5; i=i+1 )); do \
		if [[ `ps -ef |grep -v 'grep' |grep ${PROJECT}|wc -l` == "1" ]]; then \
			echo 'started'; \
			break; \
		else \
			echo 'retrying'; \
			sleep 1; \
		fi; \
	done

.PHONY: restart
# restart
restart:
	make stop
	make start

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
