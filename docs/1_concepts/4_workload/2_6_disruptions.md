# Disruption

Pods do not disappear until someone be a person or controller destroys them, or
there is unavoidable hardware or system software error.

## Involuntary disruption

We call these unavoidable cases involuntary disruption to an application. For
examples:
- A hardware failure of the physical machine backing the node.
- Cluster administrator deletes VM instance by mistake.
- Cloud provider or hypervisor failure making VM disappear.
- Kernel panic.
- The node disappear from the cluster due to cluster network partition.
- Eviction of a pod due to the node being out-of-resources.

Except for the out-of-resources condition, all these conditions should be
familiar for most users as they are not k8s specific.

## Voluntary disruption

These cinlude both actions initiated by the application owner and those
initiated by a cluster administrator. Typical application owner actions include:
- Deleting the deployment or other controller that manages the pod.
- Updating a deployment's pod template causing a restart.
- Directly deleting a pod.

Cluster administrator actions include:
- Draining a node for repair or upgrade.
- Draining a node from a cluster to scale the cluster down.
- Removing a pod from a node to permit soemthing else to fit on that node.


