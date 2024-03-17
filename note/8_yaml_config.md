# YAML Config

k8s resources are primarily configured using YAML files. We have used the
`kubectl edit` command to edit resources in the cluster on-demand. We can get a
copy of the deployment file by piping it via the terminal.

```bash
kubectl get deployment synergychat-web -o yaml > web-deployment.yaml
```

There are 5 top-level fields in the file:
- `apiVersion: apps/v1` - Specifies the version of the k8s API that we are using
  to create the object.
- `kind: Deployment` - Specifies the type of object that we are configuring.
- `metadata` - Metadata about the deployment, such as when it is created, name
  and its ID.
- `spec` - The desired state of the deployment. Most impactful edits, like how
  many replicas we want, will be made here.
- `status` - The current state of the deployment. We would not edit this
  directly. It is just for us to see what is going on with the deployment.

We can make changes to the `deployment` state using the YAML file using the
following command:

```bash
kubectl apply -f web-deployment.yaml
```

We should get a warning that lets we know that we are missing
`last-applied-configuration` annotation. This is fine since we got the warning
because we created this deployment the quick and disty way, by using `kubectl
create deployment` instead of creating a YAML file and using `kubectl apply -f`.

However, because we have now updated it with `kubectl apply`, the annotation is
now here and we would not get the warning again.
