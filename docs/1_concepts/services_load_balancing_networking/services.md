# Service

Service is a method for exposing a network application that is running as one or
more Pods in the cluster.

The aim of Services in k8s is that we do not need to modify the existing
application to use an unfamiliar service discovery mechanism. We can run code in
Pods, whether this is a code designed for a cloud native world, or an older app
we have containerized. We use a Service to make that set of Pods available on
the network so that clients can interact with it.

If we use a Deployment to run the app, that Deployment can create and destroy
Pods dynamically. From one moment to the next, we do not know how many of those
Pods are working and healthy, we might not even know what those healthy Pods are
named K8s Pods are created and destroyed to match the desired state of the
cluster. Pods are ephemeral resources as we should not expect that an individual
Pod is reliable and durable.

This leads to a problem:
If some set of Pods called "backends" provides functionality to other Pods
called "frontends" inside the cluster, how do the frontends find out and keep
track of which IP address to connect to, so that the frontend can use the
backend part of the worklod?

## Services in k8s

The Service API, part of k8s is an abstraction to help us expose groups of Pods
over a network. Each Service object defines a logical set of endpoints, usually
these endpoints are Pods along with a policy about how to make those pods
accessible. 

For example, consider a stateless image processing backend which is running with
3 replicas. Those replicas are fungible - frontends do not care that backend
they use. While the actual Pods that compose backend set may change, the
frontend clients should not need to be aware of that, nor should they need to
keep track of the set of backend themselves.

The Service abstraction enables this decoupling.

The set of Pods targeted by a Service is usually determined by a selector that
we define.

If the workload speaks HTTP, we might choose to use an Ingress to control how
web traffic reaches that workload. Ingress is not a Service type, but it acts as
the entry point for the cluster. An Ingress lets us consolidate the routing
rules into a single resource, so that we can expose multiple components of the
workload, running separately in the cluster, behind a single listener.

The Gateway APi for k9s provides extra capabilities beyond Ingress and Service.
We can add Gateway to the cluster - it is a family extension of APIs,
implemented using `CustomResourceDefinitions` and then use these to configure
access to network services that are running in the cluster.

## Cloud-native service discovery

If we are able to use k8s APIs for service discovery in the applicaiton, we can
query the API server for matching `EndpointSlices`. K8s updates the 
`EndpointSlices` for a Service whenever the set of Pods in a Service changes.
