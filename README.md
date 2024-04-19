# CD-Slack notification bot

A simple bot to find the corresponding Pull Request slack thread for each Failed/Succeeded CodeBuild CD notification, and send the notification to both the CD thread and the PR thread.

## Installation

```shell
make install
```

## Prepare the environment

Check `.env.example` for the required environment variables, then fill it to `.env` (will be auto-generated during `make install`)

## Run tests

```shell
make test.all
```

## Build Docker image

```shell
make docker.build.local.arm64
make docker.build.local.amd64
```

## Run the bot with built Docker image

```shell
make docker.run.local.arm64
make docker.run.local.amd64
```
