---
icon: lucide/rocket
---

# Getting started

Follow this guide to get fltr running with one [**Block Rule**](../reference/block-rules.md) and one [**Allow Upstream**](../reference/environment-variables.md#upstreams).

Every install path needs the same two things: a JSON file containing the Block Rules and an HTTP service to receive allowed requests.

## Install

fltr ships as a pre-built binary, a Docker image, and can be built from source.

### Pre-built binary

Download the binary for your platform from the [releases page](https://github.com/kfmndev/fltr/releases).

!!! tip
    Rename the downloaded binary to `fltr`, which is also the name used in the following steps.

!!! warning
    Before continuing, make sure to add the execute permission, which allows the binary to run at all.
    It can be added by using this command:

    ```properties
    chmod +x fltr
    ```

### Docker

Images are published to `ghcr.io/kfmndev/fltr` for `linux/amd64` and `linux/arm64`.

See [Docker](docker.md) for compose setup and image tags.

### From source

Build with [Go 1.25](https://go.dev) or newer:

```properties
go build -o fltr .
```

## Minimal setup

Create a `block_rules.json` with a single Block Rule:

```json title="block_rules.json"
{
    "credentials": ["password", "secret"]
}
```

Set the Allow Upstream and run fltr:

=== "Binary"

    !!! warning
        The following commands only work, if fltr has been added to the path. Adjust the path, depending on how you installed fltr. If the binary lives in the current directory change `fltr` to `./fltr`.

    ```sh
    FLTR_ALLOW_UPSTREAM=http://allow.example.com \
    fltr
    ```

=== "Docker"

    Due to Docker's isolation, the setup differs: the Block Rules file must be mounted, the environment variable set explicitly, and the internal port mapped to a host port.

    ```sh
    docker run -d \
      --name fltr \
      -p 8080:8080 \
      -v "$PWD/block_rules.json":/config/block_rules.json:ro \
      -e FLTR_ALLOW_UPSTREAM=http://allow.example.com \
      ghcr.io/kfmndev/fltr
    ```

    For more details or an example Docker Compose configuration, see [Docker](docker.md).

By default fltr listens on `:8080` and reads `block_rules.json` from the current directory. Both are configurable via [Environment variables](../reference/environment-variables.md).

## Try a request

With the example configuration running, send a request that matches the `credentials` Block Rule:

```sh
curl -X POST http://localhost:8080 \
    -H 'Content-Type: text/plain' \
    -H 'Title: account details' \
    -d 'password secret'
```

A request whose [**Searchable Content**](../reference/block-rules.md#matching-semantics) contains every term of a Block Rule is blocked. With no [**Block Upstream**](../reference/environment-variables.md#upstreams) configured, fltr answers `200 Request blocked, discarded`.

Send a request that matches no Block Rule to route it to the [**Allow Upstream**](../reference/environment-variables.md#upstreams):

```sh
curl -X POST http://localhost:8080 \
    -H 'Content-Type: text/plain' \
    -d 'hello service'
```

See [Block Rules](../reference/block-rules.md#matching-semantics) to learn more about the matching or have a look at the request flow in [How it works](../reference/how-it-works.md#request-flow).
[HTTP behavior](../reference/http-behavior.md) documents every response fltr produces itself.
