# Container Runtime Interface CRI

The CRI is a plugin interface which enables the kubelet to use a wide variety of
container runtimes, without having a need to recompile the cluster compoennts.

We need a working container runtime on each node in the cluster, so that the
kubelet can launch pods and their containers.

The CRI is the main protocol for the communication between the kubelet and
container runtime.

The k8s container runtime interface CRI defines the main gRPC protocol for the
communication between the node components kubelet and container runtime.

## The API

The kubelet acts as a client when connecting to the container runtime via gRPC.
The runtime and image service endpoints have to be available in the container
runtime, which can be configured separately withing the kubelet by using the
`--image-service-endpoint` command line flags.

For k8s v1.29, the kubelet prevers to use CRI v1. If a container runtime does
not support v1 of the CRI, then the kubelet tries to negotiate any older
supported version. The v1.29 kubelet can also negotiate v1alpha2, but this
version is considered deprecated. If the kubelet cannot negotiate a supported
CRI version, the kubelet gives up and does not register as a node.

## Upgrading

When upgrading k8s, the kubelet tries to automatically select the latest CRI
version on restart of the compoennt. If that fails, then the fallback will take
place as mentioned above. If a gRPC re-dial was required because the container
runtime has been upgraded, then the container runtime must also support the
initially selected version or the redial is expected to fail. This require a
restart of the kubelet.
