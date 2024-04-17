# CD-Slack notification bot (Development tutorial)

## Directory structure

```
<current_working_dir>
    |_ cmd (app entrypoints)
    |   |_ main
    |   |   |_ main.go (main's entrypoint)
    |   :
    |_ pkg (local libraries/packages)
    |   |_ mypkg1
    |   |   |_ api.go
    |   |   |_ api_test.go
    |   |   |_ utils.go
    |   |   |_ utils_test.go
    |   |_ mypkg2
    |   |   |_ mytools.go
    |   |   |_ mytools_test.go
    |   :
    |_ build (built binanries)
    |   |_ bin
    |       |_ app1 (binary)
    |       |_ app2 (binary)
    |       :
    |       :
    |_ tools.go
```
- Put your app's entrypoint in `./cmd`
  - Each entrypoint is a directory, with a `main.go` file with function `main` implemented
  - Each directory is named after the app name

- Test file has postfix `_test.go`

- Declare your development tools' dependencies in `./tools.go`, which are not included when running `go build` for production
    ```go
    //go:build tools
    // +build tools

    package main

    import (
        _ "golang.org/x/lint/golint"
        _ "gotest.tools/gotestsum"
    )

    ```

## Add new dependencies

Example:
```shell
$ go get golang.org/x/lint/golint
$ go get gotest.tools/gotestsum
```

- The dependency will be auto-added to `go.mod` (and `go.sum` will be updated accordingly)

## Build and run application

Each application has a corresponding entrypoint in `./cmd`

Build:
```shell
$ make build.entrypoint.main
```

Build and run:
```shell
$ make buildandrun.entrypoint.main
```
