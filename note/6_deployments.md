# Deployments

A `deployment` provides declarative updates for `pods` and `replicaSets`. We
describe our desired state in a `deployment`, and the deployment controller's
job is to make current state match the desired states.

When we delete a pod, we instantly see that a new pod was created in its place.
That is because the desired state described in `deployment` says we want 2 pods
running at all times. When we delete one, the deployment controller sees that
the current state does not match the desired state, so it creates a new pod to
make them match again.

## Checking deployment

To view the YAML file of the current deployment:

```bash
kubectl get deployment synergychat-web -o yaml
```

To edit the deployment in order to modify the desired state, such as changing
the number of replicas:

```bash
kubectl edit deployment synergychat-web
```

Wait until all of changes made has been completed and then verify the changes:

```bash
kubectl view pod
```

Check the proxy server:

```bash
kubectl proxy
```
