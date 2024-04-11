# Pod Lifecycle

Pod follow a defined lifecycle, starting in the `Pending` phase, moving through
`Running` if at least one of its primary containers starts OK, and then through
either the `Succeeded` or `Failed` phases depending on whether any container in
the pod terminated in failure.

Whilst a pod is running, the kubelet is able to restart containers to handle
some kind of faults. Within a pod, k8s tracks different container states and
determines what action to take to make the pod healthy again.

In the k8s API, pods have both a specification and an actual status. The status
for a pod object consists of a set of Pod conditions. We can also inject custom
readiness information into the condition data for a pod, if that is useful to
the application.

Pods are only scheduled once in their lifetime. Once a pod is scheduled or
assigned to a node, the pod runs on that node until it stops or is terminated.

## Pod lifetime

Similar as individual application containers, pods are considered to be
relatively ephemeral rather than durable entities. Pods are created, assigned
unique ID UID, and sceduled to nodes where they remain until termination
according to restart policy or deletion. If a node dies, the pods shceduled to
that node are scheduled for deletion after a timeout period.

Pods do not, by themselves, self-heal. If a pod is scheduled to a node that then
fails, the pod is deleted. Likewise, a pod would not survive a eviction due to
lack of resource or note maintenance. k8s use a higher level abstraction called
controller that handles the work of managing the relativel disposable pod
instances.

A given pod as defined by a UID is never "rescheduled" to a different node;
instead, that pod can be replaced by a new, near identical pod, with even the
same name if desired, but different UID.

When soemthing is said to have the same lifetime as a pod, such as a volume,
that means that the thing exists as long as specific pod with that exact UID
exists. If that pod is deleted for any reason, and even if an identical
replacement is created, the related thing (a volume for example) is also
destroyed and created new.

## Pod phase

A pod's `status` feld is a `PodStatus` object, which has a `phase` field. The
phase of a pod is a simple, high-level summary of where the pod is in its
lifecycle. The phase is not intended to be a comprehensive rollup of
observations of container or pod state, not is it intended to be a comprehensive
state machine.

The number and meanings of pod phase values are tightly guarded. Other than what
is documented here, nothing should be assumed about pods that have a given
`phase` values.

Possible `phase` values:

### `Pending`

The Pod has been accepted by the Kubernetes cluster, but one or more of the 
containers has not been set up and made ready to run. This includes time a Pod 
spends waiting to be scheduled as well as the time spent downloading container 
images over the network.

### `Running`

The Pod has been bound to a node, and all of the containers have been created. 
At least one container is still running, or is in the process of starting or 
restarting.

### `Succeeded`

All containers in the pod have terminated in success, and will not bet
restarted.

### `Failed`

All containers in the pod have terminated, and at least on container has
terminated in failure. That is, the container either exited with non-zero status
or was terminated by the system.

### `Unknown`

For some reason the state of the pod could not be obtained. This phase typically
occurs due to an error in communicating with the node where the pod should be
running.

> When a pod is being deleted, it is shown as `Terminating` by some kubectl
> commands. This `Terminating` status is not one of the pod phases. A pod is
> granted a term to terminate gracefully, which defaults to 30 seconds. Flag
> `--force` can be used to terminate a pod by force.

Since k8s v1.27, the kubelet transitions deleted pods, execept for static pods,
and force-deleted pods without a finalizer, to a terminal phase `Failed` or
`Succeeded` depending on the exit statuses of the pod cotainers before their
deletion from the API server.

If a node dies or is diconnected from the rest of the cluster, k8s applies a
policy for setting the `phase` of all pods on the lost node to failed.

## Container states

As well as the phase of the pod overall, k8s tracks the state of each container
inside a pod. We can use container lifecycle hooks to trigger events to run at
certain points in a container's lifecycle.

Once the scheduler assigns a pod to anode, the kubelet starts creating
containers for that pod using a container runtime. There are 3 possible
container states:
- `Waiting`
- `Running`
- `Terminated`

To check the state of a pod's containers, we can use `kubectl describe pod 
<pod_name>`. The output shows the state for each container within that pod.

### `Waiting`

If a container is not in either the `Running` or `Terminated` state, it is
`Waiting`. A container in the `Waiting` state is still running the operations it
requires in order to complate start up.

For example, pulling the container image from a container image registry, or
applying Secret data. When we use `kubectl` to query a pod with a container that
is `Waiting`, we also see a `Reason` field to summarize why the container is
in that state.

### `Running`

The `Running` status inidcates that a container is executing without issues. If
there was a `postStart` hook configured, it has already executed and finished.
When we use `kubectl` to query a pod with a container that is `Running`, we also
see information about when the container entered the `Running` state.

### `Terminated`

A container in the `Terminated` state began execution and then either ran to
completion or failed for some reason. When we use kubectl to query a pod with a
container that is `Terminated`, we see a reason, exit code and the start and
finish time for that container's period of execution.

If a container has a `preStop` hook configured, this hook runs before the
container enters the `Terminated` state.

## Container restart policy

The `spec` of a pod has a `restartPolicy` field with possible values:
- `Always`
- `OnFailure`
- `Never`

The default value is `Always`.

The `restartPolicy` for a pod applies to app containers in the pod and to
regular "init containers". "Sidecar contaienrs" ignore the pod-level
`restartPolicy` filed: in k8s, a sidecar is defined as an entry inside
"init Containers" that has its container-level `restartPolicy` set to `Always`.
For "init containers" that exit with an error, the kubelet restarts the "init
container" if the pod level `restartPolicy  is either `OnFailure` or `Always`.

When the kubelet handling container restarts according to the configured restart
policy, that only applies to restarts that make replacement containers inside
the same Pod and running on the same node.

After containers in a pod exit, the kubelet restarts them with an exponential
back off delay {10s, 20s, 40s, ...}, that is capped at 5 minnutes. Once a
container has executed for 10 minutes without any problems, the kbuelet resets
the restart backoff timer for that container.

## Pod conditions

A pod has `PodStatus`, which has an array of `PodConditions` through which the
pod has or has not passed. Kubelet manages the following `PodConditions`:
- `PodScheduled` - the pod has been scheduled to a node.
- `PodReadyToStartCotnainers` - the pod sandbox has been successfully configured
  and networking configured.
- `ContainersReader` - all containers in the pod are ready.
- `Initialized` - all "init containers" have completed successfully.
- `Ready` - the pod is able to serve requests and should be added to the load
  balancing pools of all matching services.

### Conditions field name and description
- `type` - Name of this pod condition.
- `status` - Indicates whether that condition is applicable, with possible
  values:
    - `True`
    - `False`
    - `Unknown`
- `lastProbeTime` - Timestamp for when the pod condition was last probed.
- `lastTransitionTime` - Timestamp for when the pod last transitioned from one
  status to another.
- `reason` - Machine readable upper camel case text indicating the reason for
  the condition's last transition.
- `message` - Human readable message indicating details about the last status
  transition.

### Pod readiness

The application can inject extra feedback or signals into `PodStatus:Pod`
readiness. To use this, set `readinessGates` in the pod's `spec` to specify a
list of additional conditions that the kubelet evalues for pod readiness.

Readiness gates are determined by the current state of `status.condition` fields
for the pod. If k8s cannt find such condition in the `status.conditions` field
of a pod, the status of the condition is defaulted to `False`.

```yaml
kind: Pod
...
spec:
  readinessGates:
    - conditionType: "www.example.com/feature-1"
status:
  conditions:
    - type: Ready                              # a built in PodCondition
      status: "False"
      lastProbeTime: null
      lastTransitionTime: 2018-01-01T00:00:00Z
    - type: "www.example.com/feature-1"        # an extra PodCondition
      status: "False"
      lastProbeTime: null
      lastTransitionTime: 2018-01-01T00:00:00Z
  containerStatuses:
    - containerID: docker://abcd...
      ready: true
...
```
The pod conditions we add must have names that meet the k8s label key format.

### Status for pod readiness

The `kubectl patch` command does not supporting patching object status. To set
these `status.conditions` for the pod, applications and operators should use the
`PATCH` action. We can use a k8s client library to write code that sets custom
pod conditions for pod readiness.

For a pod that uses custom conditions, that pod is evaluated to be ready
**only** when both the following statements apply:
- All containers in the pod are ready.
- All conditions specified in `readinessGates` are `True`.

When a pod's containers are ready but at least one custom condition is missing
or false, the kubelet set the pod's condition to `ContainersReader`.

### Pod network readiness

> During its early development, this condition was named `PodHasNetwork`.

After a pod get scheduled on a node, it needs to be admitted by the kubelet and
to have nay required storage volumes mounted. Once these pahses are complete,
the kubelet works with a container runtime CRI to set up a runtime sandbox and
configure networking for the pod.

If the `PodReadyToStartContainersCondition` feature gate is enabled (enabled by
default on v1.29 onwards) the `PodReadyToStartContainers` condition will be
added to the `status.conditions` field of a pod.

The `PodReadyToStartContainers` condition is set to `False` by the kubelet when
it detects a pod does not have a runtime sandbox with networking configured.
This occurs in the following scenarios:
- Early in the lifecycle of the pod, when the kubelet has not yet begun to set
  up a sandbox for the pod using the container runtime.
- Later in the lifecycle of the pod, when the pod sandbox has been destroyed due
  to either:
  - The node rebooting without the pod getting evicted.
  - For container runtimes that use virtual machines for isolation, the pod
    sandbox virtual machine rebooting, which then requires creating a new
    ssandbox and fresh container network configuration.

The `PodReadyToStartContainers` condition is set to `True` by the kubelet after
the successful completion of sandbox creation and network configuration for the
pod by the runtime plugin. The kubelet can start pulling container images and
create containers after `PodReadyToStartContainers` condition has been set to
`True`.

For a pod with "init containers", the kubelet sets the `Initiaized` condition to
`True` after the init containers have successfully completed which happens after
successful sandbox creation and network configuration by the runtime plugin. For
a pod without "init containers", the kubelet sets the `Initialized` condition to
`True` before sandbox creation and network configuration starts.

## Container probes


