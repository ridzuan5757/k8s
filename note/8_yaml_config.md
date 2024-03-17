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

# Service

Aside from deploying web service, k8s also capable to deploy API services.

## Creating a deployment configuration.

To write deployment from scratch:
- Create a new YAML configuration file.
- Add the `apiVersion` and `kind` fields.
    - The `apiVersion` is `apps/v1`.
    - The `kind` is `Deployment`.
- Add a `metadata/name` field. This should be the deployment name.
- Add a `metadata/labels/app` field and set it to be the same as the deployment
  name. This will be used to select the pods that this deployment manages.
- Add a `spec/replicas` field and set it to `1`. We can always cale up to more
  pods later.
- Add a `spec/selector/matchLabels/app` field and set it to the deployment name.
  This should be matching with the label that we set of Step 4.
- Add a `spec/template/metadata/labels/app` field and set it to the deployment
  name. Again, this should match the label that we set in step 4. Labels are
  important because they are how k8s knows which pods belong to which
  deployments.
- Add a `spec/template/spec/containers` field. This contains a list of
  containers that will be deployed.
    - Note: A hyphen is how we denote a list item in YAML.
    - Set the `name` of the container.
    - Set the `image` name. This tells the k8s where to download the Docker
      image from.

Take a look at the pods that is currently running. We should able to see pods
for the web service and a pod for the api service. However, we might notice that
the api pod is not in a `ready` state. In fact, it should be stuck in
`CrashLoopBackOff` status.

## Trashing Pods

One of the most common problem we will run into when working with k8s is pods
that keep crashing and restarting. This is called `thrashing` and it is usually
caused by one of a few things:
- The application recently had a bug introduced in the latest image version.
- The application is misconfigured and cannot restart properly.
- A dependency of the application is misconfigured and the application cannot
  start properly.
- The application is trying to use too much memory and being killed by k8s.

#### `CrashLoopBackoff`

When a pod's status is `CrashLoopBackoff`, that means the container is crashing
(the program is exiting with error code `1`). Because k8s is all about building
self-healing systems, it will automatically restart the container. However, each
time it tries to restart the container, if it crashes again, it will wait longer
and longer in between restarts. That is the reason it is called a `backoff`.
