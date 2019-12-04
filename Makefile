# Project vars
PROJECT = goaccounts
SERVICE = web
DB = pg
NG = nginx
DC_PROD = docker-compose -p $(PROJECT) -f docker-compose.yml
DC = $(DC_PROD) -f docker-compose.dev.yml
SCALE = web=3

# Main targets
up:
	$(DC) up -d --build

help:
	@echo "  up-prod    docker-compose up project $(PROJECT) with production build"
	@echo "  down-prod  docker-compose down production build"
	@echo "  up         docker-compose build and up project $(PROJECT) with dev build"
	@echo "  restart    docker-compose restart"
	@echo "  ps         docker-compose ps"
	@echo "  bash       app $(SERVICE) bash"
	@echo "  log        app $(SERVICE) log"
	@echo "  attach     app $(SERVICE) attach (for python breakpoint())"
	@echo "  bash-pg    DB $(DB) bash"
	@echo "  psql       DB $(DB) psql"
	@echo "  refactor   format code"
	@echo "  clean      clean dev staff"
	@echo "  ...        all commands in Makefile"

# Prod
up-prod:
	$(DC_PROD) up -d --build --scale $(SCALE)
down-prod:
	$(DC_PROD) down --rmi local --remove-orphans
ps-prod:
	$(DC_PROD) ps

# Dev docker targets
build:
	$(DC) build
stop:
	$(DC) stop
down:
	$(DC) down --rmi local --remove-orphans
build-no-cache:
	$(DC) build --no-cache
restart: stop down build up
	@echo "docker-compose has ben restarted!"
ps:
	$(DC) ps

refactor:
	$(DC) exec $(SERVICE) bash -c 'go fmt gopayment && go fmt gopayment/account'
	@echo "Refactor done!"

# $(SERVICE) targets
bash:
	$(DC) exec -e COLUMNS="`tput cols`" -e LINES="`tput lines`" $(SERVICE) bash
log:
	$(DC) logs -f $(SERVICE)
attach:
	docker attach $(PROJECT)_$(SERVICE)_1

# $(DB) targets
bash-pg:
	$(DC) exec -e COLUMNS="`tput cols`" -e LINES="`tput lines`" $(DB) /bin/sh
psql:
	$(DC) exec $(DB) psql accounts_db -h localhost -U accounts_user
log-pg:
	$(DC) logs -f $(DB)
# Dev mode, log all DB queries
psql-log:
	$(DC) exec $(DB) /bin/sh -c 'tail -f /var/lib/postgresql/data/log/postgresql*.log'
bash-nginx:
	$(DC) exec -e COLUMNS="`tput cols`" -e LINES="`tput lines`" $(NG) /bin/sh

# Dockerfile validations
hadolint:
	docker run --rm -i hadolint/hadolint < Dockerfile


.PHONY: all up build stop down build-no-cache restart ps bash log attach bash hadolint
