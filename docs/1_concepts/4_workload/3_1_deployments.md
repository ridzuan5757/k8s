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




