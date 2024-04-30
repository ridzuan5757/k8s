# Leases

Distributed system often have a need for **leases**, which provide a mechanism
to lock shared resources and coordinate activity between members of a set. In
k8s, the lease concept is represented by `Leases` object in the
`coordination.k8s.io` API group, which are used for system-critical capabilities
such as node heartbeats and component-level leader election.

## Node heartbeats

k8s uses the leases API to communicate kubelet node heartbeats to the k8s API
server. For every `Node`, there is `Lease` object with a maching name in the
`kube-node-lease` namespace. 

Under the hood, every kubelet heartbeat is an **update** request to this `Lease`
object, updating the `spec.renewTime` field for the Lease. The k8s control plane
uses the time stamp of this field to determine the availability of this node.

## Leader election

k8s also uses leases to ensure only one instance of a component is running at a
given thim. This is used by control plane components like
`kube-controller-manager` and `kube-scheduler` in high availability
configurations, where only one instance of the component should be actively
running while the other instances are on standby.

## API server identity

`kube-apiserver` uses the Lease API to publish its identity to the rest of the
system. While not particularly useful on its own, this provides a machanism for
clients to discover how many instaces of `kube-apiserver` are operating the k8s
control plane. Existence of `kube-apiserver` leases enables future capabilities
that may require coordination between each `kube-apiserver`.

We can inspect leases owned by each `kube-apiserver` by checking for lease
object in the `kube-system` namespace with the name
`kube-apiserver-sha2560hash`. Alternatively, we can use the label selector
`apiserver.kubernetes.io/identity=kube-apiserver`:

```bash
kubectl -n kube-system get lease -l apiserver.kubernetes.io/identity=kube-apiserver
```

```bash
NAME                                        HOLDER                                                                           AGE
apiserver-07a5ea9b9b072c4a5f3d1c3702        apiserver-07a5ea9b9b072c4a5f3d1c3702_0c8914f7-0f35-440e-8676-7844977d3a05        5m33s
apiserver-7be9e061c59d368b3ddaf1376e        apiserver-7be9e061c59d368b3ddaf1376e_84f2a85d-37c1-4b14-b6b9-603e62e4896f        4m23s
apiserver-1dfef752bcb36637d2763d1868        apiserver-1dfef752bcb36637d2763d1868_c5ffa286-8a9a-45d4-91e7-61118ed58d2e        4m43s
```

```yaml
apiVersion: coordination.k8s.io/v1
kind: Lease
metadata:
  creationTimestamp: "2023-07-02T13:16:48Z"
  labels:
    apiserver.kubernetes.io/identity: kube-apiserver
    kubernetes.io/hostname: master-1
  name: apiserver-07a5ea9b9b072c4a5f3d1c3702
  namespace: kube-system
  resourceVersion: "334899"
  uid: 90870ab5-1ba9-4523-b215-e4d4e662acb1
spec:
  holderIdentity: apiserver-07a5ea9b9b072c4a5f3d1c3702_0c8914f7-0f35-440e-8676-7844977d3a05
  leaseDurationSeconds: 3600
  renewTime: "2023-07-04T21:58:48.065888Z"
```

Expired leases from kube-apiservers that no longer exist are garbage collected
by new kube-apiserver after 1 hour. We can disable API server identity leases by
disabling `APIServerIdentity` feature gate.

## Workloads

Our won workload can define its own use of leases. For example, we might run a
custom controller where a primary or leader member performs operations that its
peers do not. We define a ease so that the controller replicas can select or
elect a leader, using the k8s API for coordination.

If we do use a lease, it is a good practice to define a name for the lease that
is obviously linked to the product or component. For example, if we have a
component named Example Foo, use a lease named `example-foo`.

If a cluster operator or another end user could deploy multiple instance of a
compoenent, select a name prefix an pick a meachnism such as hash of the name of
the deployment to avoid name colissions for the leases.

We can use another approach as long as it achieves the same outcome - dfferent
software products that do not conflict with one another.
