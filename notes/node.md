# Node
This components run on every node, maintaining running pods and providing the
k8s runtime environment.

## `kubelet`
- Agent that runs on each node in the cluster. Ensures that containers are running
in a pod.
- Takes a set of `PodSpecs` that are provided through yaml file and ensures that
  the containers described in those `PodSpecs` are running and healthy.
- Does not manage containers that were not created by k8s.

## `kube-proxy`
- Network proxy that runs on each node in the cluster, implementing k8s service.
- Maintains network rules on nodes. Allow network communictation to the pods
  from network sessions inside or outside of the cluster.

## Container runtime
- Manage execution and lifecycle of containers within k8s environment.
