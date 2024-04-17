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


