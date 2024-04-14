# Production environment

## Considerations

A production environment may require:
- access by many users
- consistent availability
- resources to adapt to changing demands

As we decide where the prodution k8s to live and the amount of management we
have to take on or hand to others, consider how the requirements for k8s cluster
are influenced by the following issues:

###### Availability

A single-machine k8s learning enviroment has a single point of failure. A highly
available cluster means considering:
- Separating the control plane from the working nodes.
- Replicating the control plane components on multiple nodes.
- Load balancing traffic to the cluster's API server.
- Having enough worker nodes available, or to be able to quickly available as
  changing workloads warrant it.

###### Scale

If we are expecting the production k8s environment to receive a stable amount of
demand, we might be able to set up for the capacity we need and be done.

However, if we expect demand to grow over time or change dramatically based on
things like season or special events, we need to plan how to scale to relieve
increased pressure from more requests to the control plane and worker nodes or
scale down to reduce unused resources.

###### Security and access management

We have full admin privileges on our own k8s learning cluster. But shared
clusters with important workloads and more than one or two users, require a more
refined approach to who and what can access cluster resources.

We can use role-based access control RBAC and other security mechanism to make
sure that users and workloads can get access the resource that they need, while
keeping workloads and the cluster itself secure.

Limits on the resources that is accessible to users or workloads can be set by
managing policies and container resources.

## Setup

In a production quality k8s cluster, the control plane manages the cluster from
services that can be spread across multiple computers in different ways. EAch
worker node, however, represents a single entity that is configured to run k8s
pods.

### Production control plane


