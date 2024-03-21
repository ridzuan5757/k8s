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

# Reource Limits

None of our current deployments have any resource limits set. We have very
little traffic, so it is not currently an issue, but in production environment,
we would want to set resource limits to ensure that our pods do not consume too
many resources. We would not want a pod to hog all the CPU and RAM on its node,
suffocating all of the other pods in the node.

## Setting limits

We can set resource limits in the deployment files.

```yaml
spec:
    containers:
    - name: <container_name>
      image: <image_name>
      resources:
        limits:
            memory: <max-memory>
            cpu: <max-cpu>
```

Memory is measured in bytes, so we can use the suffixes `Ki`, `Mi`, and `Gi` to
spcify kilobytes, megabytes and gigabytes, respectively. For example `512Mi` is
512 megabytes. 

CPU is measured in cores, so we can use the suffix `m` to specify milli-cores.
For example, `500m` is 500 milli-cores, or 0.5 cores.

# Limit break

We may have noticed that with the `testcpu` application, we never inform the
application how much CPU to use. That is because generally speaking, application
do not know how much CPU they should use. They just go as fast as they can when
they are doing computations.

Memory is different, applications allocate memory based on variety of factors,
and while an application can have its CPU throttled and just go slower, if an
application runs out of available memory, it will crash.

# Horizontal Pod Autoscaling HPA

A horizontal pod autoscaler can automatically scale the number of pods in a
deployment based on observed CPU utilization or other custom metrics. It is very
common in k8s environment to have a low number of pods in a deployment, and then
scale up the number of pods automatically as CPU usage increases. To implement
HPA:

First, delete the `replicase: x` parameter in the deployment file. This will
allow the new autoscaler to have full control over the number of pods. Create a
new YAML file for HPA.

```yaml
apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
    name: testcpu-hpa
spec:
    scaleTargetRef:
        apiVersion: apps/v1
        kind: Deployment
        name: synergychat-testcpu
    minReplicas: 1
    maxReplicas: 4
    targetCPUUtilizationPercentage: 50
```

This HPA will monitor the CPU usage of the pods in the `testcpu` deployment. Its
goal is to scale up or down the number of pods in the deployment so that the
average CPU usage of all pods is around 50%. As CPU usage increases, it will add
more pods. As CPU usage decreases, it will remove pods.

After applying the hpa, monitor the number of pods as they scale up:

```bash
kubectl apply -f testcpu-hpa.yaml
kubectl get pods
kubectl top pods
```

An hpa is just another resource, so we can use the following command to see the
current state of the autoscaler.

```bash
kubectl get hpa
```
