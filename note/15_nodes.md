# Nodes

We have talked about in a production environment, we will have multiple nodes in
our cluster and we have been using single node cluster with minikube. The nice
thing about k8s is that almost everything that we do is abstracted away from the
underlying infrastructure with the `kubectl` CLI.

## Deploying to production

### GKE, EKS, AKS
These are all managed k8s services, offered by cloud providers. They are all
pretty similar. GKE appear to be the most feature-rich out of the thress. GKE
also has auto-pilot mode that makes it so that we do not have to worry about
managing nodes at all.

The nice thing about a managed offering is that it can be configured to handle
autoscaling at the node level. This means that we can set up our cluster
automatically add and remove nodes based on the load on the cluster.

### Manual

We can also set up our own cluster manually. One of the way is by having custom
scripts that configure a cluster on top of standard EC2 instances. Then have our
own autoscaling scripts that add and remove nodes based on the load of the
cluster. It is also pissble to do the same thing on physical machines.

# Node Types

Broadly speaking, there are two types of machines in a production k8s cluster:
- Control Plane
- Worker Nodes

The control plane is responsible for managing the cluster. It is where the API
server, scheduler, and controller manager live. The control plane used to be
called "master nodes", but the term is deprecated now.

When we hear the word "node" used in isolation, it is usually referring to
worker nodes. Worker nodes are machines that actually running the containerized
applications.

# Resource Requests

A resource request is the amount of resource that a pod requests from the node
it is running on. This is different from resource limit, which is the maximum
amount of resource that a pod is allowed to consume before it is throttled or
killed.

## The need for requests

Let say we have 2 nodes:

|Node|RAM|
|---|---|
|Node 1|8GB|
|Node 2|8GB|

And we have 4 pods:

|Pod|Node|RAM|
|---|---|---|
|Pod 1|Node 1|3GB|
|Pod 2|Node 1|3GB|
|Pod 3|Node 2|3GB|
|Pod 4|Node 2|3GB|

This is valid. Currently, only 6 out of 8GB of RAM is being used on each node.
The trouble is, even though each node has 2GB of RAM left, if we try to add
another pod, and it ends up utilizing more 2GB or RAM, it will crash. This is
where resource requests can be used to solve this.

If wwe add a resource request of 3GB, k8s will know that each pod needs 3GB of
RAM to run. If we try to schedule a new pod with the request in place, k8s will
gracefully tell us it does not have enough resources to do so, or it will use a
node in a cluster that has at least 3GB of RAM available. This is implemented
either on deployment or pod kind config.

```yaml
resources:
    limits:
        memory: 4000Mi
    requests:
        memory: 4000Mi
```

After the deploymentm check the pods:

```bash
kubectl get pods
```

Say if this resource requests exceeds the machine available RAM, the pod will
stuck in "pending" state. We can use `describe` to look at the "Events" section
of the output to see the issue.

```bash
kubectl describe pod <pod-name>
```


