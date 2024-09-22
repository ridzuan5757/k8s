# Service

In Kubernetes, a Service is a method for exposing a network application that is
running as one or more Pods in the cluster.

A key aim of SErvices in Kubernetes is that we do not need to modify existing
application to use an unfamiliar service discovery mechanism. We can run code in
Pods, whether this is a code designed for a cloud-native world, or an older app
we have containerized. We use Service to make that set of Pods available on the
network so that clients can interact with it.

If we use `Deployment` to run the app, that Deployment can create and destroy
Pods dynamically.. From one moment to the next, we do not know how many of those
Pods are working and healthy. We might not even know that whose healthy Pods are
named.

Kubernetes Pods are created and destroyed to match the desired state of the
cluster. Pods are ephemeral resource (we should not expect that an individual
Pod is reliable and durable).

Each Pod gets its own IP address (Kubernetes expects network plugins to ensure
this). For a given Deployment in a cluster, the set of Pods running in one
moment in time could be different from the set of Pods running that application
a moment later.

This leads to a problem: if some set of Pods (lets call them "backends")
provides functionality to other Pods (call them "frontends") inside the cluster,
how do the frontends find out and keep track of which IP address to connect to,
so that the frontend can use the backend part of the workload?

# Services in Kubernetes

The Service API, part of Kubernetes, is an abstraction to help us expose groups
of Pods over a network. Each Service object defines a logical set of endpoints
(usually these endpoints are Pods) along with a policy about how to make those
pods accessible.

For example, consider a stateless image-processing backend which is running with
3 replicas. Those replicas are fungible - frontends do not care which backend
they use. While the actual Pods that compose the backend set may change, the
frontend clients should not need to be aware of that, nor should they need to
keep track of the set of backends themselves.

The Service abstraction enables this decoupling.

The set of Pods targeted by a Service is usually determined by a selector that
we define.

##### Ingress

If the workload speaks HTTP, we might choose to use an Ingress to control how
web traffic reaches that workload. Ingress is not a Service type, but it acts as
the entry point for the cluster. An ingress lets us consolidate the routing
rules into a single resource, so that we can expose multiple components of the
workload, running separately in the cluster, behind a single listener.

##### Gateway

The Gateway API for Kubernetes provides extra capabilities beyond Ingress and
Service. We can add Gateway to the cluster - it is a family of extension API's
implemented using CRDs - and then use these to configure access to network
services that are running in our cluster.

# Cloud-native service discovery

If we are able to use Kubernetes APIs for service discovery in the application,
we can query the API server for matching EndpointSlices. Kubernetes updates the
EndpointSlices for a Service whenever the set of Pods in a Service changes.

For non-native applications, Kubernetes offers ways to place a network port or
load balancer in between the application and the backend Pods.

Either way, the workload can use these service discovery mechanisms to find the
target it wants to connect to.

# Defining a Service

A service is an object such as Pod or ConfigMap. We can create, view or modify
Service definitions using the Kubernetes API. Usually we use a tool such as
`kubectl` to make those API calls for us.

For example, suppose we have a set of Pods that each lsiten on TCP port 9376 and
are labelled as `app.kubernetes,io/name=MyApp`. We can define a SErvice to
publish that TCP listener:

```yaml
apiVersion: v1
kind: Service
metadata:
    name: my-service
spec:
    selector:
        app.kubernetes.io/name: MyApp
    ports:
        - protocol: TCP
          port: 80
          targetPort: 9376
```

Applying this manifest creates a new service named `my-service` with the default
ClusterIP service type. The service targets TCP port 9376 on any Pod with the
`app.kubernetes.io/name:Myapp` label.

Kubernetes assigns this Service an IP address (the cluster IP), that is used by
the virtual IP address mechanism. 

The controller for that Service continuously scans for Pods that match its
selector, and then makes any necessary updates to the set of EndpointSlices for
the Service. 

The name of a Service object must be a valid RFC 1035 label name.

> [!NOTE]
> A service can map any incoming `port` to a `targetPort`. By default and for
> convenience, the `targetPort` is set to the same value as the `port` field.

# Port definitions

Port definitions in Pods have names, and we can reference these names in the
`targetPort` attribute of a Service. For example, we can bind the `targetPort`
of the Service to the Pod port in the following way:

```yaml
apiVersion: v1
kind: Pod
metadata:
    name: nginx
    labels:
        app.kubernetes.io/name: proxy
spec:
    containers:
        - name: nginx
          image: nginx:stable
          ports:
            - containerPort: 80
              name: http-web-svc

---

apiVersion: v1
kind: Service 
metadata:
    name: nginx-service
spec:
    selector:
        app.kubernetes.io/name: proxy
    ports:
        - name: name-of-service-port
          protocol: TCP
          port: 80
          targetPort: http-web-svc
```

This works even if there is a mixture of Pods in the Service using a single
configured name, with the same network protocol available via different port
numbers. This offers a lot of flexibility for deploying and evolving the
Services. For example, we can change the port numbers that Pods expose in the
next version of the backend software, without breaking clients.

The default protocol for Services is TCP, other supported protocol can also be
used:
- SCTP
- TCP
- UDP

Because many services need to expose more than one port, Kubernetes supports
multiple port definitions for a single Service. Each port definition can have
the same `protocol`, or a different one.

# Services without selectors

Services most commonly abstract access to Kubernetes Pods thanks to the
selector, but when used with a corresponding set of EndPointSlices objects and
without a selector, the Service can abstract other kinds of backends, including
ones that run outside the cluster.

For example:
- We want to have an external database cluster in production, but in the test
  environment we use our own datasets.
- We want to point the Service to a Service in a different namespace or on
  another cluster.
-  We are migrating a workload to Kubernetes. While evaluating the approach, we
   run only a portion of the backends in Kubernetes.

In any of these scenarios we can define a Service without specifying a selector
to match Pods. For example:

```yaml
apiVersion: v1
kind: Service
metadata:
    name: my-service
spec:
    ports:
        - name: http
          protocol: TCP
          port: 80
          targetPort: 9376
```

Because this Service has no selector, the corresponding EndpointSlice (and
legacy Endpoints) objects are not created automatically. We can then map the
Service to the network address and port where it is running, by adding an
EndPointSlice object manually. For example:

```yaml
apiVersion: discovery.k8s.io/v1
kind: EndpointSlice
metadata:
    # by convention, use the name of of the Service as prefix for the name of
    # the endpointslice
    name: my-service-1
    labels:
        # we should set the `kubernetes.io/service-name` label.
        # set its value to match the name of the Service
        kubernetes.io/service-name: my-service
addressType: IPv4
ports:
    # should amtch with the name of the service port defined above
    - name: http
      appProtocol: http
      protocol: TCP
      port: 9376
endpoitns:
    - addresses:
        - 10.4.5.6
    - addresses:
        - 10.1.2.3
```

# Custom EndpointSlices

When we create an EndpointSlice object for a Service, we can use any name for
the EndpointSlice. Each EndpointSlice in a naespace must have unique name. We
link an EndpointSlice to a Service by setting the `kubernetes.io/service-name`
label on that EndpointSlice.

> [!NOTE]
> The endpoint IPs must not be loopback IP (127.0.0.0/8 for IPv4, ::1/128 for
> IPv6), or link-local (169.254.0.0/16 and 224.0.0.0/24 for IPv4, fe80::/64 for
> IPv6).
> 
> The endpoint IP addresses cannot be the cluster IPs of other Kubernetes
> Services, because `kube-proxy` does not support virtual IPs as a destination.

For an EndpointSlice that we create ourselves, or in our own code, we should
also pick a value to use for the label `endpointslice.kubernetes.io/managed-by`.
If we create our own controller to manage EndpointSlices, consider using a value
similar to `my-domain.example/name-of-controller`. If we are using a third party
tool, use the name of the tool in all-lowercase and change spaces and other
punctuation to dashes (`-`). If people are directly using a tool such as
`kubectl` to manage EndpointSlices, use a name that describes this manual
management, such as `staff` or `cluster-admins`. We should avoid using the
reserved value `controller`, which identifies EndpointSlices managed by
Kubernetes' own control plane.

# Accessing a Service Without a Selector

Accessing a Service without a selector works the same as if it had a selector.
Traffic is routed to one of the two endpoints defined in the EndpointSlice
manifest.

> [!NOTE]
> The Kubernetes API server does not allow proxying to endpoints that are not
> mapped to pods. Actions such as `kubectl port-forward service/service-name
> forwarded-port:service-port` where the service has no selector will fail due
> to this constraint. This prevents the Kubernetes API server from being used as
> a proxy to endpoints the caller may not be authorized to access.

# EndpointSlices

EndpointSlices are objects that represent a subset (a slice) of the backing
entwork endpoints for a Service.

The Kubernetes cluster tracks how many endpoints each EndpointSlice represents.
If there are so many endpoint for a Service that a threshold is reached, then
Kubernetes addres another empty EndpointSlice and stores new endpoint
information there. By default, Kubernetes makes a new EndpointSlice once the
existing EndpointSlices all cotnain at least 100 endpoints. Kuberentes does not
make the new EndpointSlice until an extra endpoint needs to be added.

# Endpoints

In the Kubernetes API, an Endpoints defines a list of network endpoints,
typically referenced by a Service to define which Pods the traffic can be sent
to. The EndpointSlice API is the recommended replacement for Endpoints.

# Over-capacity endpoints

Kubernetes limits the number of endpoints that can fit in a single Endpoints
object. When there are over 1000 backing endpoints for a Service, Kubernetes
truncates the data in the Endpoints object. Because a Service can be linked with
more than one EndpointSlice, the 1000 backing endpoint limit only affects the
legacy Endpoints API.

In that case, Kuberentes selects at most 1000 possible backend endpoitns to
store into the Endpoints object, and sets an annotation on the Endpoints:
`endpoints.kubernetes.io/overcapacity:truncated`. The control plane also removes
that annotation if the number of backend Pods drops below 1000.

Traffic is still sent to backends, but any load balancing mechanism that relies
on the legacy Endpoints API only sends traffic to at most 1000 of the available
backing endpoints.

The same API limit means that we cannot manually update an Endpoints to have
more than 1000 endpoints.

# Application Protocol

The `appProtocol` field provides a way to specify an application protocol for
each Service port. This is used as a hint for implementations to offer richer
behaviour for protocols that they understand. The value of this field is
mirrored by the corresponding Endpoints and EndpointSlice objects.

This field follows standard Kubernetes label syntax. Valid values are one of:
- IANA standard service names.
- Implementation-defined prefixed names such as
  `mycompany.com/my-custom-protocol`.
- Kubernetes-defined prefixed names:

|**Protocol**|**Description**|
|---|---|
|`kubernetes.io/h2c`|HTTP/2 over cleartext as described in RFC 7540|
|`kubernetes.io/ws`|WebSocket over cleartext as described in RFC 6455|
|`kubernetes.io/wss`|WebSocket over TLS as described in RFC 6455|

# Multi-port Services

For some services, we need to expose more than one port. Kubernetes lets us
configure multiple port definitions on a Service object. When usingmultiple
ports for a Service, we must give all of the ports name so that these are
unambiguous. For example:

```yaml
apiVersion: v1
kind: Service
metadata:
    name: my-service
spec:
    selector:
        app.kubernetes.io/name: MyApp
    ports:
        - name: http
          protocol: TCP
          port: 80
          targetPort: 9376
        - name: https
          protocol: TCP
          port: 443
          targetPort: 9377
```

> [!NOTE]
> As with Kubernetes names in general, names for ports must only contain
> lowercase alphanumeric characters and `-`. Port names must also start and end
> with an alphanumeric character.
>
> For example, the names `123-abc` and `web` are valid, but `123_abc` and `-web`
> are not.

# Service Type

For some parts of the application (for example, frontends) we may want to expose
a Service onto an external IP address, one that is accessible from outisde of
the cluster. 

Kubernetes Service types allow us to specify what kind of Service we want. The
available `type` values and their behaviors are:

##### `ClusterIP`

Exposes the Service on a cluster-internal IP. Choosing this value makes the
Service only reachable from within the cluster. This is the default that is used
if we do not explicitly specify a `type` for a service. We can expose the
service to the public using an Ingress or Gateway.

##### `NodePort`

Exposes the Services on each Node's IP at a static port. To make the node port
available, Kubernetes sets up a cluster IP address, the same as if we had
requested the Service of `type: ClusterIP`.

##### `LoadBalancer`

Exposes the Service externally using an external load balancer. Kubernetes does
not directly offer a load balancing component; we must provide one, or we can
integrate the Kubernetes cluster with a cloud provider.

##### `ExternalName`

Maps the Service to the contents of the `externalName` field (for example, to
the hostname `api.foo.bar.example`). The mapping configures the cluster's DNS
server to return a CNAME record with that external hostname value. No proxying
of any kind is set up.

The `type` field in the Service API is designed as nested functionality. Each
level adds to the previous. However, there is an exception to this nested
design. We can define a `LoadBalancer` service by disabling the load balancer
`NodePort` allocation.

# `type: ClusterIP`

This default Service type assigns an IP address from a pool of IP addresses that
the default cluster has reserved for that purpose. Serveral of other types for
Services build on the `ClusterIP` type as a foundation. If we define a Service
that has the `.spec.clusterIP` set to `None`, then Kubernetes does not assign an
IP address. This turn the service into a headless service.

## Choosing own IP address

We can specify our own IP address as part of `Service` creation request. To do
this, set the `spec.clusterIP` field. For example, if we already have an
existing DNS entry that we wish to resue, or legacy systems that are configured
for a specific IP address and difficult to configure.

The IP address that being chisen must be a valid IPv4 or IPv6 range  address
from within the `service-cluster-ip-range`. CIDR range that is configured for
the API server, If we try to create a Service with an invalid `clusterIP`
address value, the API server will return a 422 HTTP status code to indicate
that there is a problem.

# `type: NodePort`

If we set the `type` field to `NodePort`, the Kubernetes control plane allocates
a port from range specified by `--service-ndoe-port-range` flag from 30000 to
32767. Each node proxies that port (same port number on every Node) into the Service.
 The Service reports the allocated port in its `.spec.ports[*].nodePort` field.

 Using a NodePort gives us the freedom to set up our own load balancing
 solution, to configure environments that are not fully supported by Kubernetes,
 or even to expose one or more nodes' IP addresses directly.

 For a node port Service, Kubernetes additionally allocates a port (TCP, UDP or
 SCTP) to match the protocol of the Service). Every node in the cluster
 configures itself to listen on that assigned port and to forard traffic to one
 of the ready endpoints associated with that Service. We will be able to contact
 the `type: NodePort` service, from outside the cluster, by connecting to any
 node using the appropriate protocol, and the appropriate port.

 ## Choosing own port

 If we want a specific port number, we can specify a value in the `nodePOort`
 field. The control plane will either allocate us that port or report that the
 API transaction failed. This means that we need to take care of possible port
 collision ourselves. We also have to use a valid port number, one that is
 inside the range configured for NodePort use. Here is an example manifest for a
 service that specifies NodePort value:

 ```yaml
apiVersion: v1
kind: Service
metadata:
    name: my-service
spec:
    type: NodePort
    selector:
        app.kubernetes.io/name: MyApp
    ports:
        - port: 80
          targetPort: 80
          nodePort: 30007
 ```

## Reserve NodePort ranges to avoid collisions

The policy for assignning ports to NodePort services applies to both the auto
assignment and the manual assignment scenarios. When a user wants to create a
NodePort service that uses a specific port, the target port may conflict with
another port that has already been assigned.

To avoid this problem, the port range for NodePort services is divided into two
bands. Dynamic port assignment uses the upper band by default, and it may use
the lower band once the upper band has been exhausted. Users can then allocate
from the lower band with lower risk of port collision.

## Custom IP address configuration for `type: NodePort` Services

We can set up nodes in the cluster to use a particular IP address for serving
node port services. We might want to do this if each node is connected to
multiple networks (for example: one network for application traffic, and another
network for traffic between nodes and the control plane).

If we want to specify particular IP addresses to proxy the port, we can set the
`--nodeport-addresses` flag for kube-proxy or the equivalent `nodePortAddress`
field of the kube-proxy configuraiton file to particular IP block(s).

This flag takes a comma-delimited list of IP blocks (e.g `10.0.0.0/8`,
`192.0.2.0/25`) to specify IP address ranges that kube-proxy should consider as
local to this node.
