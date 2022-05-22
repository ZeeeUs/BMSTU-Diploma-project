PROJECTNAME=$(shell basename "$(PWD)")

help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## build: Build docker image
build:
	@docker build -t dashboard_back .


## run: Run server in docker on port 8000
run:
	@docker run -d --rm -p 8000:8000 --name dashboard_back dashboard_back


## clean: Delete docker image
clean:
	@docker stop dashboard_back
	@docker rmi dashboard_back


## rerun: Rerun
rerun:
	@docker stop dashboard_back
	@make run


## rebuild: Rebuild and restart
rebuild:
	@make clean
	@make build
	@make run