# Minikube
- Local single node kubernetes cluster.
- Not possible to add other nodes.

# Kubeadm
- Cluster minimal size:
    - Master Node
    - Worker Node
- Worker nodes can be added as many as needed.
- Computationally expensive in laptop.
- Option to select container runtime : Docker, Podman CRI-O etc.

# Kind
- Cluster will be deployed inside docker container.
- Capable to deploy all type of clusters:
    - Single node
    - 1 master & multiple workers
    - Multiple masters & multiple workers
- Cluster are very easy to deploy, however network external access to the
  cluster are more complicated.

# K3S
- Lightweight 
- Does not use docker as default container runtime.



