# Kubernetes Network Model

Every `Pod` in a cluster gets its own unique cluster-wide IP address (one
address per IP address family). This means that we do not need to explicitly
create links between `Pods` and we almost never need to deal with mapping
container ports to host ports.

This creates a clean, backwards-compatible model where `Pods` can be treated
much like VMs or physical hosts from the perspectives of:
- Port allocation
- Port naming
- Service discovery
- Load balancing
- Application configuration
- Migration

Kubernetes imposes the following fundamental requirements on any networking
implementation:
- Pods can communicate with all other pods on any other node without NAT
- Agents on a node (system daemons, kubelet) can communicate with all pods on
  that node.

[!INFO]
> For those platforms that support `Pods` running in the host network such as
> Linux, when pods are attached to the host network of a node they can still
> communicate with all pods on all nodes without NAT.

This model is not only less complex overall, but it is principally compatible
with the desire for Kubernetes to enable low-friction porting apps from VMs to
containers. This is simalar to a VM that had an IP and could talk to other VMs
within the same project.

Kubernetes IP addresses exist at a `Pod` scope. Containers within a `Pod` share
their network namespaces - including their IP address and MAC address. This
means that containers within a `Pod` can all reach each other's ports on
`localhost`. This also means that containers within a `Pod` must coordinate port
usage, but this is no different from processes in a VM. This is called as
"IP-per-Pod" model.

It is possible to requests ports on the `Node` itself which forward to the `Pod`
(this is called host ports), but this is a very niche operation. How that
forwarding is implemented is also a detail of the container runtime. The `Pod`
itself is blind to the existence or non-existence of host ports.

Kubernetes networking addresses four concerns:
- Containers within a Pod use networking to communicate via loopback.
- Cluster networking provides communication between different Pods.
- The Service API lets us expose an application running in Pods to be reachable
  from outside of the cluster.
    - Ingress provides extra functionality specifically for exposing HTTP
      applications, websites and APIs.
    - Gateway API is an add-on that provides an expressive extensibe, and role
      oriented family of API kinds for modeling service networking.
- We can also use Services to publish services only for consumption inside the
  cluster.
