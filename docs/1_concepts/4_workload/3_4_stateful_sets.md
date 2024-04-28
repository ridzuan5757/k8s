# StatefulSets

StatefulSet is the workload API object used to manage stateful applications.
This API manages the deployment and scaling of a set of Pods, and provides
guarantees about the ordering and uniqueness of these Pods.

Like a Deployment, a StatefulSet manages Pods that are based on an identical
container spec. Unlike a Deployment, a StatefulSet maintains a sticky identity
for each of its Pods. These pods are created from the same spec, but are **not
interchangeable**. Each has persistent identifier that it maintains across any
rescheduling.

If we want to use storage volumes to provide persistence for the workload, we
can use StatefulSet as part of the solution. Although individual Pods in a
StatefulSet a susceptible to failure, the persistent Pod identifiers make it
easier to match existing volumes to the new Pods that replace any that have
failed.

## Usage

StatefulSets are valuable for applications that require one or more of the
following:
- Stable, unique network identifiers.
- Stable, persistent storage.
- Ordered, graceful deployment and scaling.
- Ordered, automated rolling updates.

In the above, stable is synonymous with persistence across Pods rescheduling. If
an application does not require any stable identifiers or ordered deployment,
deletion, or scaling, we should deploy the application using a workload object
that provides a set of stateless replicas. Deplyoment or ReplicaSet may be
better suited to the stateless needs.

## Limitations
- The storage for a given Pod must either be provisioned by a
  `PersistenceVolumeProvisioner` based on the requested storage class, or
  pre-provisioned by an admin.
- Deleting and / or scaling a StatefulSet down will not delete the volumes
  associated with StatefulSet. This is done to ensure data safety, which is
  generally more valuable than an automatic purge of all related StatefulSet
  resources.
- StatefulSets currently require a Headless Service to be responsible for the
  network identity of the Pods. It is up to our own responsibility for creating
  this service.
- StatefulSets do not provide any guarantees on the termination of Pods when
  StatefulSet is deleted. To achieve ordered and graceful termination of the
  pods in the StatefulSet, it is possible to scale the StatefulSet down to 0
  prior to deletion.
- When using RollingUpdates with the default Pod Management Policy
  `OrderedReady`, it is possible to get into a broken state that require manual
  intervention to repair.


