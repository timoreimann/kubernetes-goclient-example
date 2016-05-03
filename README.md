# kubernetes-goclient-example

Example on how to interact with Kubernetes via the [Golang client library](https://github.com/kubernetes/kubernetes/blob/release-1.2/docs/devel/client-libraries.md).

## Requirements

- A Golang version with vendoring enabled.
- `kubectl proxy` running and proxying API traffic via (the default) `localhost:8001` socket.


## Kubernetes Version

This sample projects uses a vendored copy of the Kubernetes 1.2 package (and its dependencies).

See Kubernetes' [Development Guide](https://github.com/kubernetes/kubernetes/blob/master/docs/devel/development.md) on how to properly set up Kubernetes in your `GOPATH`.

## Usage

```
go build -o client
./client <operation>
```

where `<operation>` is one of the following:

- `version`: Queries the Kubernetes server version.
- `deploy`: Deploys NGINX via the [Deployments API](http://kubernetes.io/docs/user-guide/deployments/).

There is also a convenience script `run.sh` that builds on demand (or forcefully when `-b` is given as first parameter) and executes the client with all arguments passed along. That is, invoking

`./run.sh version`

will build the client _iff_ it does not exist yet, and executes it subsequently with `version` passed as parameter.
