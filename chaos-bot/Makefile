RUN_OPTS ?= --dev --debug

include .env
export sed 's/=.*//' .env

run: install
	chaos-bot $(RUN_OPTS) discord-bot --discord-token="$(DISCORD_TOKEN)"

install:
	go install

docker:
	docker build -t ultreme/chaos-bot .
	docker run -it --rm ultreme/chaos-bot $(RUN_OPTS) discord-bot --discord-token="$(DISCORD_TOKEN)"

bumpdeps:
	go get -u ./...
	go mod tidy

tidy:
	go mod tidy
