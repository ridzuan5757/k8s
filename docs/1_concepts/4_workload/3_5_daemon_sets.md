# DaemonSet

A DaemonSet ensures that all or some Nodes run a copy of a Pod. As nodes are
added to the cluster, Pods are added to them. s nodes are removed from the
cluster, those Pods are garbage collected. Deleting a DAemonSet will clean up
the Pods it created.

Some typical uses of a DaemonSet are:
- Running a cluster of storage daemon on every node.
- Running a logs collection dameon on every node.
- Running a node monitoring daemon on every node.

In a simple case, one DaemonSet, covering all nodes, would be used for each type
of daemon. A more complex setup might use multiple DaemonSets for a single type
of daemon, but with different flags and or different memory and cpu requests for
different hardware types.


