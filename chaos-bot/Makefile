RUN_OPTS ?= --dev --debug

include .env
export sed 's/=.*//' .env

COMPILEDAEMON_OPTIONS ?= -exclude-dir=.git -color=true -build=go\ install -build-dir=.
run:
	go install github.com/githubnemo/CompileDaemon
	CompileDaemon $(COMPILEDAEMON_OPTIONS) -command="chaos-bot $(RUN_OPTS) discord-bot -discord-token=$(DISCORD_TOKEN) -manfred-channel=$(MANFRED_CHANNEL)"
.PHONY: run

install-run: install
	chaos-bot $(RUN_OPTS) discord-bot --discord-token="$(DISCORD_TOKEN)"

install:
	go install

docker:
	docker build -t ultreme/chaos-bot .
	docker run -it --rm ultreme/chaos-bot $(RUN_OPTS) discord-bot --discord-token="$(DISCORD_TOKEN)"

bumpdeps:
	go get -u ./...
	go get moul.io/pipotron@master
	go mod tidy

tidy:
	go mod tidy
