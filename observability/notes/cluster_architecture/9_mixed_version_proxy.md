# Mixed Version Proxy

k8s 1.29 icludes an alpha feature that lets an API server proxy a resource
requests to other peer API servers. This is useful when there are multiple API
server running different versions of k8s in one cluster. For example, during
long-lived rollout to a new release of k8s.

This enables cluster administrator to configure highly available clusters that
can be upgraded more safely, by directing resource requests made during the
upgrade to the correct kube-apiserver. That proxying prevents users from seeing
unexpected 404 not found errors that stem from upgrade process.

This mechanism is called the Mixed Version Proxy.

## Enabling Mixed Version Proxy

Ensure that `UnknownVersionInteroperabilityProxy` feature gate is enabled when
we start the API server.

```bash
kube-apiserver \
--feature-gates=UnknownVersionInteroperabilityProxy=true \
# required command line arguments for this feature
--peer-ca-file=<path to kube-apiserver CA cert>
--proxy-client-cert-file=<path to aggregator proxy cert>,
--proxy-client-key-file=<path to aggregator proxy key>,
--requestheader-client-ca-file=<path to aggregator CA cert>,
# requestheader-allowed-names can be set to blank to allow any Common Name
--requestheader-allowed-names=<valid Common Names to verify proxy client cert against>,

# optional flags for this feature
--peer-advertise-ip=`IP of this kube-apiserver that should be used by peers to proxy requests`
--peer-advertise-port=`port of this kube-apiserver that should be used by peers to proxy requests`

# …and other flags as usual
```

### Proxy transport and authentication between API servers
- The source kube-apiserver resues existing APIserver client authentication
  flags `--proxy-client-cert-file` and `--proxy-client-key-file` to present its
  identity that will be verified by its peer (destination kube-apiserver). The
  destination API server verifies that peer connection based on the
  configuration we specify using the `--requestheader-client-ca-file` command
  line argument.
- To authenticate the destination server's serving cert, we must configure a
  certificate authority bundle by specifying the `--peer-ca-file` command ine
  argument to the **source** aPI server.

### Configuration for peer API server connectivity

To set the network location of a kube-apiserver that peers will use to proxy
requests, use the `--peer-advertise-ip` and `--peer-advertise-port` command line
argments to kube-apiserver or specify these fields in the API server
configuration file.

If these flags are unspecified, peers will use the value from either
`--advertise-address` or `--bind-address` command line argument to the kube-api
server. If those too, are unset, the host's default interface is used.

## Mixed versin proxying

When we enable mixed version proxying, the aggregation layer loads a special
filter that does:
- When resource request reaches an API server that cannot server the API (due to
  backward imcompatibility or API is turned off) the API server attempts to send
  the request to the peer API server that can serve the requested API. It does
  so by identifying API groups / versions / resources that the local server does
  not recognize, and tries to proxy those requests to a peer API server that is
  capable handling the request.
- If the peer API fails to respond, the source API server responds with 503
  service unavailable error.

### Mechanism

When an API server receives a resource request, it first checks which API
servers can server the requested resource. This check happens using the internal
`StorageVersion` API.
- If the resource is know to the API server that received the request (for
  example, `GET /api/v1/pods/some-pod`), the request is handled locally.
- If there is no internal `StorageVersion` object found for the requested
  resource (for example, `GET /my-api/v1/my-resource`) and the configured API
  service specifies proxying to an extension API server, that proxying happens
  follwing the usual flow for extension APIs.
- If a valid internal `StorageVersion` object is found for the requested
  resource (for example, `GET /batch/v1/job`) and the API server trying to
  handle the request has `batch` API disabled, then the handling API server
  fetches the peer API servers that do serve the relevant API group / version /
  resource (`api/v1/batch` for this case) using the information in the fetched
  `StorageVersion` object. The handling API server then proxies the request to
  one of the matching peer kube-apiservers that are aware of the requested
  resource.
    - If there is no peer known for that API group / version / resource, the
      handling API server passes the request to its own handler chain which
      should eventuallay turn a 404 not found response.
    - If the handling API server has identified and selected a peer API server,
      but that peer fails to respond (for reasons such as network connectivity
      issues, or data race between the request being received and a controller
      registering the peer'ss info into the control plane), then the handling
      API server responds with a 503 service unavailable error.
