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

## Pod identity

StatefulSet Pods have a unique identity that consists of an ordinal, a stable
network identity, and stable storage. The identity sticks to the Pod, regardless
of which node it is rescheduled on.

### Ordinal Index

For a StatefulSet with N replicas, each Pod in the StatefulSet will be assigned
an integer ordinal, that is unqique over the Set. By default, Pods will be
assigned orginals from 0 up through N-1. The StatefulSet controller will also
add a pod label with this index `apps.kubernetes.io/pod-index`.

### Start Ordinal

`.spec.ordinals` is an optional field that allows us to configure the integer
ordinals assigned to each Pod. It defaults to nil. We must enable the
`StatefulSetStartOrdinal` feature gate to use this field. Once enabled, we can
configure the following options:
- `.spec.ordinals.start` : If the `.spec.ordinals.start` field is set, Pods will
  be assigned ordinals from `.spec.ordinals.start` up through
  `.spec.ordinals.start + .spec.replicas - 1`.

### Stable Network ID

Each Pod in a StatefulSet derives its hostname from the name of the StatefulSet
and the ordinal of the Pod. The pattern for the contructed hostname is
`$(statefulset name)-$(ordinal)`. The example above will create three Pods 
named `web-0, web-1, web-2`. 

A StatefulSet can use a Headless Service to control the domain of its Pods. The
domain is managed by this Service takes the form of 
`$(service-name).$(namespace).svc.cluster.local`, where `cluster.local` is the
cluster domain. As each Pod is created, it gets a matching DNS subdomain, taking
the form: `$(podname).$(governing-service-domain)`, where the governing service
is defined by the `serviceName` field on the StatefulSet.

Depending on how DNS is configured in the cluster, we may not be able to look up
the DNS name for a newly-run Pod immediately. This behaviour can occur when
other clients in the cluster have already sent queries for the hostname of the
Pod before it was created. 

Negative caching (normal in DNS) means that the results of previous failed
lookups are remembered and reused, even after the Pod is running, for at least a
few seconds.

If we need to discover Pods promptly after they acre created, we can:
- Query the k8s API directly such as using `watch` rather than relying on DNS
  lookup.
- Decreate the time of caching in k8s DNS provider (typically this means editing
  the ConfigMap for CoreDNS, which currently caches for 30 seconds).

As mentioned in the limitations section, we are responsible for creating the
Headless Service responsible for the network identity of the Pods.

Here a some examples of choices for Cluster Domain, Service name, StatefulSet
name, and how that affects the DNS names for the StatefulSet's Pods.

|**Cluster Domain**|cluster.local|cluster.local|kube.local|
|---|---|---|---|
|Service|default/nginx|foo/nginx|foo/nginx|
|StatefulSet|default/web|foo/web|foo/web|
|StatefulSet Domain|nginx.default.svc.cluster.local|nginx.foo.svc.cluster.local|nginx.foo.svc.kube.local|
|Pod DNS|web-{0..N-1}.nginx.default.svc.cluster.local|web-{0..N-1}.nginx.foo.svc.cluster.local|web-{0..N-1}.nginx.foo.svc.kube.local|
|Pod Hostname|web-{0..N-1}|web-{0..N-1}|web-{0..N-1}|

> [!NOTE]
> Cluster domain will be set to `cluster.local` unless otherwise configured.

### Stable Storage

For each VolumeClaimTemplate entry defined in a StatefulSet, each Pod receives
one PersistentVolumeClaim. In the nginx example above, each Pod receives a
single PersistentVolume with a StorageClass of `my-storage-class` and 1GB of
provisioned storage. If no StorageClass is specified, then the default
StorageClass will be used.

When a Pod is rescheduled into a node, its `volumeMounts` mount the
PersistentVolumes associated with its PersistentVolumeClaims. Note that, the
PersistentVolumes associated with the Pod's PersistentVolumeClaims are not
deleted when the Pods, or StatefulSet are deleted. This must be done manually.

### Pod Name Label

When the StatefulSet controller creates a Pod, it adds a label,
`statefulset.kubernetes.io/pod-name`, that is set to the name of the Pod. this
label allows us to attach a Service to a specific Pod in the StatefulSet.

### Pod Index Label

When the StatefulSet controller creates a Pod, the new Pod is labelled with
`apps.kubernetes.io/pod-index`. The value of this label is the ordinal index of
the Pod. this label allows us to route traffic to a particular pod index, filter
logs/metrics using the pod index label, and more. The feature gate
`PodIndexLabel` must be enabled for this feature. It is enabled by default.


