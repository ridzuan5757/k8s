# Deployments

A Deployment provides declarative updates for Pods and ReplicaSets.

We describe a desired state in a Deployment, and the Deployment Controller
changes the actual state to the desired state at a controlled rate. We can
defined Deployments to create new ReplicaSets, or to remove existing Deployments
and adopt all their resources with new Deployments.

> [!NOTE]
> Do not manage ReplicaSets owned by a Deployment.

## Use Case
- Create a Deployment to rollout ReplicaSet. The ReplicaSet creates Pods in the 
background. After the manifest file has been applied, check the status of the 
rollout to see if it succeeds or not.
- Declare the new state of the Pods. This is done by updating the `PodTemplateSpec` 
of the Deployment. A new ReplicaSet is created and the Deployment manages moving 
the Pods from old ReplicaSet to the new one at a controlled rate. Each new 
ReplicaSet updates the revision of the Deployment.
- Rollback to an earlier Deployment revision. This can be implemented if the 
current state of the Deployment is not stable. Each rollback updates the revision 
of the Deployment.
- Scale up the Deployment to facilitate more load.
- Pause the rollout of a Deployment to apply multiple fixes to tis
  PodTemplateSpec and then resume to start a new rollout.
- Use the status of the Deployment as an indicator that a rollout has stuck.
- Clean up older ReplicaSets that is not needed anymore.

## Creating a Deployment

The following is an example of a Deployment. It creates a ReplicaSet to bring up
three `nginx` pods:

```yaml
# nginx-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

In this exampple:
- A Deployment named `nginx-deployment` is created, indicated by the
  `.metadata.name` field. This name will become basis for the ReplicaSets and
  Pods which are created later.
- The Deployment creates a ReplicaSet that creates three replicated Pods,
  indicated by the `.spec.replicas` field.
- The `.spec.selector` field defines how the created ReplicaSet finds which pods
  to manage. In this case, we select a label that is defined in the Pod
  template (`app: nginx`). However, more sophisticated selection rules are
  possible, as long as the Pod template itself satisfies the rule.

  > [!NOTE]
  > THE `.spec.selector.matchLabels` field is a map of key-value pairs. A single
  > key-value in the `matchLabels` map is equivalent to an element of
  > `matchExpressions`, whose `key` field is "key", the `operator` is "in", and
  > the `values` array contains only "valus". All of the requirements, from both
  > `matchLabels` and `matchExpressions`, must be satisfied in order to match.

- The `template` field contains the following sub-fields:
    - The pods are labeled `app: nginx` using the `.metadata.labels` field.
    - The pod template's specification, or `.template.spec` field, indicates
      that the pods run one container, `nginx`, which runs the `nginx` Docker
      Hub image at version 1.14.2.
    - Create one container and name it `nginx` using the
      `.spec.template.spec.containers[0].name` field.

Make sure k8s cluster is up and running before the Deployment can be started.
- Create the Deployment by running the following command:
  
  ```bash
  kubectl apply -f nginx-deployment.yaml
  ```

- Run `kubectl get deployments` to check if the Deployment was created. If the
  Deployment is still being created, the output is similar to the following:

  ```bash
  NAME               READY   UP-TO-DATE   AVAILABLE   AGE
  nginx-deployment   0/3     0            0           1s
  ```
  
  When we inspect the Deployments in the cluster, the following fields are
  displayed:
  - `NAME` lists the names of the Deployments in the namespace.
  - `READY` displays how many replicas of the application are available to the
    users. It follows the pattern ready/desired.
  - `UP-TO-DATE` displays the number of replicas that have been updated to the
    desired state.
  - `AVAILABLE` displays how many replicas of the application are available to
    the users.
  - `AGE` displays the amount of time that the application has been running.

  Notice how the number of desired replicas is 3 according to `.spec.replicas`
  field.
- To see the Deployment rollout status, run `kubectl rollout status
  deployment/nginx-deployment`.

  The output is similar to:

  ```bash
  Waiting for rollout to finish: 2 out of 3 new replicas have been updated...
  deployment "nginx-deployment" successfully rolled out
  ```

  Notice that the DEployment has created all three replicas, and all replicas
  are up-to-date (they contain the latest Pod template) and available.

- To see the ReplicaSet `rs` created by the Deployment, run `bubectl get rs`.
  The output is similar to this:

  ```bash
  NAME                          DESIRED   CURRENT   READY   AGE
  nginx-deployment-75675f5897   3         3         3       18s
  ```

  ReplicaSet output shows the following fields:
  - `NAME` lists the names of the ReplicaSet in the namespace.
  - `DESIRED` displays the desired number of replicas of the application, which
    was defined in the Deployment template file. This is the desired state.
  - `CURRENT` displays how many replicas are currently running.
  - `READY` displays how many replicas of the application are available to the
    users.
  - `AGE` displays the amount of time that the application has been running.

  Notice that the name of the ReplicaSet is always formatted as
  `[DEPLOYMENT-NAME]-[HASH]`. This name will become bases for the Pods which are
  created. The `HASH` string is same as the `pod-template-hash` label on the
  ReplicaSet.
- To see the labels automatically generated for each Pod, run `kubectl get pods
  --show-labels`. The output is similar to:

  ```bash
  NAME                                READY     STATUS    RESTARTS   AGE       LABELS
  nginx-deployment-75675f5897-7ci7o   1/1       Running   0          18s       app=nginx,pod-template-hash=75675f5897
  nginx-deployment-75675f5897-kzszj   1/1       Running   0          18s       app=nginx,pod-template-hash=75675f5897
  nginx-deployment-75675f5897-qqcnn   1/1       Running   0          18s       app=nginx,pod-template-hash=75675f5897
  ```

  The created ReplicaSet ensures that there are three `nginx` Pods.

> [!NOTE]
> Appropriate selector and Pod template labels must be specified in a
> Deployment (ub tgus case, `app: nginx`).
> 
> Do not overlap labels or selectors with other controllers (including other
> Deployments and StatefulSets). k8s will not prevent overlapping, and if
> multiple controllers have overlapping selectors, those controllers might
> conflict and behave unexpectedly.

## pod-template-hash label

> [!CAUTION]
> Do not change this label.

The `pod-template-hash` label is added by the Deployment controller to every
ReplicaSet that has a Deployment creates or adopts.

This label ensures that child ReplicaSets of a Deployment do not overlap. It is
generated by hashing the PodTemplate of the ReplicaSet and using the resulting
hash as the label value that is added to the ReplicaSet selector, Pod template
labels, and in any existing Pods that the ReplicaSet might have.

## Updating a Deployment

> [!NOTE]
> A Deployment's rollout is triggered if and only if the Deployment's Pod
> template(that is, `.spec.template`) is changed, for example if the labels or
> container images of the template are updated. Other updates, such as scaling
> the Deployment, does not trigger a rollout.

One of the way to update the deployment is by directly updating the container
image version, for example from tag 1.14.2 to 1.16.1:

```bash
kubectl set image deployment.v1.apps/nginx-deployment nginx=nginx:1.16.1
```

The following command is also valid:

```bash
kubectl set image deployment/nginx-deployment nginx=nginx:1.16.1
```

where `deployment/nginx-deployment` indicates the Deployment, `nginx` indicates
the Container the update will take place and `nginx:1.16.1` indicates the new
image and its tag. The output is similar to:

```bash
deployment.apps/nginx-deployment image updated
```

Alternatively, we can also edit the Deployment template file and change
`.spc.tempalte.spec.containers[0].image` from `nginx:1.14.2` to `nginx:1.16.1`:

```bash
kubectl edit deployment/nginx-deployment
```

The output is similar to:

```bash
deployment.apps/nginx-deployment edited
```

To see the rollout status, run:

```bash
kubectl rollout status deployment/nginx-deployment
```

The output is similar to this:

```bash
Waiting for rollout to finish: 2 out of 3 new replicas have been updated...
```

Or this:

```bash
deployment "nginx-deployment" successfully rolled out
```

In order to get more details on the updated Deployment:
- After the rollout succeeds, we can view the Deployment by running `kubectl get
  deployments`. The output is similar to:

  ```bash
  NAME               READY   UP-TO-DATE   AVAILABLE   AGE
  nginx-deployment   3/3     3            3           36s
  ```
- Run `kubectl get rs` to see the Deployment updated the Pods by creating a new
  ReplicaSet and scaling it up to 3 replicas, as well as scaling down the old
  ReplicaSet to 0 replicas.

  ```bash
  kubect get rs
  ```

  The output is similar to:

  ```bash
  NAME                          DESIRED   CURRENT   READY   AGE
  nginx-deployment-1564180365   3         3         3       6s
  nginx-deployment-2035384211   0         0         0       36s
  ```
- Running `kubectl get pods` should now only show the new Pods:

  ```bash
  NAME                                READY     STATUS    RESTARTS   AGE
  nginx-deployment-1564180365-khku8   1/1       Running   0          14s
  nginx-deployment-1564180365-nacti   1/1       Running   0          14s
  nginx-deployment-1564180365-z9gth   1/1       Running   0          14s
  ```
  Next time we want to update these Pods, we only need to update the
  Deployment's Pod template again.

  Deployment ensures that only certain number of Pods are down while they are
  being updated. By default, it ensures that at least 75% of the desired number
  of Pods are up.

  For example, if we look at above Deployment closely, we will see that it first
  create a new Pod, then delete an old Pod, and creates another new one. It does
  not kill old Pods until a sufficient number of new Pods have come up, and does
  not create new Pods until a sufficient number of old Pods have been killed. It
  make sure that at least 3 Pods are available and that at max 4 Pods in total
  are available. In case of a Deployment with 4 replicas, the number of Pods
  would be between 3 and 5.
- Use `kubectl describe deployments` to get details of the Deployments. The
  output is similar to this:

  ```bash
  Name:                   nginx-deployment
  Namespace:              default
  CreationTimestamp:      Thu, 30 Nov 2017 10:56:25 +0000
  Labels:                 app=nginx
  Annotations:            deployment.kubernetes.io/revision=2
  Selector:               app=nginx
  Replicas:               3 desired | 3 updated | 3 total | 3 available | 0 unavailable
  StrategyType:           RollingUpdate
  MinReadySeconds:        0
  RollingUpdateStrategy:  25% max unavailable, 25% max surge
  ```
  We can see that when we first create the Deployment, it created a ReplicaSet
  and scaled it up to 3 replicas directly. When we updated the deployment, it
  created a new ReplicaSet and scaled it up to 1 and waited for it to come up.

  Then it scaled down the old ReplicaSet to 2 and scaled up the new ReplicaSet
  to 2 so that at least 3 Pods were available and at most 4 Pods were created at
  all times.

  It then continued scaling up and down the new and the old ReplicaSet, with the
  same rolling update strategy. Finally we will have 3 available replicas in the
  new ReplicaSet, and the old ReplicaSet is caled down to 0.

> [!NOTE]
> k8s does not count terminating Pods when calculating the number of
> `availableReplicas`, which must be between `replicas - maxUnavailable` and
> `replicas + maxSurge`.
> 
> As a result, we might notice that there are more Pods than expected during a
> rollout, and that the total resources consumed by the Deployment is more than
> `replicas + maxSurge` until the `terminationGracePeriodSeconds` of the
> terminating Pods expires.

### Rollover (aka multiple updates in-flight)

Each time a new Deployment is observed by the Deployment controller, a
ReplicaSet is created to bring up the desired Pods. If the Deployment is
updated, the existing ReplicaSet that controls Pods whose labels match
`match.selector` but whose template does not match `.spec.template` are scaled
down. Eventully, the new ReplicaSet is scaled to `spec.Replicas` and all old
ReplicaSets is caled to 0.

If we update a Deploymeent while an existing rollout is in progress, the
Deployment creates a new ReplicaSet as per update and start scaling that up, and
rolls over the ReplicaSet that is was scaling up previously -- it will add it to
its list of old ReplicaSets and start scaling it down.

For example, suppose that we create a Deployment to create 5 replicas of
`nginx:1.14.2`, but then update the Deployment to create 5 replicas of
`nginx:1.16.1`, when only 3 replicas of `nginx:1.14.2` has been created. In
that case, the Deployment immediately starts killing the 3 `nginx:1.14.2` Pods
that it had created, and starts creating `nginx:1.16.1` Pods. It does not wait
for the 5 replicas of `nginx:1.14.2` to be created before changing course.

### Label selector updates

It is generally discouraged to make label selector updates and it is suggested
to plan the selectors upfront. In any case, if we need to perform a label
selector update, exercise a great caution and understand the implication that
came along with this action.

> [!NOTE]
> In API version `apps/v1`, a Deployment's label selector is immutable after it
> gets created.

- Selector additions require the Pod template labels in the Deployment spec to
  be updated with the new label too, otherwise a validation error is returned.
  This change is a non-overlapping one, meaning that the new selector does not
  select ReplicaSets and Pods created with the old selector, resulting in
  orphaning all old ReplicaSets and creating a new ReplicaSet.
- Selector updates change the existing value in a selector key -- result in the
  same behaviour as additions.
- Selector removals removes an existing key from the Deployment selector -- do
  not require any changes in the Pod template labels. Existing ReplicaSets are
  not orphaned, and a new ReplicaSet is not created, but note that the removed
  label still exists in any existing Pods and ReplicaSets.

## Rolling Back a Deployment

Sometimes, a Deployment rollback might be needed; for example, when the
Deployment is not stabe, such as crash looping. By default, all of the
Deployment's rollout history is kept in the system so that we can rollback
anytime we want.

> [!NOTE]
> A Deployment's revision is created when a Deployment's rollout is triggered.
> This means that the new revision is created if and only if the Deployment's Pod
> template `.spec.template` is changed. For example, if we update the labels or
> container image of the template. 
>
> Other updates, such as scaling the Deployment, do not create a Deployment
> revision, so that we can facilitate simultaneos manual or auto-scaling. This
> means that when we roll back to an earlier revision, only the `template` part
> is being rolled back.

- Suppose that we made a typo while updating the Deployment, by putting the
  image name as `nginx:1.161` instead of `nginx:1.16.1`:

  ```bash
  kubectl set image deployment/nginx-deployment nginx=nginx:1.161
  ```

  The output is similar to this:

  ```bash
  deployment.apps/nginx-deployment image updated
  ```
- This will cause the rollout to get stuck, which can be checked via:
  
  ```bash
  kubectl rollout status deployment/nginx-deployment
  ```

  The output is similar to this:

  ```bash
  Waiting for rollout to finish: 1 out of 3 new replicas have been updated...
  ```
- We can also observe the deployment status by checking the ReplicaSet via `kubectl
  get rs`. The ouput will be similar to this:

  ```bash
  NAME                          DESIRED   CURRENT   READY   AGE
  nginx-deployment-1564180365   3         3         3       25s
  nginx-deployment-2035384211   0         0         0       36s
  nginx-deployment-3066724191   1         1         0       6s
  ```
  Observe that there is one deployment where the ready state is not matching
  with the desired state.
- Looking back at the Pods created, we will also see that 1 Pod is created by
  new ReplicaSet is stuch in an image pull loop. This can be checked using
  `kubectl get pods`. The output is similar to this:

  ```bash
  NAME                                READY     STATUS             RESTARTS   AGE
  nginx-deployment-1564180365-70iae   1/1       Running            0          25s
  nginx-deployment-1564180365-jbqqo   1/1       Running            0          25s
  nginx-deployment-1564180365-hysrc   1/1       Running            0          25s
  nginx-deployment-3066724191-08mng   0/1       ImagePullBackOff   0          6s
  ```

  > [!NOTE]
  > The Deployment controller stops the bad rollout automatically, and stops
  > scaling up the new ReplicaSet. Ths depends on the rolling update parameters
  > `maxUnavailable` that we specified. k8s by default sets this value to 25%.

  To fix this, we need to rollback to a previous revision of Deployment that is
  stable.

### Checking Rollout History of a Deployment

We can check the revision of the deployment using the following command:

```bash
kubectl rollout history deployment/nginx-deployment
```

The output is similar to this:

```bash
deployments "nginx-deployment"
REVISION    CHANGE-CAUSE
1           kubectl apply --filename=https://k8s.io/examples/controllers/nginx-deployment.yaml
2           kubectl set image deployment/nginx-deployment nginx=nginx:1.16.1
3           kubectl set image deployment/nginx-deployment nginx=nginx:1.161
```

`CHANGE-CAUSE` is copied from the Deployment annotation
`kubernetes.io/change-cause` to its revision upon creation. We can specify the
`CHANGE-CAUSE` message by:
- Annotating the Deployment with `kubectl annotate deployment/nginx-deployment
  kubernetes.io/change-cause="image updated to 1.16.1"`
- Manually editing the manifest of the resource.

To see the details of each revision, run:

```bash
kubectl rollout history deployment/nginx-deployment --revision=2
```

The output is similar to this:

```bash
deployments "nginx-deployment" revision 2
  Labels:       app=nginx
          pod-template-hash=1159050644
  Annotations:  kubernetes.io/change-cause=kubectl set image deployment/nginx-deployment nginx=nginx:1.16.1
  Containers:
   nginx:
    Image:      nginx:1.16.1
    Port:       80/TCP
     QoS Tier:
        cpu:      BestEffort
        memory:   BestEffort
    Environment Variables:      <none>
  No volumes.
```

### Rolling Back to a Previous Revision

Now that we have decided to undo the current rollout and rollback to the
previous revision:

```bash
kubectl rollout undo deployment/nginx-deployment
```

The output is similar to this:

```bash
deployment.apps/nginx-deployment rolled back
```

Alternatively, we can rollback to specific revision by specifying it with
`--to-revision`:

```bash
kubectl rollout undo deployment/nginx-deployment --to-revision=2
```

The output is similar to this:

```bash
deployment.apps/nginx-deployment rolled back
```

The Deployment is now rolled back to a previous stable revision. As we can see,
a `DeploymentRollback` event for rolling back to revision 2 is generated from
Deployment Controller. We can check if the rollback is successful and the
Deployment is running as expected:

```bash
kubectl get deployment nginx-deployment
```

The output is similar to this:

```bash
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   3/3     3            3           30m
```

Checking the description of the Deployment using `kubectl describe deployment
nginx-deployment` will output something similar to this:

```bash
Name:                   nginx-deployment
Namespace:              default
CreationTimestamp:      Sun, 02 Sep 2018 18:17:55 -0500
Labels:                 app=nginx
Annotations:            deployment.kubernetes.io/revision=4
                        kubernetes.io/change-cause=kubectl set image deployment/nginx-deployment nginx=nginx:1.16.1
Selector:               app=nginx
Replicas:               3 desired | 3 updated | 3 total | 3 available | 0 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25% max unavailable, 25% max surge
Pod Template:
  Labels:  app=nginx
  Containers:
   nginx:
    Image:        nginx:1.16.1
    Port:         80/TCP
    Host Port:    0/TCP
    Environment:  <none>
    Mounts:       <none>
  Volumes:        <none>
Conditions:
  Type           Status  Reason
  ----           ------  ------
  Available      True    MinimumReplicasAvailable
  Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  <none>
NewReplicaSet:   nginx-deployment-c4747d96c (3/3 replicas created)
Events:
  Type    Reason              Age   From                   Message
  ----    ------              ----  ----                   -------
  Normal  ScalingReplicaSet   12m   deployment-controller  Scaled up replica set nginx-deployment-75675f5897 to 3
  Normal  ScalingReplicaSet   11m   deployment-controller  Scaled up replica set nginx-deployment-c4747d96c to 1
  Normal  ScalingReplicaSet   11m   deployment-controller  Scaled down replica set nginx-deployment-75675f5897 to 2
  Normal  ScalingReplicaSet   11m   deployment-controller  Scaled up replica set nginx-deployment-c4747d96c to 2
  Normal  ScalingReplicaSet   11m   deployment-controller  Scaled down replica set nginx-deployment-75675f5897 to 1
  Normal  ScalingReplicaSet   11m   deployment-controller  Scaled up replica set nginx-deployment-c4747d96c to 3
  Normal  ScalingReplicaSet   11m   deployment-controller  Scaled down replica set nginx-deployment-75675f5897 to 0
  Normal  ScalingReplicaSet   11m   deployment-controller  Scaled up replica set nginx-deployment-595696685f to 1
  Normal  DeploymentRollback  15s   deployment-controller  Rolled back deployment "nginx-deployment" to revision 2
  Normal  ScalingReplicaSet   15s   deployment-controller  Scaled down replica set nginx-deployment-595696685f to 0
```

## Scaling a Deployment

We can scale a Deployment using the `scale` subcommand:

```bash
kubectl scale deployment/nginx-deployment --replicas=10
```

The output is similar to this:

```bash
deployment.apps/nginx-deployment scaled
```

Assuming that horizontal pod autoscaling is enabled in the cluster, we can set
up an autoscaler for the Deployment and choose the minimum and maximum number of
Pods we need to run based on the CPU utilization of the existing Pods.

```bash
kubectl autoscale deployment/nginx-deployment --min=10 --max=15 --cpu-percent=80
```

The output is similar to this:

```bash
deployment.apps/nginx-deployment scaled
```

### Proportional scaling

RollingUpdate Deployments support running multiple versions of an application at
the same time. When we or an autoscaller scales a RollingUpdate Deployment that
is in the middle of a rollout (either in progress or paused), the Deployment
controller balances the additional replicas in the existing active ReplicaSets
in order to motigate risk. This is called proportional scaling.

For example, we are running a Deployment with 10 replicas, maxSurge = 3, and
maxUnavailable = 2.
- Ensure that the 10 replicas in the Deployment are running.
  ```bash
  kubectl get deploy
  ```

  The output is similar to this:
  ```bash
  NAME                 DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
  nginx-deployment     10        10        10           10          50s
  ```
- We update to a new image which happens to be unresolvable form inside the
  cluster.
  ```bash
  kubectl set image deployment/nginx-deployment nginx=nginx:sometag
  ```

  The output si similar to this:
  ```bash
  deployment.apps/nginx-deployment image updated
  ```
- The image update starts a new rollout with new ReplicaSet, but it is blocked
  ue to the `maxUnavailable` requirement that was mentioned above. This can be
  checked using the `rs` subcommand:

  ```bash
  kubectl get res
  ```

  The output is similar to this:

  ```bash
  NAME                          DESIRED   CURRENT   READY     AGE
  nginx-deployment-1989198191   5         5         0         9s
  nginx-deployment-618515232    8         8         8         1m
  ```
- Then a new scaling request for the Deployment comes along. The autoscaler
  increments the Deployment replicas to 15. The Deployment controller needs to
  deide where to add these 5 new replicas. If we were not using proportional
  scaling, all 5 of them would be added in the new ReplicaSet.

  With proportional scaling, we spread the additional replicas across all
  ReplicaSets. Bigger propertions go to the ReplicaSets with the most replicas
  and lower proportions to go to ReplicaSets with less replicas. Any leftovers
  are added to the ReplicaSet with most replicas. ReplicaSets with zero replicas
  are not scaled up.

In the example above, 3 replicas are added to the old ReplicaSet and 2 replicas
are added to the new ReplicaSet. The rollout process should eventuallly move all
replicas to the new ReplicaSet, assuming the new replicas become healthy. We can
confirm this using `kubectl get deploy` which should returns:

```bash
NAME                 DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment     15        18        7            8           7m
```

The rollout status checked using `kubectl get rs` confirms how the replicas were
added to each ReplicaSet:

```bash
NAME                          DESIRED   CURRENT   READY     AGE
nginx-deployment-1989198191   7         7         0         7m
nginx-deployment-618515232    11        11        11        7m
```

## Pausing and Resuming a rollout of a Deployment

When we update a Deployment, or plan to, we can pause rollouts for that
Deployment before we trigger one or more updates. When we are ready to apply
those changes, we resume rollouts for the Deployment. This approach allows us to
apply multiple fixes in between pausing and resuming without triggering
unecessary rollouts.

For example, with a Deployment that was created, we can obtain the deployment
details using `kubectl get deploy`. Suppose that the outcome of the subcommand
is something like this:

```bash
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
nginx     3         3         3            3           1m
```

We can check the rollout status using `kubectl get rs` subcommand which will
output something similar to this:

```bash
NAME               DESIRED   CURRENT   READY     AGE
nginx-2142116321   3         3         3         1m
```

We can pause the deployment using the following command:

```bash
kubectl rollout pause deployment/nginx-deployment
```

The output is similar to this:

```bash
deployment.apps/nginx-deployment paused
```

Then we can perform different sutffs, say we want to update the image of the
Deployment, we can do something like this:

```bash
kubectl set image deployment/nginx-deployment nginx=nginx:1.16.1
```

The output is similar to this:

```bash
deployment.apps/nginx-deployment image updated
```

Notice that no new rollout started, which can be verified using the following
command:

```bash
kubectl rollout history deployment/nginx-deployment
```

The output is similar to this:

```bash
deployments "nginx"
REVISION  CHANGE-CAUSE
1   <none>
```

We can check the rollout status using `kubectl get rs` to verify that the
existing ReplicaSet has not changed. The output will be something similar to
this:

```bash
NAME               DESIRED   CURRENT   READY     AGE
nginx-2142116321   3         3         3         2m
```

We can make as many updates as we wish, for example update the resources that
will be used:

```bash
kubectl set resources deployment/nginx-deployment \
    -c=nginx --limits=cpu=200m,memory=512Mi
```

The output is similar to this:

```bash
deployment.apps/nginx-deployment resource requirements updated
```

The initial state of the Deployment prior to pausing its rollout will continue
its function, but new updates to the Deployment will not have any effect as long
as the Deployment rollout is paused.

Eventually, we will resume the Deployment rollout and obser a new ReplicaSet is
coming up with all the new updates:

```bash
kubectl rollout resume deployment/nginx-deployment
```

The output is similar to this:

```bash
kubectl rollout resume deployment/nginx-deployment
```

The output is similar to this:

```bash
deployment.apps/nginx-deployment resumed
```

Watch the status of the rollout until it is done:

```bash
kubectl get rs -w
```

The output is similar to this:

```bash
NAME               DESIRED   CURRENT   READY     AGE
nginx-2142116321   2         2         2         2m
nginx-3926361531   2         2         0         6s
nginx-3926361531   2         2         1         18s
nginx-2142116321   1         2         2         2m
nginx-2142116321   1         2         2         2m
nginx-3926361531   3         2         1         18s
nginx-3926361531   3         2         1         18s
nginx-2142116321   1         1         1         2m
nginx-3926361531   3         3         1         18s
nginx-3926361531   3         3         2         19s
nginx-2142116321   0         1         1         2m
nginx-2142116321   0         1         1         2m
nginx-2142116321   0         0         0         2m
nginx-3926361531   3         3         3         20s
```

Get the status of the latest rollout using `kubectl get rs`. The output is
similar to this:

```bash
NAME               DESIRED   CURRENT   READY     AGE

nginx-2142116321   0         0         0         2m
nginx-3926361531   3         3         3         28s
```

> [!NOTE]
> We cannot rollback a paused Deployment until we resume it.

## Deployment status

A Deployment enters various states during its lifecycle. It can be progressing
while rolling out a new ReplicaSet, it can be completes, or it can fail to
progress.

### Progressing Deployment

k8s marks a Deployment as **progressing** when one of the following tasks is
performed:
- The Deployment creates a new ReplicaSet.
- The Deployment is scaling up its newest ReplicaSet.
- The Deployment is scaling down its older ReplicaSets.
- New Pods become ready or available (ready for at least `MinReadySeconds`).

When the rollout becomes `progressing`, the Deployment controller adds a
condition with the following attributes to the Deployment's
`.status.conditions`:

```json
{
    type: Progressing,
    status: "True",
    reason: NewReplicaSetCreated | 
    reason: FoundNewReplicaSet | 
    reason: ReplicaSetUpdated
}
```

We can monitor the progress for a Deployment by using `kubectl rollout status`.
