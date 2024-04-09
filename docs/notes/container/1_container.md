# Containers

Each container that we run is repeatable; the standardization from having
dependencies included means that we get the same behaviour wheerever we run it.

Containers decouple applications from the underlying host infrastructure. This
makes deployment easier in different cloud or OS environments.

Each node in a k8s cluster runs the containers that form the pods assigned to
that node. Containers in a pod are co-located and co-scheduled to run on the
same node.

## Container images

A container image is a ready-to-run software package containing everything
needed to run an application: the code and any runtime it requires, application
and system lbiraries, and default values for any essential settings.

Containers are intended to be stateless and immutable:
- We should not change the code of the container that is already running.
- If we have containerized application and want to make changes, the correct
  process is to build a new image that includes the change, then recreate the
  container to start from the updated image.

## Container runtimes

A fundamental compoennt that empower k8s to run containers effectively. It is
responsibe for managing the execution and lifecycle of containers within k8s
environment.

k8s supports container runtimes such as:
- containerd
- CRI-O
- any other implementation of k8s CRI

Usually, we can allow our cluster to pick the default container runtime for a
pod. If we need to use more than one container runtime in the cluster, we can
specify the RuntimeClass for a pod to make sure that k8s run those containers
using a particular container runtime.

We can also use RuntimeClass to run different pods with the same container
runtime but different settings.
