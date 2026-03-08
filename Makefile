COMPOSE_FILE ?= docker-compose.yaml
PROFILE ?=
COMPOSE_CMD ?= $(shell if docker compose version >/dev/null 2>&1; then echo "docker compose"; elif command -v docker-compose >/dev/null 2>&1; then echo "docker-compose"; fi)
DOCKER_COMPOSE = $(COMPOSE_CMD) -f $(COMPOSE_FILE)

.PHONY: help profiles urls up up-auto-down down restart logs ps build config pull guard-compose

help:
	@echo "Usage:"
	@echo "  make <target> PROFILE=<profile> [COMPOSE_FILE=<compose-file>]"
	@echo ""
	@echo "Examples:"
	@echo "  make up PROFILE=automatic-instrumentation"
	@echo "  make up PROFILE=initial COMPOSE_FILE=custom-automatic-instrumentation/docker-compose.yaml"
	@echo ""
	@echo "Targets:"
	@echo "  profiles - list all profiles found in the compose file"
	@echo "  urls     - list detected service URLs from published ports"
	@echo "  up       - start services for the selected profile"
	@echo "  up-auto-down - run in foreground and execute down on Ctrl+C/exit"
	@echo "  down     - stop services for the selected profile"
	@echo "  restart  - restart services for the selected profile"
	@echo "  logs     - follow logs for the selected profile"
	@echo "  ps       - show service status for the selected profile"
	@echo "  build    - build services for the selected profile"
	@echo "  pull     - pull images for the selected profile"
	@echo "  config   - show rendered compose config for the selected profile"

guard-profile:
	@if [ -z "$(PROFILE)" ]; then \
		echo "Error: PROFILE is required. Example: make up PROFILE=automatic-instrumentation"; \
		exit 1; \
	fi

guard-compose:
	@if [ -z "$(COMPOSE_CMD)" ]; then \
		echo "Error: neither 'docker compose' nor 'docker-compose' was found."; \
		exit 1; \
	fi

profiles: guard-compose
	@awk '\
	/^[[:space:]]*profiles:[[:space:]]*$$/ { in_profiles=1; next } \
	in_profiles && /^[[:space:]]*-[[:space:]]+/ { \
		v=$$0; gsub(/^[[:space:]]*-[[:space:]]+/, "", v); gsub(/[[:space:]]+$$/, "", v); print v; next \
	} \
	in_profiles && /^[[:space:]]*[a-zA-Z0-9_.-]+:[[:space:]]*$$/ { in_profiles=0 } \
	in_profiles && /^[^[:space:]]/ { in_profiles=0 } \
	' $(COMPOSE_FILE) | sort -u

urls: guard-profile guard-compose
	@set -eu; \
	ids="$$( $(DOCKER_COMPOSE) --profile $(PROFILE) ps -q )"; \
	if [ -z "$$ids" ]; then \
		echo "No running containers found for profile '$(PROFILE)'."; \
		exit 0; \
	fi; \
	echo "Service URLs (detected from published ports):"; \
	for id in $$ids; do \
		service="$$(docker inspect -f '{{ index .Config.Labels "com.docker.compose.service" }}' $$id 2>/dev/null || true)"; \
		[ -z "$$service" ] && continue; \
		pub=""; \
		if [ "$$service" = "jaeger" ]; then \
			pub="$$(docker port $$id 16686/tcp 2>/dev/null | head -n1 || true)"; \
		fi; \
		if [ -z "$$pub" ]; then \
			line="$$(docker port $$id 2>/dev/null | awk '/\/tcp -> / {print; exit}' || true)"; \
			[ -z "$$line" ] && continue; \
			pub="$${line##* -> }"; \
		fi; \
		host="$${pub%:*}"; \
		port="$${pub##*:}"; \
		[ "$$host" = "0.0.0.0" ] && host="localhost"; \
		[ "$$host" = "::" ] && host="localhost"; \
		echo "  - $$service: http://$$host:$$port"; \
	done

up: guard-profile guard-compose
	$(DOCKER_COMPOSE) --profile $(PROFILE) up -d
	@$(MAKE) --no-print-directory urls PROFILE=$(PROFILE) COMPOSE_FILE=$(COMPOSE_FILE)
	@echo "#############################"
	@echo "Started profile '$(PROFILE)'"
	@echo "To follow logs: 'make logs PROFILE=$(PROFILE)'"
	@echo "To stop: 'make down PROFILE=$(PROFILE)'"


down: guard-profile guard-compose
	$(DOCKER_COMPOSE) --profile $(PROFILE) down

restart: down up

logs: guard-profile guard-compose
	$(DOCKER_COMPOSE) --profile $(PROFILE) logs -f --tail=200

ps: guard-profile guard-compose
	$(DOCKER_COMPOSE) --profile $(PROFILE) ps

build: guard-profile guard-compose
	$(DOCKER_COMPOSE) --profile $(PROFILE) build

pull: guard-profile guard-compose
	$(DOCKER_COMPOSE) --profile $(PROFILE) pull

config: guard-profile guard-compose
	$(DOCKER_COMPOSE) --profile $(PROFILE) config
