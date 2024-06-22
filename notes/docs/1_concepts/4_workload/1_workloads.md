# Workloads

A workload is an application running on k8s. Whether the workload is a single
component or several that work together, on k8s that us run inside a set of
pods. In k8s, a pod represent a set of running containers in the cluster.

k8s pods have defined lifecycle. For example, once a pod is running in the
cluster, then a critical fault on the node where that pod is running means that
all the pods on that node will fail. k8s treats that level of failure as final:
we would need to create a new pod to recover, event if the node later become
healthy.

However, we do not need to manage each pod directly. We can use **workload
resources** that manage a set of pods. These resources configure controllers
that make sure the right number of right kind of pod are running, to match the
desired state.

k8s provides several built-in workload resources:

## Deployment and ReplicaSet
- This replace the legacy resource `ReplicationController`
- Deployment is a good fit for managing stateless application workload on the
  cluster, where any pod in the deployment is interchangeable and can be
  replaced if needed.

## StatefulSet
- We can run one or more related pods that perform state tracking.
- For example: If our wokrload records data persistently, we can run a
  `StatefulSet` that matches each pod with a `PersistentVolume`. The code
  running in the pods for that `StatefulSet` can replicate data to other pods in
  the same `StatefulSet` to improve overall resilience.

## DaemonSet
- Defines pods that provide facilities that are local to nodes.
- Everyt time a node is added to the cluster that match specification in a
  `DaemonSet`, the control plane schedules a pod for that `DaemonSet` onto the
  new node.
- Each pod in a `DaemonSet` performs a job similar to a system daemon on a
  classic Unix / POSIX server.
- A `DaemonSet` might be fundamental to the operation of the cluster, such as
  plugin to run cluster networking, it might help us to manage the node, or it
  could provide optional behaviour that enahances the container platform that is
  currently running.

## Job and CronJob
- Provide different ways to define tasks that run to a completion and then stop.
  We can use `Job` to define a task that runs to completion, just once.
- `CronJob` can be used to run the same `Job` multiple times according a
  schedule.

> We can also use third-party workload resource that provide additiona
> behaviours.
>
> Using a custom resource definition, we can add in a third-party workload
> resource if we want a specific behaviour that is not part of k8s core.
>
> For example, if we want to run a group of pods for the application but stop
> work unless all the pods are available, (perhaps for some high-throughput
> distributed task), then we can implement or install an extension that does
> provide that feature.
