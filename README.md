# CD-Slack notification bot

## Development

Check [development tutorial](./README.dev.md)


## Prepare Docker-registry (to push the image)

If you want to build for multi-platform, you need to setup a private Docker registry to push all the cross-platform manifest to the registry

Create the file `secrets/registry_password.txt` storing your Private Registry's password which stores the final images

- This password is used for development/build purpose
- The default registry host address and username can be found at the `Makefile` target `docker.build.multiarch`

Create and modify `/etc/buildkit/buildkitd.toml` into:
```toml
debug = true

[registry]

  [registry."<registry_address>"]
    ca = ["/etc/buildkit/certs.d/<registry_address>/ca.crt"]

    [[registry."<registry_address>".keypair]]
      cert = "/etc/buildkit/certs.d/<registry_address>/client.cert"
      key = "/etc/buildkit/certs.d/<registry_address>/client.key"
```

where:

- `<registry_address>` is the current private registry host (ie. `pi-nas.local:5000`, `hub.docker.io`,...)
- `/etc/buildkit/certs.d/<registry_address>/` is the directory of the `CA` and `client` certificates for `ssl` connection to `<registry_address>`


Create symbolic-link for `cert.d` 

- For `macos`

    ```shell
    $ sudo ln -s ~/.docker/certs.d /etc/buildkit/certs.d
    ```
- For `linux`

    ```shell
    $ sudo ln -s /etc/docker/certs.d /etc/buildkit/certs.d
    ```

You have to create the `~/.docker/certs.d/<registry_address>` with the following structure (if not existing):

- For `macos`

    ```
    ~/.docker/certs.d
        |_ <registry_address>
            |_ ca.crt
            |_ client.cert
            |_ client.
    ```

- For `linux`

    ```
    /etc/docker/certs.d
        |_ <registry_address>
            |_ ca.crt
            |_ client.cert
            |_ client.
    ```


## Build the image

### Local build

```shell
$ make docker.build.local.arm64 cmd_name=<target_cmd_name>
$ make docker.build.local.amd64 cmd_name=<target_cmd_name>
```

### Local run

```shell
$ make docker.run.local.arm64 cmd_name=<target_cmd_name>
$ make docker.run.local.amd64 cmd_name=<target_cmd_name>
```

### Multi-architecture build (and push to a private registry)

Create a cross-platform builder (if not created yet):
```shell
$ docker.build.multiarch.builder.create
```

Run multi-arch build:
```shell
$ make docker.build.multiarch cmd_name=<target_cmd_name>
```

(Optional) Remove the builder:
```shell
$ make docker.build.multiarch.builder.remove
```

## Pull and run the image

If you built the image with `make docker.build.local.*`:

- You do not need to pull and re-tag the image, just run it directly

If you built the image with `make docker.build.multiarch`:

- You have to pull and re-tag the image

    ```shell
    $ docker login \
        <registry_address> \
        -u my_username \
        --password-stdin
    $ docker pull \
        --platform arm64 \
        <registry_address>/<your_image_name>:latest 
    $ docker tag \
        <registry_address>/<your_image_name>:latest \
        <your_image_name>:latest
    ```

### Build and deploy to remote host

Build and deploy to remote host (auto remove existing container)
```shell
$ make docker.buildanddeploy.main
```

To only deploy with existing image, run
```shell
$ make docker.deploy cmd_name=<target_cmd_name>
```
