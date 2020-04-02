include .env
export sed 's/=.*//' .env

install:
	go install

run: install
	chaos-bot -t $(DISCORD_TOKEN)

docker:
	docker build -t ultreme/chaos-bot .
	docker run -it --rm ultreme/chaos-bot -t $(DISCORD_TOKEN)
