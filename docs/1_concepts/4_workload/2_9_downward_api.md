# Downward API

Downward API allows containers to consume information about themselves or the
cluster without using k8s client or API server.

An example is an existing application that assumes a particular well-known
environment variable holds a unique identifier. One possibility is to wrap the
application, but that is tedious and error-prone, and it violates the goal of
low coupling.

A better option would be to use the pod's name as identifier, and inject the
pod's name into the well-known environment variable.

In k8s, there are 2 ways to expose pod and container fields to a running
container:
- as environment variables
- as files in a `downwardAPI` volume

Together, these 2 ways of exposing pod and container fields are called the
downward API.

## Available fields

Only some k8s API fields are available through the downward API. This section
list which fields can be made available.

Information from available pod-level fields can be passed using `fieldRef`. At
the API level, the `spec` for pod always defines at least one container.
Information from available container-level fields can be passed using
`resourceFieldRef`.

### Information available via `fieldRef`

For some pod-level fields, we can provide them to a container either as an
environment variable or using `downwardAPI` volume. The fields available via
either mechanism are:

|Field|Description|
|---|---|
|`metadata.name`|the pod's name|
|`metadata.namespace`|the pod's namespace|
|`metadata.uid`|the pod's unique ID|
|`metadata.annotations['<KEY>']`|value of the pod's annotation named `<KEY>`|
|`metadata.labels['<KEY>']`|the text value of the pod's level named `<KEY>`|

The following information is available through environemnt variables **but not
as a `downwardAPI` volume `fieldRef`*:

|Field|Description|
|---|---|
|`spec.serviceAccountName`|the name of the pod's service account|
|`spec.nodeName`|the name of the node where the pod is executing|
|`status.hostIP`|the primary IP address of the node to which pod is assigned|
|`status.hostIPs`|
the IP addresses is a dual-stack version of `status.hostIP`,the first is always
the same as `status.hostIP`. The field is avaialble if `PodHostIps` feature gate
is enabled.|

