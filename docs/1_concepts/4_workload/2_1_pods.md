# Pods

Pods are the smalles deployable units of computing that can be created and
managed in k8s.

A pod is a group of one or more containers, with:
- shared storage
- network resources
- specification for how to run containers

A pod's contents are always co-located and co-scheduled, and run in a shared
context. A pod models can be application-specific "logical host": it contains
one or more application containers which are relatively tightly coupled.

In non-cloud contexts,  applications executed on the same physical or virtual
machine are analogous to cloud applications executed on the same logical host.

> `init containers` - one or more initialization containers that must run to
> completion before any app containers run
> `ephemeral containers` - a type of containers that can be run temporarily in a
> pod

As well as application containers, a pod can contain init containers that run
during pod startup. We can also inject ephemeral containers for debugging a
running pod.

> Container runtime need to be installed into each node in the cluster so that
> the pods can run there.

The shared context of a pod is a set of Linux napespaces, cgroups, and
potentially other facets of isolation - the same things that isolate a
container. Within a pod's context, the individual applications may have further
sub-isolations applied.

A pod is similar to a set of containers with shared snamespaces and shared
filesystem volumes. Pods in k8s are used in 2 main ways:

###### Pods that run a single container

The "one-container-per-pod" model is the most common k8s use case. In this case,
we can think a pod as a wrapper around a single container. k8s manages pods
rather than containers directly.

Grouping multiple co-located and co-managed containers in a single pod is
relatively advanced use case. We should use this pattern only in specific
instances in which the container is tightly coupled.

We do not need to run multiple containers to provide replication for relatively
advaned use case. This pattern should only be used on specific instances in
which the containers are tightly coupled.

We do not need to run muliplt containers to provide replication for reliance and
capacity.

## Using pods
