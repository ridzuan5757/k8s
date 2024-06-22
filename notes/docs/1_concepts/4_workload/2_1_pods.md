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

## Pod update and replacement

When the pod template for a workload resource is changed, the controller creates
new pods based on the updated template instead of updating the existing pods.

k8s does not prevent us from managing pods directly. It is possible to update
some fiels of a running pod in place. However, pod update operations like `patch` 
and `replace` have some limitations:
- Most of the metadata about pod is immutable. We cannot change the field value:
    - `namespace`
    - `name`
    - `uid`
    - `creationTimestamp`
    - `generation` field is unique. It only accepts updates that increment the
      field's current value.
- If the `metadata.deletionTimestamp` is set, no new entry can be added to the
  `metadata.finalizers` list.
- Pod updates may not change fields other than:
    - `spec.containers[*].image`,
    - `spec.initContainers[*].image`
    - `spec.activeDeadlineSeconds`
    - `spec.tolerations` - For this field value, we can only add new entries.
- When pdating the `spec.ActiveDeadlineSeconds` field, 2 types of updates are
  allowed:
    - setting the unassigned field to positive number.
    - updating the field from positive number to a smaller non-negative number.

## Resource sharing and communication

Pods enable data sharing and communication among their constituent containers.

### Storage in pods

A pod can specify a set of shared storage volumes. All containers in the pod can
access the shared volumes, allowing those containers to share data. Volumes also
allow persistent data in a pod to survive in a case one of the containers within
needs to be restarted.

### Pod networking

EAch pod is assigned a unique IP address for each address family. Every
container in a pod shares the network namespace, including the IP address and
network ports. Inside a pod, and only then the containers that belong to the pod
can communicate with one aother using `localhost`.

When containers in a pod communicate with entities outside the pod, they must
coordinate how they use the shared network resources such as ports. Within a
pod, containers share an IP address and port space, and can find each other via
`localhost`.

The containers in a pod can also communicate with each other using standard
inter-process communications like SystemV semaphores or POSIX shared memory.
Containers in different pods have distinct IP addresses and can not communicate
by OS-level IPC without special configuration. Containers that want to interact
with a container running in a different pod can use IP networking to
communicate.

Containers within the pod see the system hostname as being the same as the
configured `name` for the pod.

## Priveleged mode for containers

> The container runtime must support the concept of a privileged container for
> this setting to be relevant.

Any container in a pod can run in privileged mode to use operating system
administrative capabilities that would otherwise be unaccessible. This is
available for both Windows and Linux.

### Linux privileged containers

A ny container in a pod can enable privileged mode using the `privileged` flag
on the security context of the container spec. This is useful for containers
that want to use operating system administrative capabilities such as
manipulating network stack or accessing hardware devices.

### Window privileged containers

We can create windwos hostprocess pod by setting `windowsOption.hostProcess`
flag on the security context of the pod spec. All containers in these pods must
run as Windows HostProcess containers. HostProcess pods run directly on the host
and can also be used to perform administrative tasks as is done with Linux
privileged containers.

## Static pods

Static pods are managed directly by the kubeet daemon on a specific node without
the API server observing them. Whereas most pods are managed by the control
plane (for example, Deployment), for static pods, the kubelet directly supervises
each static pod (and restarts it if it fails).

Static pods are always bound to one kubelet on specific node. The main use for
static pods is to run a self-hosted control plane; in other words, using kubelet
to supervise the individual control plane components.

The kubelet automatically tries to create a mirror pod on the k8s API server for
each static pod. This means that the pods running on a node are visible on the
API server, but cannot be controlled from there.

> The `spec` of static pod cannot refer to other API objects such as
> ServiceAccount, ConfigMap, Secret etc.

## Pods with multiple containers

Pods are designed to support multiple cooperating processes as contaienrs that
form a cohesive unit of service. The containers in a pod are automatically
co-located and co-scheduled on the same physical or virtual machine in the
cluster.

The containers can share resources and dependencies, communicates with one
another and coordinate when they are terminated.

Pods in k8s cluster are used in 2 main ways:

### Pods that run a single container

The "one-container-per-pod" model is the most common k8s use case; in this case
a pod is acting as a wrapper around a single container. k8s manages pods rather
than managing the containers directly.

### Pods that run multipe containers that need to work together

A pod can encapsulate an application composed of multiple co-located containers
that are tightly coupled and need to share resources. These co-located
containers form a single cohesive unit of service.

For example, one container serving data stored in a shared volume to the public
while a separate sidecar container refreshes or updates those files. The pod
wraps these containers, storage resources, and an ephemeral network identity
together as a single unit.

Some pods have init containers as well as app containers. By default, init
containers run and complete before the app containers are started. Sidecar
containers is also possible for providing auxiliary services to the main
appliation pod.

### Feature State

Enabled by default, the `SidecarContainers` feature gate allows us to specify
`restartPolicy: Always` for init containers. Setting the `always` restart policy
ensures that the containers where we set it are treated as sidecars that are
kept running during the entire lifetime of the pod. Containers that we explicily
define as sidecar containers start up before the main applicationn pod and
remain running until the pod is shut down.

## Container probes

A probe is a diagnostic performed periodically by the kubelet on a container. To
perform a diagnostic, the kubelet can invoke different actions:
- `ExecAction` performed with the help of the container runtime
- `TCPSocketAction` checked directly bu the kubelet
- `HTTPGetAction` checked directly by the kubelet
