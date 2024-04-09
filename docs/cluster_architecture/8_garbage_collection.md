# Garbage Collection

Garbage collection GC is a collective term for the various mechanisms k8s uses
to clean up cluster resources. This allows the clean up of resources:
- Terminated pods
- Completed jobs
- Object without owner references
- Unused containers and container images
- Dynamically provisioned `PersistentVolumes` with a `StorageClass` reclaim
  policy of delete.
- Stale or expired certificate signing request
- Nodes deleted in the following scenarios:
    - On a cloud when the cluster uses cloud controller manager
    - On premises when the cluster uses addon similar to a cloud controller
      manager
- Node lease Object

## Owner and dependents

Many object in k8s link to each other through owner references. Owner references
tell the control plane which objects are dependent on each others. k8s uses
owner references to give the control plane, and other API clients, the
opportunity to clean up related resources before deleting an object. In most
cases, k8s manages owner references automatically.

Ownership is different from the labels and selectors mechanism that some
resources also use. For example, consider a service that creates `EndpointSlice`
objects. The Service uses labels to allow the control plane to determine which
`EndpointSlice` object are used for that Service.

In addition to the labels, each `EndpointSlice` that is managed on behalf of a
service has an owner reference. Owner references help different parts of k8s
avoid interfering with objects they do not control.

- Cross-namespace owner references are disallowed by design. Namespaced
  dependents can specify cluster-scoped or namespaced owners. A namespaced owner
  must exist in the same namespace as the dependent. If it does not, the owner
  reference is treated as absent, and the dependent is subject to deletion once
  all owners are verified absent.
- Clster scoped dependents can only specify cluster-scoped owners. In v1.20+, if
  a cluster-scoped dependent specifies a namesaced kind as an owner, it is
  treated as having unresolvable owner referencem and is not able to be garbage
  collected.
- In v1.20+, if the garbage collector detects an invalid cross-namespace
  `ownerReference` or a cluster-scped dependent with an `ownerReference`
  referencing a namespaced ind, a warning event with a reason of
  `OwnerRefInvalidNamespace` and an `involvedObject` of the invalid dependent is
  reported. We can check for that kind of event by running `kubectl get events
  -A --field-selector=reason=OwnerRefInvalidNamespace`.

## Cascading deletiion

k8s checks for and deltes object that no longer have owner references, like the
pods left behind when we delete `ReplicaSet`. When we delete an object, we can
control whether k8s deletes the object's dependents automatically, in a process
called **cascading deletion**. There are 2 types of cascading deletion:
- Foreground cascading deletion
- Background cascading deletion

We can also control how and when garbage collection deletes resources that have
owner references using k8s finalizers.

#### k8s Finalizers
Namespaced key that tells k8s to wait until specific conditions are met before
it fully deletes an object marked for deletion.

### Foreground cascading deletion

The owner object we are deleting first enters deletion in progress state. In
this state, the following happens to the owner object:
- The k8s API server sets the object's `metadata.deletionTimestamp` field to the
  time the object is marked for deletion.
- The k8s API server also sets the `metadata.finalzers` field to
  `foregroundDeletion`.
- The object remains visble throught the k8s API until the deletion process is
  complete.

After the owner object enters the deletion progress state, the controller
deletes the dependent. After deleting all the dependent objects, the controller
deletes the owner object. At this point, object is no longer visible in k8s API.

During foreground cascading deletion, the only dependents that block owner
deletion are those that have the `ownerReference.blockOwnerDeletion=true` field.

### Background cascading deletion

In background cascading deletion, the k8s API server deletes the owner object
immediately and the controller cleans up the dependent objects n the background.
By default, k8s uses background cascading deletion unless we manually use
foreground deletion or choose to orphan the dependent objects.

### Orphaned dependents

When k8s deltes an owner object, the dependents left behind are called orphan
objects. By default, k8s deltes dependent objects however this could be
overrided.

## Garbage collection of unused containers and images

The kubelet performs gc on unused images every 2 minutes and on unused
containers on every minute. We should avoid using external gc tools, as these
can break the kubelet behaviour and remove containers that should exist.

To configure options for unused container and image garbage collection, tune the
kubelet using a configuration file and change the parameters related to garbage
collection using the `KubeletConfiguration` resource type.

### Container image lifecycle

k8s manages the lifecycle of all images through its image manager, which is
part of the kubelet, wth the cooperation of cadvisor. The kubelet consisders the
following disk usage limits when making gc decisions:
- `HighThresholdPercent`
- `LowThresholdPercent`

Disk usage above the configured `HighThresholdPercent` value triggers gc, which
deletes images in order based on the last time they were used, starting with the
oldest first. The kubelet deletes images until disk usage reaches the 
`LowThresholdPercent` value.

### Garbage collection for unused container images

We can specify the maximum time a ocal image can be unused for, regardless of
disk usage. This is a kubelet setting that we configure for each node.

To configure the setting, enable the `ImageMaximumGCAge` feature gate for the
kubelet, and also set a value for the `ImageMaximumGCAge` field in the kubelet
cofiguration file.

The value is specified as a k8s duration; for example, we can set the
configuration to `3d12h` which means 3 days and 12 hours.

### Container garbage collection

The kubelet gc collects unused containers based on the following variables that
can be defined:
- `MinAge` - minimum age at which the kubelet can gc a container. Disable by
  setting it to `0`.
- `MaxPerPodContainer` - the maximum number of dead containers each pod can
  have. Disable by setting to less than `0`.
- `MaxContainers` - the maximum number od dead containers the cluster can have.
  Disable by setting to less than `0`.

`MaxPerPodContainer` and `MaxContainers` may potentially conflict with each
other in situations where retaining the maximum number of containers per pod
(`MaxPerPodContainer`) would go outside allowable total of global dead
containers (`MaxContainers`). 

In this situation, the kubelet adjusts `MaxPerPodContainer` to address the
conflict. A worst-case scenatio would be to downgrade `MaxPerPodContainer` to `1` 
and evict the oldest containers. Additionally, containers owned by pods that
have been deleted are removed once they are older than `MinAge`.

The kubelet only fc the containers it manages.

## Configuring garbage collection

We can tune garbage collection of resources by configuring options specific to
controllers managic those resources.
