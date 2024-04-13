# Pod quality of service classes

k8s classifies the pods that we run and allocates each pod into a specific
**quality of service (QoS) class**. k8s uses that classification to influence
how different pods are handled. k8s does this classification based on the
resource requests of the containers in that pod, along with how those requests
relate to resource limits. This is known as QoS class.

k8s assigns every pod a qos class based on the resource requests and limits of
its component containers. qos classes are used by k8s to decide which pods to
evict from node experiencing node pressure.

The possible qos classes are:
- `Guaranteed`
- `Burstable`
- `BestEffort`

When a node runs out of resources, k8s will first evict `BestEffort` pods
running on that node, followed by `Burstable` and finally `Guaranteed` pods.

When this eviction is due to resource pressure, only pods exceeding resource
requests are candidates for eviction.

### `Guaranteed`

Pods that are `Guaranteed` has the strictest resource limits and are least
likely to face eviction. They are guaranteed not to be killed until they exceed
their limits or there are no ower-priority pods that can be preempted from the
node. They may not acquire resources beyond the specified limits. These pods can
also make use of exclusive CPUs using the `static` CPU management policy.

#### Criteria
- Every container in the pod must have a memory limit and a memory request.
- For every container in the pod, the memory limit must equal to the memory
  request.
- Every container in the pod must have a CPU limit and a CPU request.
- For every container in the pod, the CPU limit must equal to the CPU request.

### `Burstable`

Pods that are `Burstable` have some lower-bound resource guarantees based on the
request, but does not have a specific limit. If a limit is not specified, it
defaults to a limit equivalent to the capacity of the node, which allows the
pods to flexibly increase their resources if resources are available. In the
event of Pod eviction due to Node resource pressure, these pods are evicted only
after all `BestEffort` pods are evicted. Because a `Burstable` pod can include a
container that has no resource limits or requests, a pod that is `Burstable` can
try to use any amount of node resource.

#### Criteria
- The pod does not neet the criteria for qos class `Guaranteed`.
- At least one container in the pod has memory or cpu request or limit.

### `BestEffort`

Pods in `BestEffort` qos class can use node resources that are not specifically
assigned to pods in other qos classes. For example, if we have a node with 16
cpu cores available to the kubelet, and we assign 4 cpu cores to a `Guaranteed`
pod, then a pod in `BestEffort` qos can try to use any amount of the remaining
12 cpu cores.

The kubelet preers to evict `BestEffort` pods if the node comes under resource
pressure.

#### Criteria

A pod has a qos class of `BestEffort` if it does not meet the criteria for
either `Guaranteed` or `Burstable`. In other words, a pod is `BestEffort` only
if none of the containers in the pod have a cpu limit or cpu request. Containers
in a pod can request other resources (not cpu or memory) and still be classified
as `BestEffort`.

