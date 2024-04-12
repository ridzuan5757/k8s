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

A probe is a diagnostic performed periodically by the kubelet on a container. To
perform a diagnostic, the kubelet either executes code within the container, or
makes a network request.

### Check mechanisms

There are four different ways to check a container using a probe. Each probe
must define exactly one of these 4 mechanisms:

###### `exec`

Executes a specified command inside the container. The diagnostic is considered
successful if the command exits with a status code of `0`.

###### `grpc`

Perform a remote procedure call using gRPC. The target should implement gRPC
health checks. The diagnostic is considered successful if the `status` of the
response is `SERVING`.

###### `httpGet`

Perform HTTP `GET` request against the pod's IP address on a specified port and
path. The diagnostic is considered successul if the response has a status code
greater than or equal to 200 and less than 400.

###### `tcpSocket`

Perform  TCP check against the pod's IP address on specified port. The
diagnostic is considered successful if the port is open. If the remote system
(the container) closes the connection immediately after it opens, this counts as
healthy.

> Unlike other mechanisms, `exec` probe's implementation involves the creation
> or forking of multiple processes each time when executed. As a requeslt, in
> case of the clusters having higher pod densities, lower intervals of
> `intialDelaySeconds`, `periodSeconds`, configuring any probe with exec
> mechanism might introduce an overhead on the cpu usage of the node. In such
> scenarios, consider using the alternative probe mechanism to avoid the
> overhead.

### Probe outcome

Each probe has one of 3 results:
- `Success` - The container passed the diagnostic.
- `Failure` - The container failed the diagnostic.
- `Unknown` - The diagnostic failed. No action should be taken, and the kubelet
  will make further checks.

### Type of probe

The kubelet can optionally perform and reach to 3 kinds of probes on running
containers:

###### `livenessProbe`

Indicates whether the container is running. If the liveness probe fails, the
kubelet kills the container and the container is subjected to its restart
policy. If a container does not provide a liveness probe, the default state is
`Success`.

###### `readinessProbe`

Indicates whether the container is ready to respond to requests. If the
readiness probe fails, the endpoints controller removes the pod's IP address
from the endpoints of all services that match the pod.

The default state of readiness before the initial delay is `Failure`. If a
container does not provide a readiness probe, the default state is `Success`.

###### `startupProbe`

Indicates whether the application within the container is started. All other
probes are disabled if a startup probe is provided, until it succeeds. If the
startup probe fails, the kubelet kills the container, and the container is
subjected to its restart policy. If a container does not provide startup probe,
the default state is `Success`.

#### When `livenessProbe` should be used.

If the process in the container is able to crash on its own whenever it
encounters an issue or becomes unhealthy, we do not necessarily need a liveness
probe. The kubelet will automatically perform the correct action in accordance
to the pod's restart policy.

If we would like the container to be killed and restarted if a probe fails, then
specify liveness probe, adn specify a `restartPolicy` of `Always` or
`OnFailure`.

#### When `readinessProbe` should be used.

If we would like to start sending traffic to a pod only when a probe succeeds,
specify a readiness probe. In this case, the readiness probe might be the same
as the liveness probe, but the existence of the readiness probe in the spec
means that the pod will start without receiving any traffic and only start
receiving traffic after the probe starts succeeding.

If we want the container to be able to make itself down for maintenance, we can
specify a rediness probe that checks an endpoint specific to readiness that is
different from the liveness probe.

If the app has strict dependency on backend services, we can implement both
liveness and readiness probe. The liveness probe passes when the app itself is
healthy, but the readiness probe additionally checks that each required backend
service available. This helps us avoid directing traffic to pods that can only
respond with error messages.

If the container needs to work on loading large data, configuration files, or
migration startup, we can use startup probe. However, if we want to dectect the
difference between an app that has failed and an app that is still processing
its startup data, a readiness probe might be more suitable.

> If we want to be able to drain requests when the pod is deleted, we do not
> necessarily need a readiness probe; on deletion, the pod automatically put
> itself into an unready state regardless of whether the readiness probe exists.
> The pod remains in the unready state while it waits for the containers in the
> pod to stop.

#### When `startupProbe` should be used.

Startup probes are useful for pods that have containers that take a long time to
come into service. Rather than set a long liveness interval, we can configure a
separate configuration for probing the container as it starts up, allowing a
time longer than the liveness interval would allow.

If the container usually starts in more than `initialDelaySeconds +
failureThreshold * periodSeconds`, we should specify a startup probe that checks
the same endpoint as the liveness probe.

The default for `periodSeconds` is 10s. We whould then set its `failureThreshold` 
high enough to allow the container to start, without changing the default values
of the liveness probe. This helps to protect against deadlocks.

## Termination of pods

Because pods represent processes running on nodes in the cluster, it is
important to allow those processes to gracefully terminate when they are no
longer needed rather than being abruptly stopped with a `KILL` signal and having
no change to clean up.

The design aim is for us to be able to request deletion and know when processes
terminate, but also be able to ensure that deletes eventually complete. When we
request deletion of a pod, the cluster records and tracks the intended grace
period before the pod is allowed to be forcefully killed. With that forceful
shutdown tracking in place the kubelet attempts graceful shutdown.

Typically, with this graceful termination of the pod, kubelet makes requests to
the container runtime to attempt to stop the containers in the pod by first
sending a `SIGTERM` signal, with a grace period timeout, to the main process in
each container. The requests to stop the containers are processed by the container 
runtime asynchronously. There is no guarantee to the order of processing for these
requests.  

Many container runtimes respect the `STOPSIGNAL` value defined in the container
image and, if different, send the container image configured `STOPSIGNAL`
instead of `SIGTERM`. Once the grace period has expired, the `KILL` signal is
sent to any remaining processes, and the pod is then deleted from the API
server.

If the kubelet or the contaienr runtime's management service is restarted while
waiting for processes to terminate, the cluster retries from the start including
the full original grace period.

An example flow:
- `kubectl` is used too manually delete specific pod, with default frace period
  of 30 seconds.
- The pod in the API server is updated with the time beyond which the pod is
  considered "dead" along with the grace period. If `kubectl describe` is used
  to check the pod that is being deleted, that pod will show up as
  `Terminating`. On the node where the pod is running: as soon as the kubelet
  sees that a pod has been marked as terminating (a graceful shutdown duration
  has been set), the kubelet begin the local pod shutdown process.
    - If one of the pod's container has defined a `preStop` hook and the
      `terminatingGracePeriodSeconds` in the pod spec is not set to 0, the
      kubelet runs that hook inside of the container. The default 
      `terminatingGracePeriodSeconds` is 30 seconds.
    - If the `preStop` hook is still running after the grace period expires, the
      kubelet will request a small, one-off grace period extension of 2 seconds.
    - If the `preStop` hook needs longer to complete than the default grace
      period allows, `terminatingGracePeriodSeconds` value must be modified to
      suit this.
    - The kubelet triggers the container runtime to send a `SIGTERM` signal to
      process 1 inside each container.
    - The containers in the pod receive the `SIGTERM` signal at different times
      and in arbitrary order. If the order of shutdowns matters, consider using
      a `preStop` hook to synchronize.
- At the same time as the kubelet is starting graceful shutdown of the pod, the
  control plane evaluates whther to remove that shutting down pod from
  `EndpointSlices` and `Endpoints` object, where those objects represent a Service
  with a configured selector. `ReplicaSets` and other workload resources no
  longer treat the shutting-down pod as a valid, in-service replica.
    - Pods that shut down slowly should not continue to serve regular traffic
      and should start terminating and finish processing open connections. Some
      applications need to go beyond finishing open connections and need more
      graceful termination, for example, session draining and completion.
    - Any endpoints that represent the terminating pods are not immediately
      removed from `EndpointSlices`, and status indicating terminating state is
      exposed from the `EndPointsSlices` API and the legacy `Endpoints` API. 
    - Terminating endpoints always have their `ready` status as `false` (for
      backward compatibility versions before v1.26), so load balancers will not
      use it for regular traffic.
    - If traffic draining on terminating pod is needed, the actual readiness can
      be checked as condition `serving`. 
        - When the grace period expires, the kubelet triggers forcible shutdown.
          The container runtime sends `SIGKILL` to any processes still running in
          any container in the pod. The kubelet also cleans up a hidden `pause`
          container if that container runtime uses one.
        - The kubelet transitions the pod into a terminal phasse (`Failed` or
          `Succeeded` depending on the state of its containers). This step is
          guaranteed since version v1.27.
        - The kubelet triggers forcible removal of pod obect from the API
          server, by setting grace period to 0 (immediate deletion).
        - The API server deletes the pod's API object, which is then no longer
          visible from any client.
        

> If we do not have the `EndpointSliceTerminatingCondition` feature gate enabled
> in the cluster (the gate is on by default from k8s v1.22 and locked to default
> in v1.26), then the k8s control plane removes a pod from any relevant
> EndpointSlices as soon as the pod;s termination grace period begins. The
> behaviour above is described when the feature gate 
> `EndpointSliceTerminatingCondition` is enabled.

> Beginning kwith k8s v1.29, if the pod includes one or more sidecar contaienrs
> (init cotnainer with an `Always` restart policy), the kubelet will delay
> sinding the `SIGTERM` sgnal to these sidecar containers until the last main
> container has fully terminated. The sidecar controllers will be terminated in
> the reverse order they are defined in the pod spec. This ensures that sidecar
> containers continue serving the other containers in the pod until they are no
> longer needed.
>
> Note that slow termination of a main container will also delay the termination
> of the sidecar containers. If the grace period expires before the termination
> process is complete, the pod may enter emergency termination. In this case,
> all remaining containers in the pod will be terminated simultaneously with a
> short grace period.
>
> Similarly, if a pod has a `preStop` hook that exceeds the termination grace
> period, emergency termination may occur. In general, if we have used `preStop`
> hooks to control the termination order without sidecar containers, we can no
> remove them and allow the kubelet to manage sidecar termination automatically.

### Forced pod termination

> Forced deletions can be potentially disruptive for some workloads and their
> pods.

By default, all deletes are graceful within 30 seconds. The `kubectl delete`
command supports the `--grace-priod=<seconds>` option which allows us to
override the default and specify our own value.

Setting the grace period to `0` forcible and immediately deletes the pod from
the API server. If the pod was still running on a node, that forcible deletion
triggers the kubelet to begin immediate cleanup.

> Additional flag `--force` must be specified along with `--grace-period=0` in
> order to perform force deletions.

When a force deltion is performed, the API server does not wait for confirmtatin
from the kubelet that the pod has been terminated on the node it was running on.
It removes the pod in the APO immediately so that a new pod can be created with
the same name. On the node pods that are set to terminate immediately will still
be given a small grace period before being force killed.

> Immediate deletion does not wait for confirmation that the running resource
> has been terminated. The resource may continue to run on the cluster
> indefinitely.


