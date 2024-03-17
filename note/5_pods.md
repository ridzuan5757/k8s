# Pods

A `pod` is the smallest and simplest unit in k8s object model that we create or
deploy. It represents one (or sometimes more) running container(s) in a cluster.
In a simple web application, we might have one single pod: the web server.

As the traffic grows, we might deploy the same code to multiple pods to handle
the increased load. Serveral pods, one codebase. In much mode complex backend
system, we might have serveral pods for the web server and serveral pods that
handle video processing: multiple pods, multiple codebases.

**pods are just wrappers around containers**. We can think of it as a Docker
container k8s black magic wrapped on top of it. The container is the actual
application, and the pod is the k8s abstraction that manages the container and
the resources it needs to run.

### Adding pods replica

```bash
kubectl edit deployment <deployment_name>
```

This will open the deployment configuration yaml file. Under the `spec` section,
we can modify the `replicas` field in order to adjust numer of replicas running
for that deployment. Verify the change by running:

```bash
kubectl get pods
```

## Ephemeral

Pods die, they die often, and sometimes without warning. The temporary nature of
pods is one of the defining feature of k8s. Unlike traditional virtual machines
or barebone servers that might run indefinitely until hardware failure, pods are
designed to be spun up, torn down and restarted at moment's notice.
- The ephemerality of pods provides flexibility and resilience. If a pod
  encounters a problem, it can be easily terminated and replaced with a new
  healthy instance. This model not only allows for high availability but also
  promotes immutability. Instead of manuall patching or updating existing
  environments, we can simply just spin up new versions of the entire
  environment.
- As a developer, it is crucial to understand that it is rarely a good idea to
  store persistent data on pod. They can be terminated and replaced, and any
  locally saved data will be lost. The image inside the pod should be capable to
  be restarted from scratch often.

### Deleting pods

Get a list of running pods:

```bash
kubectl get pods
```

Print the logs of the older pod:

```bash
kubectl logs <pod_name>
```

Kill the older pod:

```bash
kubectl delete pod <pod_name>
```

Verify the change:

```bash
kubectl get pods
```

As the older pod is deleted, a new pod will be created based on the number of
replicas set in the config.


## Unique IP address

Every pod in k8s cluster has unique internal-to-k8s IP address. By giving each
pod unique IP, k8s simplifies communication and service discovery within
cluster. Pods within the same node or across different nodes can easily
communicate.

All the rsource inside k8s cluster are virtualized. So, if the IP address of a
pod is not the same as the IP address of the node its running on. However, its
virtual IP address is only accessible from within the cluster.

In order to obtain more pods information:

```bash
kubectl get pods -o wide
```

We can see that there are unique IP address of each pod.

#### Proxy server

```bash
kubectl proxy
```

This will start a proxy server on our local machine, probably on
`http://localhost:8001`. Assuming that this is the host, navigate to
`http://localhost:8001/api/v1/namespaces/default.pods` in the browser and we
should be able so see JSON blob describing the pods that are running.
