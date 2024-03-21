# `top` Command

WE have already learned how to look at logs for k8s pods, but sometimes that is
not enough when it comes to debugging. Sometimes we want to know about the
resources that a pod is using.

To get metrics working, we need to enable the `metrics-service` addon. Run:

```bash
minikube addons enable metrics-server
```

Take a look inside the `kube-system` namespace:

```bash
kubectl -n kube-system get pod
```

We whould see a new "metrics-server" pod. It might take a couple of minutes to
get started, but once that pod is ready, we should be able to run:

```bash
kubectl top pod
```

We should see something like:

```bash
NAME                               CPU(cores)   MEMORY(bytes)   
synergychat-api-76b796b58d-x5wpk   1m           14Mi            
synergychat-web-846d86c444-d9c8q   1m           15Mi            
synergychat-web-846d86c444-sk6n4   1m           15Mi            
synergychat-web-846d86c444-w2pqg   1m           15Mi
```

The `kubectl top` command (similar to unix top) will show the resources that
each pod is using. In the example above, each pod is using 1milliCPU and 15MB of
memory.

# Vertical and Horizontal Scaling

Generally speaking, there are two ways to scale an application: vertically and
horizontally. Scaling in this context means increasing the capacity of an
application. For example, maybe we have a web server, and to handle roughly 1000
requests per second, it uses about:
- 1/2 of a CPU core
- 1GB of RAM

If we want to scale up to handle 2000 requests per second, we could double up
the CPU and RAM:
- 1 CPU core
- 2GB of RAM

This is called vertical scaling because we are increasing the capacity of the
application by increasing the resources avaialble to it. We are scaling up.
Scaling up usually works until it does not. We can only scale up as much as the
hardware will allow (maximum number of CPUs and amount of RAM the node has).

The other way is to scale horizontally. Instead of increasing the resources
available to the application, we increase the number of instances of the
application pods. Pods can be distributed across nodes, so we can scale
horizontally until we run out of nodes. When working in system like k8s, it is
better to scale horizontally than vertically.
