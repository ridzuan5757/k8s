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

## Writing a DaemonSet Spec

### Create a DaemonSet

The file below describes a DaemonSet that runs the `fluentd`-`elasticsearch`
Docker image:

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
    name: fluentd-elasticsearch
    namespace: kube-system
    labels:
        k8s-app: fluentd-logging
spec:
    selector:
        matchLabels:
            name: fluentd-elasticsearch
    template:
        metadata:
            labels:
                name: fluentd-elasticsearch
        tolerations:
        # these tolerations are to have the daemonset runnable on control plane
        # remove them if the control plane should not run Pods
        - key: node-role.kubernetes.io/control-plane
          operator: Exists
          effect: NoSchedule
        - key: node-role.kubernetes.io/master
          operator: Exists
          effect: NoSchedule
        containers:
        - name: fluentd-elasticsearch
          image: quay.io/fluentd_elasticsearch/fluentd:v2.5.2
          resources:
            limits:
                memory: 200Mi
            requests:
                cpu: 100m
                memory: 200Mi
          volumeMounts:
          - name: varlog
            mountPath: /var/log
        # it may be desirable to set a high priority class to ensure that a
        # DaemonSet Pod preempts running Pods
        priorityClassName: important
        terminationGracePeriodSeconds: 30
        volumes:
        - name: varlog
          hostPath:
            path: /var/log
```

### Required Fields

As with all other k8s config, a DaemonSet needs the following fields:

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
spec:
```

The name of a DaemonSet object must be a valid DNS subdomain name.

### Pod Template

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
spec:
    template:
```

`.spec.template` is a pod template. It has exactly the same schema as a Pod,
except it is nested and does not contain the following fields:

```yaml
apiVersion:
kind:
```

In addition to required fields for a Pod, a Pod template in a DaemonSet has to
specify appropriate labels. **A Pod Template in a DaemonSet must have a
`RestartPolicy` equal to `Always`, or be unspecified, which defaults to
`Always`.

### Pod Selector

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
spec:
    template:
    selector:
```

`.spec.selector` field is a pod selector. It works the same as the
`.spec.selector` of a Job.

We must specify a pod selector that matches the labels of the `.spec.template`.
Also once a DaemonSet is created, its `.spec.selector` can not be mutated.
Mutating the pod selector can lead to the unintentional oprphaning of Pods, and
it was found to be confusing to users.

The `.spec.selector` is an object consisting of two fields:
- `matchLabels` - works the same as `.spec.selector` of a
  ReplicationController.
- `matchExpresstion` - allows to build more sophisticated selectors by
  specifying key, list of values and an operator that relates the key and
  values.

When the two are specified, the result is AND-ed. The `.spec.selector` must
match the `.spec.template.metadata.labels`. Config with these two not matching
will be rejected by the API.

### Running Pods on select Nodes

If we specify the following field:

```yaml
spec:
    template:
        spec:
            nodeSelector:
```

The DaemonSet controller will create Pods on nodes which match that node
selector. Likewise if we specify the following field:

```yaml
spec:
    template:
        spec:
            affinity:
```

The DaemonSet controller will create Pods on nodes that match that node
affinity. If none are specified, then the DaemonSet controller will create Pods
on all nodes.

## DaemonPod Schedule

A DaemonSet can be used to ensure that all elligible nodes run a copy of a Pod.
The DaemonSet controller creates a Pod for each eligible node and adds the
`.spec.affinity.nodeAffinity` field of the Pod to match the target host.

After a Pod is created, the default scheduler typically takes over and then
binds the Pod to the garget host by setting the `.spec.nodeName` field. If the
new Pod cannot fit on the node, the default scheduler may preempt (evict) some
of the existing Pods based on the priority of the new Pod.

> [!NOTE]
> If it is important that the DaemonSet pod run on each node, it is often
> desirable to set the `.spec.template.spec.priorityClassName` of the DaemonSet
> to a PriorityClass with a higher priority to ensure that this eviction occurs.
