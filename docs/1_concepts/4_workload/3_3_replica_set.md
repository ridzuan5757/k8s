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


