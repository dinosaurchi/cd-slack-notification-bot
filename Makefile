# Use bash instead of shell (default)
SHELL := /bin/bash
OS_NAME := $(shell uname)
DOCKER_IMAGE_VERSION=0.0.1
DOCKER_IMAGE_NAME=codebuild-github-bot

define print_title
	echo -e "\n>>>>> $(1) <<<<<<\n"
endef

-include ./.env

#####################
# Template required #
#####################
install: install-tools
install:
	@if [ ! -f .env ]; then \
		echo "Creating .env file..."; \
		cp .env.example .env; \
	fi
lint:
	# Disable fix flag to see the errors
	@golangci-lint run --fix
check:
	@golangci-lint run --timeout 5m
	# When the command line specifies a single main package,
	# build writes the resulting executable to output.
	# Otherwise build compiles the packages but discards the results,
	# serving only as a check that the packages can be built.
	@echo -e "Try to build to make sure it works"
	@CGO_ENABLED=1 go build ./...
test: test.all


######################################################################
# Setup                                                              #
# -------------------------------------------------------------------#
# Reference: https://marcofranssen.nl/manage-go-tools-via-go-modules #
######################################################################
install-tools:
	@echo -e "Download go.mod dependencies"
	@go mod download

	@echo -e "Adding githooks to autorun CI before each commit"
	@chmod 700 ./.githooks/pre-commit
	@git config core.hooksPath .githooks

	# Use `go install` instead of `go get` to avoid updating the `go.mod`
	# when installing tools
	@echo -e "Installing tools from tools.go"
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

	@echo -e "Installing golangci-lint seperately to not have to deal with the large depedency tree"
	@if [ "$(github_ci)" = "false" ] || [ "$(github_ci)" = "" ]; then \
		if [ "$(OS_NAME)" = "Darwin" ]; then \
			echo -e "Install golangci-lint for Mac"; \
			brew install golangci/tap/golangci-lint; \
			brew upgrade golangci/tap/golangci-lint; \
		elif [ "$(OS_NAME)" = "Linux" ]; then \
			echo -e "Install golangci-lint for Linux"; \
			GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2; \
		else \
			echo -e "Not supported OS: ${OS_NAME}"; \
			exit 1; \
		fi \
	elif [ "$(github_ci)" = "true" ]; then \
		echo -e "Skipped installing golangci-lint on Github Actions CI to avoid go.sum conflict"; \
	else \
		echo -e "Not supported github_ci value: $(github_ci)"; \
		exit 1; \
	fi \

	@ golangci-lint version

######################
# Development and CI #
######################
test.all:
	@$(eval env_vars:=$(shell grep -v '^#' .env | xargs -0 | tr '\n' ' '))
	@$(eval command:="$(env_vars) CGO_ENABLED=1 CI=false gotestsum --format pkgname --no-summary=skipped -- ./... $(flags)")
	@bash -c $(command)

test.pkg.%:
	@$(eval env_vars:=$(shell grep -v '^#' .env | xargs -0 | tr '\n' ' '))
	@$(eval command:="$(env_vars) CGO_ENABLED=1 CI=false gotestsum --format pkgname --no-summary=skipped -- ./pkg/$*/... $(flags)")
	@bash -c $(command)

ci:
	@go env GOCACHE
	@echo -e "Run linting check"
	@golangci-lint run --timeout 5m

	# When the command line specifies a single main package,
	# build writes the resulting executable to output.
	# Otherwise build compiles the packages but discards the results,
	# serving only as a check that the packages can be built.
	@echo -e "Try to build to make sure it works"
	@CGO_ENABLED=1 go build ./...

	@echo -e "Run CI tests"
	@CGO_ENABLED=1 CI=true go test -timeout=5m -v ./...

build.entrypoint.%:
	@CGO_ENABLED=1 go build -o build/bin/$* ./cmd/$*

buildandrun.entrypoint.%: build.entrypoint.%
	@if [ -z "$(flags)" ]; then \
		echo -e "Missing 'flags' argument"; \
		exit 1; \
	fi
	@./build/bin/$* $(flags)

################
# Remote utils #
################

remote.upload: key_path=~/.ssh/pi_databot_rsa
remote.upload: host_user=pi
remote.upload: host_address=pi-databot.local
remote.upload:
	@if [ "$(file_path)" = "" ]; then \
		echo -e "Must provide file_path argument"; \
		exit 1; \
	fi && \
	if [ "$(dest_path)" = "" ]; then \
		echo -e "Must provide dest_path argument"; \
		exit 1; \
	fi && \
	rsync -h -P -e "ssh -i $(key_path)" -a $(file_path) $(host_user)@$(host_address):$(dest_path)

################
# Docker utils #
################
docker.login:
	@$(call print_title,Login to Docker Registry) && \
	source .env.docker && \
	(docker login $(registry_host) -u $(registry_user) --password-stdin < $(registry_password_file))

docker.build.local:
	@if [ "$(cmd_name)" = "" ]; then \
		echo -e "Must provide cmd_name argument"; \
		exit 1; \
	fi && \
	if [ "$(platform)" = "" ]; then \
		echo -e "Must provide platform argument"; \
		exit 1; \
	fi && \
	$(call print_title,Build local docker image) && \
	docker build \
		--platform linux/$(platform) \
		--pull \
		--build-arg CMD_NAME=$(cmd_name) \
		-f ./Dockerfile \
		-t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION}-$(platform) \
		-t ${DOCKER_IMAGE_NAME}:latest-$(platform) \
		. && \
	docker image prune -f --filter label=stage=builder

docker.build.local.arm64: cmd_name=main
docker.build.local.arm64: platform=arm64
docker.build.local.arm64: docker.build.local

docker.build.local.amd64: cmd_name=main
docker.build.local.amd64: platform=amd64
docker.build.local.amd64: docker.build.local

docker.run.local:
	@if [ "$(cmd_name)" = "" ]; then \
		echo -e "Must provide cmd_name argument"; \
		exit 1; \
	fi && \
	if [ "$(platform)" = "" ]; then \
		echo -e "Must provide platform argument"; \
		exit 1; \
	fi && \
	$(call print_title,Run local docker container) && \
	docker run \
		--name ${DOCKER_IMAGE_NAME} \
		--env-file .env \
		-v `pwd`/state:/bot-state \
		-d \
		${DOCKER_IMAGE_NAME}:latest-$(platform)

docker.run.local.arm64: cmd_name=main
docker.run.local.arm64: platform=arm64
docker.run.local.arm64: docker.run.local

docker.run.local.amd64: cmd_name=main
docker.run.local.amd64: platform=amd64
docker.run.local.amd64: docker.run.local
