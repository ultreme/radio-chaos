include .env
export sed 's/=.*//' .env

install:
	go install

run: install
	chaos-bot -t $(DISCORD_TOKEN)
