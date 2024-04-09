# Nodes
- Workload : application running on k8s.
- Pod: set of running containers in a cluster.
- k8s runs the workload by placing containers into pods to run on nodes.
- Node could be virtual or physical machine, depending on cluster.
- Each node is managed by control plane and contains service necessary to run
  pods.

## Node components
- kubelet
- container runtime
- kube-proxy

## Management
- 2 main ways to have node added to api server:
    - kubelet on a node self-registers to a control plance.
    - manually adding node object
- After node object is created or kubelet on a node self-registers, the control
  plane checks whether the new node object is valid.
- k8s creates node object internally. k8s checks that a kubelet has registered
  to the API server that matches `metadata.name` field of the node. If the node
  is healthy (all necessary services are running), then it is eligible to run a
  pod. Otherwise, that node is ignored for any cluster activity until it
  becomes helathy.

## Node name uniqueness
- `name` identifies a node.
- 2 nodes cannot share the same name at the same time.
- k8s also assumes that resource with the same name is pointing to the same
  object. In case of a node, it is implicitly assumed that na instance using the
  same name will have the same state:
    - network settings
    - root disk contents
    - node labels
- this may lead to inconsistencies if an instance was modified without changing
  its name. 
- if the node needs to be replaced or updated significantly, the existing node
  object needs to be removed from the API server first and re-added after the
  update.

## Node self-registration
- When kubelet flag `--register-flag` is true, the kubelet will attempt to
  register itself with the API server. This is the preferred pattern used.

## Manual node administration
- We can create and modify node objects using kubectl.
- When we create node object manually, set `--register-node=false`.
- We can modify node object regarless of the setting of `register-node`. For
  example we can set labels on an existing node and mark it unschedulable.
- We can use labels on nodes in conjuction with node selectors on pods to
  control scheduling. We can constrain a pod to only eligible to run on a subset
  of the available nodes.
- Marking a node as unschedulable prevents the scheduler form placing new pods
  onto that node but does not affect existing pods on the node. This is useful
  as a preparatory step before a node reboot or other maintenance.
- To mark a node unschedulable:

```bash
kubectl cordon <node_name>
```

## Node status
- Node status contains the following information:
    - Addresses
    - Conditions
    - Capacitu and allocatable
    - Info
- We can use `kubectl` to view node status and other details

```bash
kubectl describe node <node_name>
```

## Node heartbeats
- Heartbeats send by k8s nodes help our cluster determine the availability of each
nodes and to take action when failures are detected.
- For nodes, there are 2 forms of heartbeats:
    - Updates to the `.status` of a node.
    - `Lease` objects within `kube-node-lease` namespace. Each node has an
      associated lease object.

## Node controller
- CIDR block: Classless inter-domain routing - method for allocating IP
  addresses for IP routing.
- Node controller (control loops) is a k8s control plane component that manages
  various aspects of nodes.
- The node controller has multiple roles in a node's life. 
    - Assigning CIDR block to the node when it is registered if the CIDR assignment
    is turned on.
    - Keeping the node controller's internal list of nodes up to date with the
      cloud provider;s list of available machines. When running a cloud
      environment and whenever node is unhealthy, the node controller asks the
      cloud provider if the VM for that node is still available. If not, the
      node controller deletes the node from the list of the nodes.
    - Monitoring node's health:
        - In the case that a node becomes unreachable, updating `Ready`
          condition in the node's `.status` field. In this case the node
          controller sets the `Ready` condition to `Unknown`.
        - If a node remains unreachable: triggering API-initiated eviction for
          all of the pods on the unreachable node. By default, the node
          controller waits 5 minutes between marking the node as `Unknown` and
          submitting the first eviction request.
- By default, the node controller chesk the state of each node every 5 seconds.
  This period can be configured using `--node-monitor-period` flag on the
  `kube-controller-manager` component.

# Rate limtis on eviction
- In modes cases, the node controller limits the eviction rate to 
  `--node-eviction-rate` per second (default is 0.1). It won't evict pods from 
  more than 1 node per 10 seconds.
- The node eviction behaviour changes when a node in a given availability zone
  become sunhealthy. The node controller checks what percentage of nodes in the
  zone are unhalthy at the same time.
    - If the fraction of unhealthy nodes is at least
      `--unhealthy-zone-threshold` (default is 0.55), then eviction rate is
      reduced.
    - If the cluster is small `--large-cluster-size-threshold` (default to 0.55),
      then evictions are stopped.
    - Otherwise, the eviction rate is reducted to
      `--secondary-node-eviction-rate` that is default to 0.01 per second.
- The reason these policies are implemented per availability zone is because one
  availability zone might become paritioned from the control plane while the
  others remain connected. If our cluster does not span multiple cloud provider
  availability zones, then the eviction mechanism does not take per-zone
  unavailability into account.
- A key reason for spreading nodes across availability zones is so that the
  workload can be shifted to healthy zones when one entire zone goes down.
  Therefore, if all nodes in a zone is unhealthy, then the node controller
  evicts at the normal rate of `--node-eviction rate`. The corner case is when
  all zones are completely unhealthy (none of the nodes in the cluster are
  healthy). In such case, the node controller assumes that there is some rpoblem
  with connectivity between the control plane and the nodes, and does not
  perform any evictions. If there has been outage and some nodes reappear, the
  node controller does evic pods from the remaining nodes that are unhealthy or
  unreachable.
- The node controller is also responsible for evicting pods running on nodes
  with `NoExecute` taints, unless those pods tolerate that taint. The node
  controller also add taints corresponding to node problems like node
  unreachable or not ready. This means that the scheduler would not place pods
  onto unhealthy nodes.

## Resource capacity tracking
- Node objects track information about the node's resource capcity. For example,
  the amount of memory available and the number of CPUs. Nodes that self-register
  report their capacity during registration. If we manually add a node, then we
  need to set the node's capacity information when we add it.
- The k8s scheduler ensures that there are enough resources for all the pods on
  a node. The scehduler checks that the sum of requests of containers on the
  node is no greater than the node's capacity. The sum of requests includes all
  containers managed by the kubelet, but excludes any containers start directly
  by the container runtime, and also excludes any processes running outside of
  the kubelet's control.

## Node topology
- If we have enabled `TopologyManager` feature gate, then the kubelet can use
  topology hints when making resource assignment decisions.

## Graceful node shutdown
- The kubelet attempts to detect node system shutdown and terminates pods
  running on the node.
- Kubelet ensures that pods follow the normal pod termination process during the
  node shutdown. During the node shutdown, the kubelet does not accept new pods,
  even if those pods are already bound to the node.
- The graceful node shutdown feature depends on systemd since it takes advantage
  of systemd inhibitor locks to delay the node shutdown with a given duration.
- Graceful node shutdown is controlled with the`GracefulNodeShutdown` feature
  gate which is enabled by default in v1.21.
- By default, both configuration options described below, `shutdownGracePeriod`
  and `shutdownGracePeriodCriticalPods` are set to zero, thus not activating the
  graceful node shutdown functionality. To activate the feature, the two kubelet
  config settings should be configured appropriately and set to non-zero values.
- Once the systemd detects or notifes node shutdown, the kubelet sets a
  `NotReady` condition on the node, with the `reason` set to `"node is shutting
  down"`. The kube-scheduler honors this condition and does not schedule any
  pods onto the affected node; other third-party schedulers are expected to
  follow the same logic. This means that new pods would not be scheduled onto
  that node and therefore none will start.
- kubelet also rejects pods during the `PodAdmission` phase if an ongoing node
  shutdown has been detected, so that even pods with a toleration for 
  `node.kubernetes.io/not-ready:NoSchedule` do not start there.
- At the same time when kubelet is setting that condition on its node via the
  API, the kubelet also begins terminating any pods that are running locally.
- During graceful shutdown, kubelet terminates pods in 2 phases:
    - Terminate regular pods running on the node.
    - Terminal critical pods running on the node.
- Graceful node shutdown feature is configured with 2 `KubeletConfiguration`
  options:
    - `shutdownGracePeriod` - specifies the total duration that the node should
      delay the shutdown by. this is the total grace period for pod termination
      for both regular and critical pods.
    - `shutdownGracePeriodCriticalPods` - specifies the duration used to
      terminal critical pods during node shutdown. Value should be less than
      `shutdownGracePeriod`.
    - For example, if `shutdownGracePeriod=30s` and 
     `shutdownGracePeriodCriticalPods=10s`, kubelet will delay the node shutdown
     by 30 seconds. During the shutdown, the first 20 seconds would be reserved
     for gracefully temrinating normal pods, and the last 10 seconds would be
     reserved for the termination of critical pods.
- There are cases when node termination was cancelled by the system or perhaps
  manualler by an administrator. In either of those situations, the node will
  return to the `Ready` state. However, pods which already started the process
  of termination will not be restored by kubelet and will need to be
  rescheduled.
- When pods were evicted during the graceful node shutdown, they are marked as
  shutdown. Running `kubectl get pods` shows the status of the evicted pods as
  `Terminated`. And `kubectl describe pod` indicates that the pod was evicted
  becuase of node shutdown.

## Pod priority based graceful node shutdown.
- To provide more flexibility during graceful node shutdown around the ordering
  of pods during shutdown, graceful node shutdown honors the `PriorityClass` for
  pods, provided that we enabled this feature in our cluster. The feature allows
  cluster administrator to explicityly define the ordering of pods during
  graceful node shutdown based on priority classes.
- The graceful node shutdown feature, shut downs pods in 2 phases, non-critical
  pods, followed by critical pods. If additional flexibility is needed to
  explicitly define the ordering of pods during shutdown in a more granular way,
  pod priority based graceful shutdown can be used.
- When graceful node shutdown honors pod priorities, this makes it possible to
  do graceful node shutdown in multiple phases, each phase shutting down a
  particular priority class of pods. The kubelet can be configured with the
  exact phases and shutdown time per phase.
- Assuming the following custom pod priority classes in a cluster:

|Pod priority class name|Pod priority class value|
|---|---|
|custom-class-a|100000|
|custom-class-b|10000|
|custom-class-c|1000|
|regular/unser|0|

- Within the kubelet configuration the settings for
  `shutdownGracePeriodPodPriority` could look like:

|Pod priority class value|Shutdown period|
|---|---|
|100000|10 seconds|
|10000|180 seconds|
|1000|120 seconds|
|0|60 seconds|

- The corresponding kubelet config yaml file would be:
```yaml
shutdownGracePeriodByPodPriority:
  - priority: 100000
    shutdownGracePeriodSeconds: 10
  - priority: 10000
    shutdownGracePeriodSeconds: 180
  - priority: 1000
    shutdownGracePeriodSeconds: 120
  - priority: 0
    shutdownGracePeriodSeconds: 60
```
- The above table implies that any pod with priority value 100000 will get just
  10 seconds to stop, any pod, any pod with value between 100000 and 10000 will
  get 180 seconds to stop, any pod with value 1000 and 10000 will get 120
  seconds to stop. Finally, all other pods will get 60 seconds to stop.
- One does not have to specify values corresponding to all of the classes. For
  example, we could instead use these settings:

|Pod priotiy class value|Shutdown Period|
|---|---|
|100000|300 seconds|
|1000|120 seconds|
|0|60 seconds|

- In the above case, the pods with `custom-class-b` will go into the same bucket
  as `custom-class-c` for shutdown.
- If there are no pods in a particular range, then the kubelet does not wait for
  pods in that priority range. Instead, the kubelet immediately skips to the
  next priority class value range.
- If this feature is enabled and no configuration is provided, then no ordering
  action will be taken.
- Using this feature requires enabling the
  `GracefulNodeShutdownBasedOnPodPriority` feature gate, and setting
  `ShutdownGracePeriodPodPriority` in the kubelet config to the desired
  configuration containing the pod priority class values and their respective
  shutdown periods.
- Metrics `graceful_shutdown_start_time_seconds` and
  `graceful_shutdown_end_time_seconds` are emitted under the kubelet subsystem
  to monitor node shutdowns.

# Non-graceful node shutdown handling

- A node shutdown action may not be detected by kubelet's Node Shutdown Manager,
  either because the command does not trigger the inhibitor locks mechanism used
  by kubelet or because of a user error, i.e, the `ShutdownGracePeriod` and
  `shutdownGracePeriodCriticalPods` are not configured properly.
- When a node is shutdown but not detected by kubelet's Node Shutdown Manager, the
  pods that are part of a `StatefulSet` will be stuck in terminating status on the
  shutdown node and cannot move to a new running node. This is because kubelet on
  the shutdown node is not available to delete the pods so the `StatefulSet`
  cannot create a new pod with the same name. 
- If there are volumes used by the pods, the `VolumeAttachments` will not be
  deleted from the original shutdown node so the volumes used by these pods
  cannot be attached to a new running node. As a result, the application running
  on the `StatefulSet` cannot function properly. 
- If the original shutdown node comes up, the pods will be deleted by kubelet
  and new pods will be created on a different running node. If the original
  shutdown node does not come up, these pods will be stuck in terminating status
  on the shutdown node forever.

## Node Tainting
- To mitigate the above situation, a user can manually add the tain
  `node.kubernetes.io/out-of-service` with either `NoExecute` or `NoSchedule`
  effect to a Node, making it out-of service. If the
  `NodeOutOfServiceVolumeDetach` feature gate is enabled on
  `kube-controller-manager` and a `Node` is marked out-of-service with this
  taint, the pods on the node will be forcefully deleted if there are no
  matching tolerations on it and volume detach operations for the pod
  terminating on the node will happen imediately. This allows the pods on the
  `out-of-service` node to recover quickly on different node.
- During a non-graceful shutdown, pods are terminated in 2 phases:
    - Fore delete the pods that do not have matching `out-of-service`
      tolerations.
    - Immediately perform detach volume operation for such pods.
    - Note:
        - Before adding the taint `node.kubernetes.io/out-of-service`, it should
          be verified that the node is already in shutdown or power off state,
          not in the middle of restarting.
        - The user is required to manually remove the out-of-service taint after
          the pods are moved to a new node and the user has checked that the
          shutdown node has been recovered since the user was the one who
          originally added the taint.

# Swap memory Management

- To enable swap on a node, the `NodeSwap` feature gate must be enabled on the
  kubelet, and the `--fail-swap-on` command line flag or `failSwapOn`
  configuration setting must be set to false.
- When the memory swap feature is turned on, k8s data such as the content of
  secret objects that were written to tmpfs now could be swapped to disk.
- A user can also optionally configure `memorySwap.swapBehavior` in order to
  specify how a node will use swap memory. For example:

```yaml
memorySwap:
    swapBehavior: UnlimitedSwap
```

- `UnlimitedSwap` (default) - k8s workloads can use as much swap memory as they
  request, up to the system limit.
- `LimitedSwap` - The utilization of swap memory by k9s workloads is subjects to
  limitations. Only pods of burstable qos are permitted to employ swap.
- IF configuration for `memorySwap` is not specified and the feature gate is
  enabled, by default the kubelet will apply the same behaviour as the
  `UnlimitedSwap` setting.
- With `LimitedSwap`, pods that do not fall udner the bursable qos classificatin
  (i.e `BestEffort / Guaranteed` qos pods) are prohibited from utlizing swap
  memory. To maintain the aforementioned security and node health guarantees,
  these pods are not permitted to use swap memory when `LimitedSwap` is in
  effect.

Prior to detailing the calculation of the swap limit, it is necessary to define
the following terms:
- `noteTotalMemory` - The total amount of physical memory available on the node.
- `totalPodsSwapAvailable` - The total amount of swap memory on the node that is
  available for use by pods. Some swap memory may be reserved for system use.
- `containerMemoryRequest` - The container's memory request.


Swap limitation is configured as:
```bash
(containerMemoryRequest / nodeTotallMemory) * totalPodsSwapAvailable
```

- For container within burstable qos pods, it is possible to opt-out of swap
  usage by specifying memory requests that are equal to memory limits. 
- Container configured in this manner will not have access to swap memory.
- Swap is supported only with cgroup v1, cgroup v1 is not supported.









