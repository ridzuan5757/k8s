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
