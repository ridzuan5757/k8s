# Container Environment

k8s container environment provides several important resources to containers:
- A filesystem, which is a combination of an image and one or more volumes.
- Information about the container itself.
- Information about other object in the cluster.

### Container Information

The `hostname` of a container is the name of the pod in which the container is
running. It is available through the `hostname` command or the `gethostname`
function call in libc.

The pod name and namespace are available as environment variables throguht the
downward API.

User defined environment variables from the pod definition are also available to
the container,as are any environment variables specified statically in the
container image.

### Cluster Information

A list of all services that were running when a container was created is
available to that container as environment variables. This list is limited to
services within the same namespace as the new container's pod and k8s control
plane service.

For a service named `foo` that maps to a container named `bar`, the following
variables are defined.

```bash
FOO_SERVICE_HOST=<the host the service is running on>
FOO_SERVICE_PORT=<the port the service is running on>
```

Services have dedicated IP addresses and are available to the container via DNS,
if DNS addon is enabled.
