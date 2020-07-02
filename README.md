# xdscli

Command-line tool for interacting with xDS servers

# Usage

```bash
xdscli --help
xDS protocol client to talk with management servers like Istio Pilot

Usage:
  xdscli [options] <xds> [flags]

Flags:
      --api-version string            version of xDS protocol (default "v2")
      --dial-timeout duration         dial timeout for client connections (default 2s)
      --error-detail string           the error reason that update configuration cannot be applied, using non-empty string means the discovery response will be rejected by xdscli
      --grpc-max-call-recv-size int   maximum message size that a gRPC call can accept (default 536870912)
  -h, --help                          help for xdscli
      --initial-version-info string   the version_info received with the most recent successfully processed response
      --node string                   the node making the request
      --node-metadata string          comma splitted key value pairs reresent node metadata
      --resource-names strings        list of resources to subscribe to
      --servers strings               xDS server addresses
  -v, --version                       show the version of xdscli
      --watch                         continually watch the config update
      --write-out string              set the output format (json, yaml, simple) (default "simple")
```

# Examples

```bash
xdscli eds --servers 127.0.0.1:8910 --resource-names "outbound|0||product-page.default.svc.cluster.local" --write-out json
```
