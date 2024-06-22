# ReplicaSet

A ReplicaSet's purpose is to maintain a stable set of replica Pods running at
any given time. As such, it is often used to guarantee the availability of a
specified number of identical Pods.

## Mechanism

A ReplicaSet is defined with fields, including a selector that specifies how to
identify Pods it can acquire, a number of replicas indicating how many Pods it
should be maintaining, and a pod template specifying the data of new Pods it
should create to meet the number of replicas criteria.

A ReplicaSet then fulfills its purpose by creating and deleting Pods as needed
to reach the desired number. When a ReplicaSet needs to create new Pods, it uses
its Pod template.

A ReplicaSet is linked to its Pods via the Pods' `metadata.ownerReferences`
field, which specifies what resource the current object is owned by. All Pods
acquired by a ReplicaSet have their owning ReplicaSet's identifying information
withinb their ownerReferences field. It is through this link that the ReplicaSet
knows the state of the Pods it is maintaining and plans accordingly.

A ReplicaSet identifies new Pods to acquire by using its selector. If there is a
Pod that has no OwnerReference or the OwnerReference is not a Controller and it
matches a ReplicaSet's selector, it will be immediately acquired by said
ReplicaSet.

## When to use a ReplicaSet

A ReplicaSet ensures that a specified number of pod replicas are running at any
given time. However, a Deployment is a higher-level concept that manages
ReplicaSets and provides declarative updates to Pods along with a lot of other
useful features. 

Therefore, **we should use Deployments instead of directly using ReplicaSets**,
unless custom orchestration update is required or the Deployment itself does not
require qupdates at all.

This actually means that we may neve need to manipulate ReplicaSet objects - we
will be using a Deployment instead and define the application in the `spec`
section.

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: frontend
  labels:
    app: guestbook
    tier: frontend
spec:
  # modify replicas according to your case
  replicas: 3
  selector:
    matchLabels:
      tier: frontend
  template:
    metadata:
      labels:
        tier: frontend
    spec:
      containers:
      - name: php-redis
        image: us-docker.pkg.dev/google-samples/containers/gke/gb-frontend:v5
```

Saving this manifest into `frontend.yaml` and submitting it to k8s cluster will
create the defined ReplicaSet and the Pods that it manages.

```bash
kubectl apply -f frontend.yaml
```

We can then get the current ReplicaSet deployed via `kubectl get rs`. The front
end created will be shown like this:

```bash
NAME       DESIRED   CURRENT   READY   AGE
frontend   3         3         3       6s
```

We can also check the state of the ReplicaSet via `kubectl describe rs/frontend`
and we will see output similar to:

```bash
Name:         frontend
Namespace:    default
Selector:     tier=frontend
Labels:       app=guestbook
              tier=frontend
Annotations:  <none>
Replicas:     3 current / 3 desired
Pods Status:  3 Running / 0 Waiting / 0 Succeeded / 0 Failed
Pod Template:
  Labels:  tier=frontend
  Containers:
   php-redis:
    Image:        us-docker.pkg.dev/google-samples/containers/gke/gb-frontend:v5
    Port:         <none>
    Host Port:    <none>
    Environment:  <none>
    Mounts:       <none>
  Volumes:        <none>
Events:
  Type    Reason            Age   From                   Message
  ----    ------            ----  ----                   -------
  Normal  SuccessfulCreate  13s   replicaset-controller  Created pod: frontend-gbgfx
  Normal  SuccessfulCreate  13s   replicaset-controller  Created pod: frontend-rwz57
  Normal  SuccessfulCreate  13s   replicaset-controller  Created pod: frontend-wkl7w
```

We can also check the Pods being brought up from this manifest file via `kubectl
get pods`. The information shown will be similar to:

```bash
NAME             READY   STATUS    RESTARTS   AGE
frontend-gbgfx   1/1     Running   0          10m
frontend-rwz57   1/1     Running   0          10m
frontend-wkl7w   1/1     Running   0          10m
```

We can also veirfy that the owner referene of these pods is set to the frontend
ReplicaSet. To do thism, get the yaml of one of the Pods running using `kubectl
get pods frontend-gbgfx -o yaml`. The output will show that that the
ReplicaSet's infor set in the metadata's ownerReferences field:

```bash
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: "2024-02-28T22:30:44Z"
  generateName: frontend-
  labels:
    tier: frontend
  name: frontend-gbgfx
  namespace: default
  ownerReferences:
  - apiVersion: apps/v1
    blockOwnerDeletion: true
    controller: true
    kind: ReplicaSet
    name: frontend
    uid: e129deca-f864-481b-bb16-b27abfd92292
...
```

## Non-Template Pod acquisitions

While Pods created using Pod manifest file, it is strongly recommended to make
sure that bare Pods **do not have labels which can match the selector** of one
of the ReplicaSets. The reason for this is because a ReplicaSet is not limit to
owning Pods specified by its template - it can acquire other Pods in the manner
specified in the previous sections.

For example, consider the following ReplicaSet:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod1
  labels:
    tier: frontend
spec:
  containers:
  - name: hello1
    image: gcr.io/google-samples/hello-app:2.0

---

apiVersion: v1
kind: Pod
metadata:
  name: pod2
  labels:
    tier: frontend
spec:
  containers:
  - name: hello2
    image: gcr.io/google-samples/hello-app:1.0
```

As those Pods do not have a Controller or any object as their owner reference
and match the selector of the frontend ReplicaSet, they will be immediately be
acquired by it,

Supposed that we create Pods after the frontend ReplicaSet has been deployed and
has set up its initial Pod replicas to fullfill its replica count requirement.

```bash
kubectl apply -f pod.yaml
```

The new Pods will be acquired by the ReplicaSet, and then immediately terminated
as the REplicaSet would be over its desired count. By fetching the Pods using
`kubectl get pods`, the output will show that the new Pods are either already
terminated, or in process of being terminated:

```bash
NAME             READY   STATUS        RESTARTS   AGE
frontend-b2zdv   1/1     Running       0          10m
frontend-vcmts   1/1     Running       0          10m
frontend-wtsmm   1/1     Running       0          10m
pod1             0/1     Terminating   0          1s
pod2             0/1     Terminating   0          1s
```

If we create the Pods first using `kubectl apply -f rs.yaml`, and then create
the ReplicaSet via `kubectl apply -f frontend.yaml`, we shall see that the
ReplicaSet has acquired the Pods and has only created new one according to its
spec until the number of its new Pods and the original matches its desired
count. As fetching the pods via `kubectl get pods`, the following output will be
obtained:

```bash
NAME             READY   STATUS    RESTARTS   AGE
frontend-hmmj2   1/1     Running   0          9s
pod1             1/1     Running   0          36s
pod2             1/1     Running   0          36s
```

In this manner, a ReplicaSet can own a non-homogeneous set of Pods.

## Writing a ReplicaSet manifest

As with all other k8s API obhects, a ReplicaSet needs the `apiVersion`, `kind`
and `metadata` fields. For ReplicaSet, the `kind` is always a ReplicaSet. When
the control plane creates a new Pods for a ReplicaSet, the `.metadata.name` of
the ReplicaSet is part of the basis for naming those Pods.

The name of a ReplicaSet must be a valid DNS subdomain value, but this can
produce unexpected results for the Pod hostnames. For best compatibility, the
name should follow more restrictive rules for a DNS label. A ReplicaSet also
needs a `.spec` section.

### Pod Template

The `.spec.template` is a pod template which is also required to have labels in
place. In our `frontend.yaml` example we had one label: `tier: frontend`. Be
careful not to overlap with the selectors of other controllers, lest they try to
adopt this Pod.

For the template's restart policy field, `.spec.tempalte.spec.restartPolicy`,
the only allowed value is `Always`, which is the default.

### Pod Selector

The `.spec.selector` field is a label selector. As discussed earlier these are
the labels used to identify potential Pods to acquire. In our frontend.yaml
example, the selector was:

```yaml
matchLabels:
    tier: frontend
```

In the ReplicaSet, `.spec.template.metadata.labels` must match `.spec.selector`,
or it will be rejected by the API.

> [!NOTE]
> For 2 ReplicaSets specifying the same `.spec.selector` but different
> `.spec.template.metadata.labels` and `.spec.template.spec` fields, each
> ReplicaSet ignores the Pods created by the other ReplicaSet.

### Replicas

We can specify how many Pods should run concurrently by setting
`.spec.replicas`. The ReplicaSet will create/delete its Pods to match this
number. If we do not specify `.spec.replicas`, then it defaults to 1.

## Working with ReplicaSets

### Deleting a ReplicaSet and its Pods

To delete a REplicaSet and all of its Pods, use `kubectl delete`. The garbage
collector automatically deletes all of the dependent Pods by default.

When using the REST API or the `client-go` library, we must set
propagationPolicy to `Background` or `Foreground` in the `-d` option. For
example:

```bash
kubectl proxy --port=8080
curl -X DELETE  'localhost:8080/apis/apps/v1/namespaces/default/replicasets/frontend' \
  -d '{"kind":"DeleteOptions","apiVersion":"v1","propagationPolicy":"Foreground"}' \
  -H "Content-Type: application/json"
```

### Deleting just a ReplicaSet

We can delete a ReplicaSet without affecting any of its Pods using `kubectl delete` 
with the `--cascade=orphan` option. When using the REST API or the `client-go`
library, we must set `propagationPolicy` to `Orphan`. For example:

```bash
kubectl proxy --port=8080
curl -X DELETE  'localhost:8080/apis/apps/v1/namespaces/default/replicasets/frontend' \
  -d '{"kind":"DeleteOptions","apiVersion":"v1","propagationPolicy":"Orphan"}' \
  -H "Content-Type: application/json"
```

Once the original is deleted, we can create a new ReplicaSet to replace it. As
long as the old and new `.spec.selector` are the same, then the new one will
adopt the old Pods. Howeverm it will not make any effort to make existing Pods
match a new, different pod template. To update Pods to a new spec in a
controlled way, use a Deployment, as ReplicaSets do not support rolling update
directly.

### Isolating Pods from a ReplicaSet

We can remove Pods from a ReplicaSet by changing their labels. This technique
may be used to remove Pods from service for debugging, data recovery, etc. Pods
that are removed in this way will be replaced automatically (assuming that the
number of replicas is not also changed).

### Scaling a ReplicaSet

A ReplicaSet can be easily scaled up or down by simply updating `.spec.replicas`
field. The ReplicaSet controller ensures that a desired number of Pods with a
matching label selector are available and operational.

When scalign down, the REplicaSet controller chooses which pods to delete by
sorting the available pods to prioritize scaling down pods based on the
following general algorithm:
- Pending and unschedulable pods are scaled down first.
- If `controller.kubernetes.io/pod-deletion-cost` annotation is set, then the
  pod with the lower value will come first.
- Pods on nodes with more replicas come before pods on nodes with fewer
  replicas.
- If the Pods' creation times differ, the pod that was created more recently
  comes before the older pod (the creation times are bucketed on an integer log
  scale when the `LogarithmicScaleDown` feature gate is enabled).

If all of the above match, then the selection is random.

### Pod deletion cost

Using the `controller.kubernetes.io/pod-deletion-cost` annotation, users can set
a preference regarding which pods to remove first when downscaling a ReplicaSet.

The annotation should be set on the pod, the range is [-2147483648, 2147483647].
It represents the cost of deleting a pod comapred to other pods belinging to the
same ReplicaSet. Pods with lower deletion cost are preferred to be deleted
before pods with higher deletion cost.

The implicit value for this annotation for pods that don't set is 0. Negative
values are permitted. Invalid value will be rejected by the API server.

This is a beta feature and it is enabled by default. We can however disableit
using the feature gate `PodDeletionCost` in both kube-apiserver and
kube-controller-manager.

> [!NOTE]
> - This is honored on a best-effort basis, so it does not offer any guarantees
>   on pod deletion order.
> - Users should avoid updating the annotation frequently, such as updating it
>   based on a metric value, because doing so will generate a significant number
>   of pod updates on the api server.

### Example Use Case

The different pods of an application could have different utilization levels. On
scale down, the application may prefer to remove the Pods with lower
utilization. To avoid frequently updating the Pods, the application should
update `controller.kubernetes.io/pod-deletion-cost` once before issuing a scale
down (setting the annotation to a value proportional to pod utilization level).
This works if the application itself controls the down scaling; for example, the
driver pod of a Spark deployment.

### ReplicaSet as a Horizontal Pod Autoscaler Targer

A ReplicaSet can also be a target for Horizontal Pod Autoscalers (HPA). That is,
a ReplicaSet can be auto-scaled by an HPA. Here is an example HPA targeting the
ReplicaSEt we created in the previous example.

```yaml
# hpa.yaml

apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
    name: frontend-scaler
spec:
    scaleTargetRef:
        kind: ReplicaSet
        name: frontend
    minReplicas: 3
    maxReplicas: 10
    targetCPUUtilizationPercentage: 50
```

Submitting this manifest to a k8s cluster should create the defined HPA that
autoscales the target ReplicaSet depending on the CPU usage of the replicated
Pods.

```bash
kubectl apply -f ./hpa.yaml
```

Alternatively, we can use `kubectl autoscale` command to accomplish the same
result.

```bash
kubectl autoscale rs frontend --max=10 --min=3 --cpu-percent=50
```

## Alternatives to ReplicaSet

### Deployment (recommended)

Deployment is an object which can own ReplicaSets and update them and their Pods
via declarative, servver-side rolling updates. While ReplicaSets can be used
independently, today they are mainly used by Deployments as a machanism to
orchestrate Pod creation, deletion and updates.

When we use Deployments we do not have to worry about managing the ReplicaSets
that they create. Deployments own and manage their ReplicaSets. As such, it is
recommended to use Deployments when we want ReplicaSets.

### Bare Pods

Unlike the case where a user directly created Pods, a ReplicaSet replaces Pods
that are deleted or terminated for any reason, such as in the case of node
failure or disruptive node maintenance, such as a kernel upgrade.

For this reason, it is reommended that ReplicaSet is used even if the
application requires only a single Pod. Think of it similarly to a process
supervisor, only it supervises multiple Pods across multiple nodes instead of
individual processes on a single node. A ReplicaSet delegates local container
restarts to some agent on the node such as Kubelet.

### Job

We can use `Job` instead of ReplicaSet for Pods that are expected to terminate
on their own.

### DaemonSet

`DaemonSet` can be used instead of ReplicaSet for Pods that provide a
machine-level function, such as machine monitoring and machine logging. These
Pods have a lifetime that is tied to a machine lifetime: the Pod needs to be
running on the machine before other Pods start, and are safe to terminate when
the machine is otherwise ready to be rebooted/shutdown.

### ReplicationController

ReplicaSets are the successor to ReplicationControllers. The two serve the same
purpose, and behave similarly, except that a ReplicationController does not
support set-based selector requirements as described in the user guide. As such,
ReplicaSets are preferred over ReplicationControllers.
