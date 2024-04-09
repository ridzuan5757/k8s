# Control Plane Components
- Make global decisions about the cluster:
    - Scheduling.
    - Detecting and responding to cluster events.
- Can be run on any machine in the cluster.

#### Note
- Setup  scripts typically start all control plane components on the same
  machine.
- Do not run user containers in this machine.

## `kube-apiserver`
- Expose k8s API / front-end of the control plane.
- Designed to scale horizontally / scales by deploying more instances / run
  serveral instances of `kube-apiserver` and balance traffic between those
  instances.

## `etcd`
- Consistent and highly available key value store used as k8s backing store for
  all cluster data.
- Make sure to have back up plane for the data.

## `kube-schedulaer`
- Watches for newly created pods with no assigned node and selects node for them
  to run on.
- Scheduling decision factors:
    - Individual and collective resource requirements.
    - Hardware/software/policy constraints.
    - Affinity and anti-affinity specifications.
    - Data locality.
    - Inter-workload interferences.
    - Deadlines.

## `kube-controller-manager`
- Controller : A control loop that watch shared state of cluster through the api
  server and make changes attempting to move the current state towards the
  desired state.
- Run controller process.
- Each controller is a separate process, but to reduce complexity, they are all
  compiled into a single binary and run as a single process.
- Controller types:
    - Node controller - responsible for noticing and responding when nodes go
      down.
    - Job controller - watches for job object that represent one-off tasks, then
      creates pods to run those tasks to completion.
    - EndpointSlice controller - populate EndpointSlice objects to provide link
      between services and pods.
    - ServiceAccount controller - create default ServiceAccounts for new
      namespaces.


