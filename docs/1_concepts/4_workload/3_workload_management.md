# Workload Management

k8s provides several built-in APIs for declarative management of the workloads
and the components of those workloads.

Ultimately, the applications running as containers inside pods. However,
managing individual pods would requires lot of effort. For example, if a pod
fails, we probably want to run a new pod to replace it via k8s.

We use the k8s API to create the workload object that represents a higher
abstraction level than a pod, and then the k8s control plane automatically
manages pod objects on our behalf, based on the specification for the workload
object that has been define using the manifest YAML file.

The built-in APIs for managing workloads are:

**Deployment** (and, indirectly **ReplicaSet**), the most common way to run
application in the cluster. Deployment is a good fit for managing stateless
application workload for the cluster, where any pod in the deployment is
interchangeable and can be replaced if needed. (Deployment are a replacement for
the legacy ReplicationController API).

A **StatefulSet** lets us manage one or more pods - all running the same
application code - where the pod rely in having a distinct identity. This is
different from Deployment where the pods are expected to be interchangeable. The
most common use for a StatefulSet is to be able to make a link between its pods
and their persistent storage. For example, we can run StatefulSet that
associates each pods with a PersistentVolume. If one of the pods in the
StatefulSet fails, k8s makes a replacement pod that is connected to the same
PersistentVolume.

A **DaemonSet** defines pods that provides facility that are local to a specific
node. For example, a driver that lets containers on that node access a storage
system. We use DaemonSet when the driver, or other node-level service, has to
run on the node where it is useful. Each pod in a DaemonSet performs a role
similar to a system daemon on a classic Unix / POSIX server. A DaemonSet might
be fundamental to the operation of the cluster, such as a plugin to let that
code access cluster networking, it might help us to manage the node, or it could
provide less essential facilities that enhance the container platform we are
running. We can run DaemonSets and their pods across every node in the cluster,
or  acess just a subset (for example, only install the GPU accelerator driver on
the node that have GPU installed).

We can use a **Job** or a **CronJob** to define tasks that run to completion and
then stop. A Job reprenets a one-off task, where each CronJob repeats according
to a schedule.
