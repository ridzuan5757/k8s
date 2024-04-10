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

The following is an example of a pod which consists of a container running the
image `nginx:latest`

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80
```

To create a pod shown above, run the following command:

```bash
kubectl apply -f https://k8s.io/examples/pods/simple-pod.yaml
```

Pods are generally not created directly and are created using workload
resources.

### Worload resources for managing pods

Usually we do not need to create pods directly, even singleton pods. Insteadm we
create them using workload resources such as Deployment or Jon. If the pods need
to track state, consider the StatefulState resource.

Each pod is meant to run a single instance of given application. If we want to
scale the applicaiton horizontally to provide more overall resources by running
more instances, we should use multiple pods, one for each instance.

In k8s, this is typically referred as replication. Replicated pods are usually
created and managed as a group by a workload resource and its controller.

Pods natively provide 2 kinds of shared resources for their constituent
containers:
- network
- storage

## Working with pods

Most of the time, we will not creating individual pods directly in k8s, even a
singleton pods. This is because pods are designed as relatively ephemeral,
disposable entities. When a pod gets created either directly by ourselves or
indirectly by the controller, the new pod is scheduled to run on a node in the
cluster. The pod remains on that node until either:
- The pod finishes execution
- The pod object is deleted
- The pod is evicted for lack of resources
- The node is failing

> Restartng a container in a pod should not be confused with restarting a pod. A
> pod is not a process, but an environment for running containers. A pod
> persists until it is deleted.

The name of a pod must be a valid DNS subdomain value, but this can produce
unexpected resouts for the pod hostname. For best compatibility, the name should
follow the more restrictive reuls for a DNS label.

### Pod OS

We should set the `.spec.os.name` field to either `windows` or `linux` to
indicate the OS on which we want the pod to run. These 2 ae the only operating
systems supported for now by k8s.

In k8s v1.29, the value set for this field has no effect on scheduling of the
pods. Setting the `.spec.os.name` helps to identify the pod OS authoritatively
and is used for validation.

The kubelet refuses to run a pod where we have specified a pod OS, if this is
not the same as the operating system for the node where that kubelet is running.
The pod security standards also use this field to avoid enforcing policies that
are not relevant to the operating system.

### Pods and controllers

We can use workload resources to create and manage multiple pods. A controller
for the resource handles replication and rollout and automatic healing in case
of a pod failure. For example if a node fails, a controller notices that pods on
that node have stopped working and creates replacement pod. The scheduler places
the replacement pod onto a healthy node.

Here are some examples of workload resurces that manage one or more pods:
- Deployment
- StatefulSet
- DaemonSet

### Pod templates

Controllers for workload resources create pods from a pod template and manage
these pods on our behalves.

Pod templates are specification for creating pods, and are included in workload
resources such as Deployments, Jobs and DaemonSets.

Each controller for a workload resource use the `PodTemplate` inside the
wokrload object to make actual pods. The `PodTemplate` is part of the desired
state of whatever workload resource used to run the application.

The sample below is a manifest for a simple job with a `template` that starts
one container. The contaiiner in that pod prints a message then pause.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: hello
spec:
  template:
    # This is the pod template
    spec:
      containers:
      - name: hello
        image: busybox:1.28
        command: ['sh', '-c', 'echo "Hello, Kubernetes!" && sleep 3600']
      restartPolicy: OnFailure
    # The pod template ends here
```

Modifying the pod template or switching a new pod template has no direct effect
on the pods that already exist. If we are changing the pod template for a
workload resource, that resource needs to create replacement pods that use the
updated emplate.

For example, the StatefulSet controller ensures that the running pods match the
current pod template for each StatefulSet object. If we edit the StatefulSet to
change its pod template, the StatefulSet starts to create new pods based on the
updated template. Eventually, all of the old pods are replaced with new pods and
the update is complete.

Each workload resource implements its own rules for handling changes to the pod
template.

On nodes, the kubelet does not directly observe or manage any of the details
around pod templates and updates; those details are abstracted away That
abstraction and separation of concerns simplifies system semantics, and makes it
feasible toe xtend the cluster's behaviour without changing existing code.


