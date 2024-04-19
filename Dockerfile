FROM golang:1.20-alpine AS BUILD_IMAGE
LABEL stage=builder

WORKDIR /src

# Compile the project
COPY ./go.mod .
COPY ./go.sum .
COPY ./cmd ./cmd
COPY ./pkg ./pkg

RUN go mod download

ARG CMD_NAME

# The compiled binary contains debug information. It is still impossible to debug
# on target machines due to lack of access. So we can safely remove it by compiling
# with the necessary flags or using the strip utility. The process is called stripping
# and should be quite familiar for Linux lovers
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/main ./cmd/$CMD_NAME


# Final image with the executable binary
# FROM scratch
FROM alpine:latest

WORKDIR /bin
COPY --from=BUILD_IMAGE /bin/main /bin/main

RUN mkdir -p /bot-state

ENTRYPOINT [ "/bin/main" ]
CMD [ "/bot-state" ]
