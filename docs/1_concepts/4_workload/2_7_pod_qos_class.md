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



