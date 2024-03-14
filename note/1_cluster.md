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

