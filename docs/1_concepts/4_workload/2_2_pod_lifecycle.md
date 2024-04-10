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

|Value|Description|
|---|---|
|`Pending`| The pod has been accepted by the k8s cluster, but one or more of the
containers has not been set up and made ready to run. This includes time a pod
spends waiting to be scheduled as well as the time spent downloading container
images over the network.|
|`Running`|The pod has been bound to a node, and all of the containers have been
created. At least one container is still running, or is in the process of
starting or restarting.|


