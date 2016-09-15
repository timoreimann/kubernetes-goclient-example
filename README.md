# kubernetes-goclient-example

Example on how to interact with Kubernetes via the [Golang client library](https://github.com/kubernetes/kubernetes/blob/release-1.2/docs/devel/client-libraries.md).

## Requirements

A Golang version with support for vendoring (1.5+) is required.

## Kubernetes Version

This sample project includes a vendored copy of the [client-go 1.4 package](https://github.com/kubernetes/client-go) (and its dependencies). It uses [Glide](https://glide.sh/) for dependency management.

## Usage

```
go build -o client
./client <operation>
```

where `<operation>` is one of the following:

- `version`: Queries the Kubernetes server version.
- `deploy`: Deploys NGINX via the [Deployments API](http://kubernetes.io/docs/user-guide/deployments/) and exposes a [Service](http://kubernetes.io/docs/user-guide/services/).

The client tries to reach the API server through `http://127.0.0.1:8001` which happens to match the default tunnel set up by `kubectl proxy`. Custom overriding settings for the API server URL, a bearer token, and the path to a CA file can be injected through the environment variables `SERVER`, `TOKEN`, and `CA_FILE`, respectively.

There is also a convenience script `run.sh` that builds on demand (or forcefully when `-b` is given as first parameter) and executes the client with all arguments passed along. That is, invoking

`./run.sh version`

will build the client _iff_ it does not exist yet, and executes it subsequently with `version` passed as parameter.
