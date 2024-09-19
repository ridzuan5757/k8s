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


