PROJECTNAME=$(shell basename "$(PWD)")


help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'


## build: Build docker image
build:
	@docker build -t api_go .


## run: Run server api_go + postgres + redis
run:
	@@docker compose up


## clean: Delete docker image
clean:
	@docker stop api_go
	@docker rmi api_go
