# Cluster
- K8s coordinates a highly available cluster of computers that are connected to
  work as a single unit.
- K8s automates the distribution and scheduling of application containers across
  cluster in a more efficient way.


K8s cluster consists of 2 types of resource:
- **Control Plane** coordinates the cluster.
- **Nodes** are the workers that run the applications.

## Control Plane
- Responsible for managing the cluster.
- This covers:
    - Scheduling applications.
    - Maintaining applications' desired state.
    - Scaling applications.
    - Rolling out new update.

## Node
Node is a virtual machine or physical computer that serves as a worker machine
in a K8s cluster.
- Each node has **kubelet**, which is an agent for managing the node and
  communicating with the k8s control plane.
- The node should also have tools for handling containers operations (
  containerd, CRI-O).
- K8s clusters that handles production traffic should have a minimum of 3 nodes
  because if one node goes down, both an etcd (distributed reliable key-value
  store) and control plane instance are lost, compromising redundancy.

When deploying apps on k8s, we tell the control plane to start the application
containers. The control plane shecdules the containers to run on the cluster's
nodes.
