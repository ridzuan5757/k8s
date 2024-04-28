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

## Components

The example below demonstrates the components of a StatefulSet.

```yaml
# service.yaml

apiVersion: v1
kind: Service
metadata:
    name: nginx
    labels:
        app: nginx
spec:
    ports:
    - port: 80
      name: web
    clusterIP: None
    selector:
        app: nginx
```

```yaml
# statefulset.yaml

apiVersion: apps/v1
kind: StatefulSet
metadata:
    name: web
spec:
    selector:
        matchLabels:
            app: nginx
    serviceName: "nginx"
    replicas: 3
    minReadySeconds: 10
    template:
        metadata:
            labels:
                app: nginx
        spec:
            terminationGracePeriodSeconds: 10
            containers:
            - name: nginx
              image: registry.k8s.io/nginx-slim:0.8
              ports:
              - containerPort: 80
                name: web
              volumeMounts:
              - name: www
                mountPath: /usr/share/nginx/html
    volumeClaimTemplates:
    - metadata:
        name: www
      spec:
        accessModes: [ "ReadWriteOnce" ]
        storageClassName: "my-storage-class"
        resources:
            requests:
                storage: 1Gi
```

> [!NOTE]
> This example uses the `ReadWriteOnce` access mode, for simplicity. For
> production use, the k8s project recommonds using the `ReadWriteOncePod` access
> mode instead.

In the example above:
- A headless service, named `nginx`, is used to control the network domain.
- The StatefulSet, named `web`, has a Spec that indicates that 3 replicas of the
  nginx container will be launched in unique Pods.
- The `volumeClaimTemplates` will provide stable storage using PersistentVolumes
  provisioned by a PersistentVolume Provisioner.

The name of a StatefulSet object must be a valid DNS label.

### Pod Selector

We must set the `.spec.selector` field of a StatefulSet to match the labels of
its `.spec.template.metadata.labels`. Failing to specify a matching Pod Selector
will result in a validation error during StatefulSet creation.

### Volume Claim Templates

We can set the `.spec.volumeClaimTemplates` field to create a
PersistentVolumeClaim. this will provide stable storage to the StatefulSet if
either:
- The StorageClass specified for the volume claim is set up to use dynamic
  provisioning, or
- The cluster already contains a PersistentVolume with the correct StorageClass
  and sufficient available storage space.

### Minimum ready seconds

`.spec.minReadySeconds` is an optional field that specifies the minimum number
of seconds for which a newly created Pod should be running and ready without any
of its containers crashing, for it to be considered available.

This is used to check progression of a rollout when using a RollingUpdate
strategy. The field defaults to 0 (the Pod will be considered available as soon
as it ready).




