# Runtime Class

RuntimeClass is a feature for selecting the container runtime configuration. The
container runtime configuration is used to run a pod's container.

## Motivation

We can set a different RuntimeClass between different pods to provide balance of
performance versus security. For example, if part of the workload deserves a
high level of information security assurance, we might choose to schedule those
pods so that they run in a container runtime that uses hardware virtualization.

We would then benefit from the extra isolation of the alternative runtime, at
the expense of some additiona overhead.

We can also use RuntimeClass to run different pods with the same container
runtime but with different settings.

## Setup

### Configure the CRI implementation on nodes

The configuration availablle thrpugh RuntimeClass are Container Runtime
Interface CRI implementation dependent.

> RuntimeClass assumes a homogeneous node configuration across the cluster by
> default which means that all nodes are configured the same way with respect to
> container runtimes. To support heterogenous node configuration, scheduling
> should be used.

The configuration hass a corresponding `handler` name, referenced by the
RuntimeClass. The handler must be a valid DNS label name.

### Create the corresponding RuntimeClass resources

The configurations setup in the previous step should each have an associated
`handler` name, which identifies the configuration. For each handler, create a
corresponding RuntimeClass object.

The RuntimeClass resource currently only has 2 significant fields:
- The RuntimeClass name `metadata.name`
- `handler`

The object definition looks like this:

```yaml
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  # The name the RuntimeClass will be referenced by.
  # RuntimeClass is a non-namespaced resource.
  name: myclass 
# The name of the corresponding CRI configuration
handler: myconfiguration 
```

The name of a RuntimeClass object must be a valid DNS subdomain name.

> It is recommended that RuntimeClass write operations create / update / patch /
> delete be restricted to the cluster administrator. This is typically the
> default configuration.

## Usage

Once the RuntimeClass are configured for the cluster, we can specify
`runtimeClassName` in the pod's specification to use it. For example:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  runtimeClassName: myclass
  # ...
```

This will instruct the kubelet to use the named RuntimeClass to run this pod. If
the named RuntimeClass does not exist, or the CRI cannot run the corresponding
handler, the pod will enter the `Failed` terminal phase. Look for a
corresponding event for an error message.

If no `runtimeClassName` is specified, the default RuntimeHandler will be used,
which is equivalent to the behaviour when the RuntimeClass feature is disabled.

### CRI Configuration

#### `containerd`

Runtime handlers are configured throught containerd's configuration at:

```bash
/etc/containerd/config.toml
```

Valid handlers are configured under the runtime section:

```bash
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.${HANDLER_NAME}]
```

#### `CRI-O`

Runtime handlers are configured throught CRI-O's configuration at:

```bash
/etc/crio/crio.conf
```

Valid handlers are configured under the `crio.runtime` table:

```bash
[crio.runtime.runtimes.${HANDLER_NAME}]
  runtime_path = "${PATH_TO_BINARY}"
```

## Scheduling

By specifying the `scheduling` field for a RuntimeClass, we can set constraints
to ensure that Pods running with this RuntimeClass arescheduled to nodes that
support it. If `scheduling` is not set, this RuntimeClass is assumed to be
supported by all nodes.

To ensure pods land on nodes supporting a specific `RuntimeClass`, that set of
nodes should have a common label which is then selected by the
`runtimeclass.scheduling.nodeSelector` field.

The RuntimeClass's `nodeSelector` is merged with the pod's `nodeSelector` in
admission, effectviely taking the intersection of the et of nodes selected by
each. If there is a conflict, pod will be rejected.

If the supported nodes are tainted to prevent other RuntimeClass pods from
running on the node, we can add `toleratins` to the RuntimeClass. As with the
`nodeSelector`, the tolerations are merged with the pod's tolerations in
admission, effectively taking the union set of the ndoes tolerated by each.

## Pod Overhead

We can specify `overhead` resources that are associated with running a pod.
Declaring overhead allows the cluster including the scheduler to account for it
when making decisions about pods and resources.

Pod overhead is defined in RuntimeClass through the `overhead` field. Through
the use of this field, we can specify the overhead of running pods utilizing
this RuntimeClass and ensure these overheads are accounted for in k8s.
