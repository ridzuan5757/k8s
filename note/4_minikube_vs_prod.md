# Minikube vs production

`Minikube` is a great tool for learning k8s, but it is not production-scale k8s
cluster. The primary difference is that `Minikube` runs on a single-node
cluster, whereas production clusters typically multi-node distributed systems.

# Distributed systems

Whenever we are dealing with system that involves multiple machines talking to
each other over a network, we are dealing with a distributed system. Distributed
systems are inherently complex and k8s is no different, however k8s provides
layer of abstraction on its user.

# Resource and nodes

k8s job is to run software applications, but applications require resources such
as:
- CPU
- Memory
- Disk space

k8s job is to manage these resources and allocate them to the applications that
are running on it.

Consider the following situation:

#### 3 nodes (machine)
|Node|RAM|
|---|---|
|node 1|16GB|
|node 2|8GB|
|node 3|8GB|

#### 5 pods (application)
|App|Required RAM|
|---|---|
|app 1|12GB|
|app 2|2GB|
|app 3|5GB|
|app 4|4GB|
|app 5|4GB|

k8s looks at the resources required by each application and decies which node to
run it on. In this case it might do something like:

|Node|Apps|RAM leftover|
|---|---|---|
|node 1|app1(12GB) app2(2GB)|2GB|
|node 2|app4(4GB) app5(4GB)|0GB|
|node 3|app3(5GB)|3GB|

What happens if we get new application that requires 10GB of RAM? The cluster
does not have resources to run on it. But we can just add another node to the
cluster and let k8s figure out to run it.

# Limitation of `Minikube`

Since we only get 1 node, this setup is no longer working once our machine is
out of resources. k8s clusters are running in production that have thousands of
nodes which is lot of resources would be needed to manage such infrastructure.
