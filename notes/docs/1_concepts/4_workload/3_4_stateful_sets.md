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

## Deployment and Scaling Guarantees

- For a StatefulSet with N replicas, when Pods are being deployed, they are
  created subsequently, in order from 0 to N-1.
- When Pods are being deleted, they are terminated in reverse order, from N-1 to
  0.
- Before a scaling operation is applied to a Pod, all of its predecessors must
  be Running and Ready.
- Before a Pod is terminated, all of its successors must be completely shutdown.

The StatefulSet should not specify a `pod.Spec.TerminationGracePeriodSeconds` of
0. This practice is unsafe and strongly discouraged.

When the nginx example above is created, three Pods will be deployed in the
order, web-0, web-1, and web-2. web-1 will not be deployed before web-0 is
Running and Ready, and web-2 will not be deployed until web-1 is Running and
Ready. If web-0 should fail, after web-1 is Running and Ready, but before web-2
is launched, web-2 will not be launched until web-0 is successfully relaunced
and becomes Running and Ready.

If a user were to scale the deployed example by patching the StatefulSet such
that `replicas=1`, web-2 would be terminated first. web-1 would not be
terminated until web-2 is fully shutdown and deleted. If web-0 were to fail
after web-2 has been terminated and is compeltely shutdown, but prior to web-1's
termination, web-1 would not be terminated until web-0 is Running and Ready.

### Pod Management Policies

StatefulSet allows us to relax its ordering guarantees while preserving its
uniqueness and identity via its `.spec.podManagementPolicy` field.

### OrderedReady Pod Management

`OrderedReady` pod management is the default for StatefulSets. It implements the
behaviour described above.

### Parallel Pod Management

`Parallel` pod management tells the StatefulSet controller to launch or
terminate all Pods in parallel, and to not wait for Pods to become Running and
Ready or completely terminated prior to launching or terminating another Pod.
This option only affects the behaviour for scaling operations. Updates are not
affected.

## Update strategies

A StatefulSet's `.spec.updateStrategy.type` allow us to configure and disable
automated rolling update for containers, labels, resource request/limits, and
annotations for the Pods in a StatefulSet. There are 2 possible values:

### `OnDelete`

The StatefulSet controller will not automatically update the Pods in a
StatefulSet. Users must manually delete Pods to cause the controller to create
new Pods that reflect modifications made to a StatefulSet's `.spec.template`.

### `RollingUpdate`

This update strategy implements automated, rolling updates for the Pods in a
StatefulSet. This is the default update strategy.

## Rolling Updates

When a StatefulSet's `.spec.updateStrategy.type` is set to `RollingUpdate`, the
StatefulSet controller will delete and recreate each Pod in the StatefulSet. It
will proceed in the same order as Pod termination (from the largest ordinal to
the smallest), updating each Pod one at a time.

The k8s control plane waits until an updated Pod is Running and Ready prior to
updating its predecessors. If we have set `.spec.minReadySeconds`, the control
plane additionally waits that amount of time after the Pod turns ready, before
moving on.

### Partitioned rolling updates

The `RollingUpdate` update strategy can be partitinoned, by specifying:

```yaml
spec:
    updateStrategy:
        rollingUpdate:
            partition:
```

If a partition is specified, all Pods with an ordinal that is greater than or
equal to the partition will be updated when the StatefulSet's `.spec.template`
is updated. All pods with an ordinal that is less than partition will not be
updated, and, even if they are deleted, they will be recreated at the previous
version.

If a StatefulSet's `.spec.updateStrategy.rollingUpdate.partition` is greater
than its `.spec.replicas`, updates to its `.spec.template` will not be
propagated to its Pods. In most cases we will not need to use a partition, but
they are useful if we want to stage and update, roll out a canary, or perform a
phased roll out.

### Maximum unavailable Pods

We can control the maximum number of Pods that can be unavailable during an
update by specifying this field:

```yaml
spec:
    updateStrategy:
        rollingUpdate:
            maxUnavailable:
```

The value can be an absolute number or a percentage of desired Pods. Absolute
number is calculated from the percentage value by rounding it up. This field
cannot be 0. The default setting is 1.

This field applies to all Pods in the range of `0` to `replicas - 1`. If there
is any unavailable Pod in the range `0` to `replicas - 1`, it will be counted
towards `maxUnavailable`.

> [!NOTE]
> The `maxUnavailable` field is in Aplha stage and it is honored only by API
> servers that are running with the MaxUnavailableStatefulSet feature gate
> enabled.

### Forced rollback

When using Rolling Updates with the default Pod Management Policy 
(`OrderedReady`) it is possible to get into a broken state that requires manual 
intervention to repair.

If we update the Pod template to a configuraiton that never becomes Running and
Ready (for example, due to a  bad binary or application-level configuration
error), StatefulSet will stop the rollout and wait.

In this state, it is not enough to revert the Pod template to a good
configuraiton. Due to a known issue, StatefulSet will continue to wait for the
broken Pod to become Ready (which never happens) before it will attempt to
revert it back to the working configuration.

After verting the template, we must also delete any Pods that StatefulSet had
already attempted to run with the bad configuration. StatefulSet will then begin
to recreate the Pods using the reverted template.

## PErsistentVolumeClaim retention

The optional `.spec.persistentVolumeClaimRetentionPolicy` field controls if and
how PVCs are deleted during the lifecycle of a StatefulSet. We must enable the
`StatefulSetAutoDeletePVC` feature gate on the API server and the controller
manager to use this field. Once enabled, there are two policies we can configure
for each StatefulSet:

### `whenDeleted`

Configures the volume retention behaviour that applies when the StatefulSet is
deleted.

### `whenScaled`

Configures the volume retention behaviour that applies when the replica count of
the StatefulSet is reduced; for example, when scaling down the set.

For each policy that we can configure, we can set the value to either `Delete`
or `Retain`.

### `Delete`

The PVCs created from the StatefulSet `volumeClaimTemplate` are deleted for each
Pod affected by the policy. With the `whenDeleted` policy all PVCs from the
`volumeClaimTemplate` are deleted after their Pods have been deleted. With the
`whenScaled` policy, only PVCs corresponding to Pod replicas being scaled down
are deleted, after their Pods have been deleted.

### `Retain` (default)

PVCs from the `volumeClaimTemplate` are not affected when their Pods is deleted.
This is the behaviour before this new feature.

Bear in mind that these policies only apply when Pods are being removed due to
StatefulSet being deleted or scaled down. For example, if a Pod associated with
StatefulSet fails due to node failure, and the control plane creates a
replacement Pods, the StatefulSet retains the existing PVC. The existing volume
is unaffected, and the cluster will attach it to the node where the new Pod is
about to launch.

The default for policies is `Retain`, matching the StatefulSet behaviour before
this new feature. Here is an example policy.

```yaml
apiVersion: apps/v1
kind: StatefulSet
spec:
    persistentVolumeClaimRetentionPolicy:
        whenDeleted: Retain
        whenScaled: Delete
```

The StatefulSet controller adds owner references to its PVCs, which are then
deleted by the garbage collector after the Pod is terminated. This enables the
Pod to cleanly unmount all volumes before the PVCs are deleted (and before the
backing PV and volume are deleted, depending on the retain policy). When we set
the `whenDeleted` policy to `Delete`, an owner reference to the StatefulSet
instance is placed on all PVCs associated with the StatefulSet.

The `whenScaled` policy must delete PVCs only when a Pod is scaled down, and not
when a Pod is deleted for another reason. When reconciling, the StatefulSet
controller compares its desired replica count to the actual Pods present on the
cluster. Any StatefulSet Pod whose id greater than the replica count is
condemned and marked for deletion. If the `whenScaled` policy is `Delete`, the
condemned Pods are first set as owners to the associated StatefulSet template
PVCs, before the Pod is deleted. This causes the PVCs to be garbage collected
after only the condemned Pods have terminated.

This means that if the controller crashes and restarts, no Pod will be deleted
before its owner reference has been updated appropriate to the policy. If a
condemned Pod is force-deleted while the controller is down, the owner reference
may or may not have been setup, depending on the when the controller crashed. It
may take several reconcile loops to update the owner references, so some
condemned Pods may have set up owner references and others may not.

For this reasons it is recommended to wait for the controller to come back up,
which will verify owner references before terminating Pods. If that is not
possible, the operator should verify the owner references on PVCs to ensure the
expected objects are deleted when Pods are force-deleted.

## Replicas

`.spec.replicas` is an optional field that specifies the number of desired Pods.
It defaults to 1.

Should we manually scale a deployment, example via `kubectl scale statefulset
statefulset --replicas=X`, and then we update that StatefulSet based on a
manifest (for example, by running `kubectl apply-f statefulset.yaml`), when
applying that manifest overwrites the manual scaling that we previously did.

If a HorizontalPodAutoscaler or any similar API for horizontal scaling is
managing scaling for a Statefulset, do not set `.spec.replicas`. Instead, allow
the k8s control plane to manage the `.spec.replicas` field automatically.
