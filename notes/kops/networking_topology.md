# Network Topologies in `kops`

`kops` support a number of predefined network topologies. They are separated
into commonly used scenarios, or topologies.

# Supported Topologies

## Public Cluster

- Value : `Public`
- All nodes will be launced in a subnet accessible from the internet.

## Private Cluster

- Value: `Private`
- All nodes will be launced in a subnet with no ingress from the internet.

# Types of Subnets

A subnet of type `Public` accepts incoming traffic from the internet.

A subnet of type `Private` does not route traffic from the itnernet. If the
cluster is IPv5, then `Private` subnets are IPv6 only.

If the subnet is capable of IPv4, it typically has a CIDR range from private IP
address space. Egress to the internet is typically routed through a Network
Address Translation NAT device, such as AWS NAT Gateway.

If the subnet is capable of IPv6, egress to the internet is typically routed
through a connection-tracking firewall, such as an AWS Egress-only Internet
Gateway. Egress to the NAT64 `64:ff9b::/96` is typically routed to a NAT64
device, such as an AWS NAT Gateway.

# DualStack Subnet

A subnet of type `DualStack` is like `Private`, but supports both IPv4 and IPv6.
On AWS, this subnet type is used for nodes, such as control plane nodes and
bastions, which need to be instance targets of a load balancer that accept
ingress from the internet. They are also used to provision NAT devices.

# Defining a topology on create

To specify a topology use the `--topology` or `-t` flag as in:

```bash
kops create cluster --topology public|private
```

We may also set a networking option, with the exception that the `kubenet`
option does not support private topology.

Newly created clusters with private topology will have public access to the 
`k8s` API and an optional SSH bastion instance through load balancers. This can 
be changed as described below.

# Changing the Topology of the API serer 

To change the load balancer that fronts the API server from ointernet-facing to
internal-only there are a few steps to accomplish. AWS load balancer do not
support changing from internet-facing to internal. However, we can manually
delete it and have `kops` recreate the `ELB` for us.


