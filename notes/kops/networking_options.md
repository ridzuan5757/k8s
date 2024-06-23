# `k8s` Networking Options

`k8s` has a networking model in which Pods and Services have their own IP
addresses. As Pods and Services run on servers with their own IP addresses and
networking, the `k8s` networking model is an abstraction that sits separately
from the underlying servers and networks. A number of options, listed below are
available which implement and manage this abstraction.

## Supported networking options

The following list the various networking providers with regards to `kops`
version. As of `kops` 1.26, the default network provider is Cilium. Prior to the
default is Kubenet.

- AWS VPC
- Calico
- Canal
- Cilium
- Cilium ENI
- Flannel udp
- Flannel vxlan
- Kopeio
- Kube-router
- Kubenet
- Lyft VPC
- Romana
- Weave

## Specifying network option for cluster creation

We can specify the network provider via the `--networking` command line switch.
However, this will only give a default configuration of the provider. Typically
we would often modify the `spec.networking` section of the cluster spec to
configure the provider further.

## Container Network Interface CNI

CNI provide specification and libraries for writing plugins to configure network
interfaces in Linux containers. `k8s` has built in support for CNI networking
compoents. Several CNI providers are currently built into `kops`.

- AWS VPC
- Calico
- Canal
- Cilium
- Flannel
- Kube-router

`kops` makes it easy for cluster operators to choose one of these options. The
manifests for the providers are included with `kops`, and we simply use
`--networking <provider_name>`. Replace the provider name with the name lsited
in the provider's documentation when we run `kops cluster create`. For example:

```bash
kops create cluster --networking calico
```

Later, when we run `kops get cluster -o yaml` we will see the option we choose
configured under `spec.networking`.

## Advanced

`kops` makes a best-effort attempt to expose as many configuration options as
possible for the upstream CNI options that it supports within the `kops` cluster
spec. However, as upstream CNI options are always changing, not all options may
be available, or we may wish to use a CNI option which `kops` does not support.

There may also be edge cases to operating a given CNI that were not considered
by the `kops` maintainers. Allowing `kops` to manage the CNI installation is
sufficient for the vast majority of production clusters; however, if this is not
true, then `kops` provides an escape hatch that allows us to take greater
control over the CNI installation.

When using the flag `--networking cni` on `kops create cluster` or
`spec.networking: cni {}`, `kops` will not install any CNI at all, but expect
that we manually install it.

If we try to create a new cluster in this mode, the master nodes will come up in
`not ready` state. We then be able to deploy any CNI DaemonSet by following the
vanilla `k8s` install instructions. Once the CNI DaemonSet has been deployed,
the master nodes should enter `ready` state and the remaining nodes should join
the cluster shortly thereafter.


