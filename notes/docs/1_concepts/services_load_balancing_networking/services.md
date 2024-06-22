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

The Gateway API for k9s provides extra capabilities beyond Ingress and Service.
We can add Gateway to the cluster - it is a family extension of APIs,
implemented using `CustomResourceDefinitions` and then use these to configure
access to network services that are running in the cluster.

## Cloud-native service discovery

If we are able to use k8s APIs for service discovery in the application, we can
query the API server for matching `EndpointSlices`. K8s updates the 
`EndpointSlices` for a Service whenever the set of Pods in a Service changes.

For non-native applications, k8s offers way to place network port or load
balancer in between the application and the backend Pods.

Either way, the workload can use these service discovery mechanism to find the
target it wants to connect to.

## Service Definition

A Service is an object the same way that a Pod or ConfigMap is an object. We can
create, view or modify Service definitions using the k8s API. Usually we use a
tool such as `kubectl` to make those API calls for us.

For example, suppose we have a set of Pods that each listen on TCP port 9376 and
are albelled as `app.kubernetes.io/name=MyApp`. We can define a Service to
publish that TCP listener.

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

Applying this manifest creates a new Service named "my-service"m with the
default ClusterIP service type. The Service targets TCP port 9376 on any Pod
with the `app.kubernetes.io/name: MyApp` label.

K8s assigns this Service an IP address (the cluster IP), that is used by the
virtual IP address mechanism.

The controller for that Service continuously scans for Pods that match its
selector, and then makes any necessary updates to the set of EndPointSliaces for
the Service.

The name of a Service object must be a valid RFC 1035 label name.

> [!NOTE]
> A Service can map any incoming port to a targetPort. By default and for
> convenience, the targetPort is set to the same value as the port field.

## Port definitions

Port definitions in Pods have names, and we can reference these names in the
`targetPort` attribute of a Service. For example, we can bind the targetPort of
the Service to the Pod port in the following way:

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

The default protocol for Services is TCP; we can also use any other supported
protocol:
- SCTP
- TCP (default)
- UDP

Because many Services need to expose more than one port, k8s supports multiple
port definitions for a single Service. Each port definition can have the same
protocol, or a different one.

## Services without selectors

Services most commonly abstract access to k8s Pods thanks to the selector, but
when used with a corresponding set of `EndpointSlices` objects and without a
selector, the `Service` can abstract other kinds of backends, including ones that
run outside the cluster.

For example:
- We want to have an external database cluster in production, but in the test
  environment we use our own databases.
- We want to point the `Service` to a `Service` in a different `Namespace` or on
  another cluster.
- Migrating workload to k8s. While evaluating the approach, we might only run
  portion of the backends in k8s.

In any of these scenarios, we can define a `Service` without specifying a
selector to match Pods.

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

Because this `Service` has no selector, the corresponding `EndpointSlice` and
legacy `Endpoints` objects are not created automatically. We can map the
`Service` to the network address and port where it is running by adding a
`EndpointSlice` object manually. 

```yaml
apiVersion: discovery.k8s.io/v1
kind: EndpointSlice
```
