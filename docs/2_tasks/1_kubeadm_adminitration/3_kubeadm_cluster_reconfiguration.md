# Reconfiguring a kubeadm cluster

kubeadm does not support automated ways of reconfiguring components that were
deployed on managed nodes. One way of automating this would be using a custom
operator.

To modify the components configuration, we must manually edit associated cluster
objects and files on disk.

## Preparation
- Cluster that was deployed using kubeadm.
- Administrator credentials `/etc/kubernetes/admin.conf` and network
  connectivity to a running kube-apiserver in the cluster from a shot that has
  kubectl installed.
- Text editor installed on all hosts.

## Reconfiguring the cluster

kubeadm writes a set of cluster wide component configuration options in
ConfigMaps and other objects. These objects must be manually edited. The command
`kubectl edit` can be used for that.

The `kubectl edit` command will open a text editor where we can edit and save
the object directly.

We can use the environment variables `KUBECONFIG` and `KUBE_EDITOR` to specify
the location of the kubectl consumed kubeconfig file and preferred text editor.

For example:

```bash
KUBECONFIG=/etc/kubernetes/admin.conf KUBE_EDITOR=nano kubectl edit <parameters>
```

> [!NOTE]
> Upon saving any changes to these cluster objects, components runnin on nodes
> may not be automatically updated. The steps below will be used to perform that
> update manually.

> [!WARNING]
> Component configuration in ConfigMaps is stored as unstructured data YAML
> string. This means that validation will not be performed upon updating the
> contents of a ConfigMap.
>
> We have to be careful to follow the documented API format for a particular
> component configuration and avoid introducing typos and YAML indentation
> mistakes.

### Applying cluster configuration changes

#### Updating the `ClusterConfiguration`

During cluster creation and upgrade, kubeadm writes its `ClusterConfiguration`
in a ConfigMap called `kubeadm-config` in the `kube-system` namespace.

To change a particular option in the `ClusterConfiguration` we can edit the
ConfigMap with this command:

```bash
kubectl edit cm -n kube-system kubeadm-config
```

The configuration is located under the `data.ClusterConfiguration` key.

> [!NOTE]
> The `ClusterConfiguration` includes a variety of options that affect the
> configuration of individual components such as kube-apiserver, kube-scheduler,
> kube-controller-manager, CoreDNS, etcd and kube-proxy. Changes to the
> configuration must be reflected on node components manually.

#### Reflecting `ClusterConfiguration` changes on control plane nodes

kubeadm manages the control plane components as static pod manifests located in
the directory `/etc/kubernetes/manifests`. Any changes to the
`ClusterConfiguration` under the `apiServer`, `controllerManager`, `scheduler`
or `etcd` keys must be reflected in the associated files in the manifests
directory on a control plane node. Such changes may include:
- `extraArgs` - require updating the list of flags passed to a component
  container.
- `extraMounts` - require updated the volume mounts for a component container.
- `*SANs` - requires writing new certificates with updated Subject Alternative
  Names.

Before proceeding with these changes, make sure directory `/etc/kubernetes/`.

To write new certificates we can use:

```bash
kubeadm init phase certs <component-name> --config <config-file>
```

To write new manifests files in `/etc/kubernetes/manifests`, we can use:

```bash
# For Kubernetes control plane components
kubeadm init phase control-plane <component-name> --config <config-file>
# For local etcd
kubeadm init phase etcd local --config <config-file>
```

> [!NOTE]
> Updating a file in `/etc/kubernetes/manifests` will tell the kubelet to
> restart the static pod for the corresponding component. Try doing these
> changes one node at a time to leave the cluster without downtime.

### Applying kubelet configuration changes

#### Updating the `KubeletConfiguration`

During cluster creation and upgrade, kubeadm writes its `KubeletConfiguration`
in a ConfigMap called `kubelet-config` in the `kube-system` namespace.

We can edit the ConfigMap with this command:

```bash
kubectl edit cm -n kube-system kubelet-config
```

The configuration is located under the `data.kubelet` key.

#### Reflecting the kubelet changes

To reflect the change on kubeadm nodes, we must do the following:
- Log in to a kubeadm node.
- Run `kubeadm upgrade node phase kubelet-config` to download the latest
  `kubelet-config` ConfigMap contents into the local file
  `var/lib/kubelet/config.yaml`.
- Edit the file `/var/lib/kubelet/kubeadm-flags.env` to apply additional
  configuration with flags.
- Restart the kubelet service with `systemctl restart kubelet`.

> [!NOTE]
> Do these changes one node at a time to allow workloads to be rescheduled
> properly.

> [!NOTE]
> During `kubeadm upgrade`, kubeadm downloads the `KubeletConfiguration` from
> the `kubelet-config` ConfigMap  and overwrite the conents of
> `/var/lib/kubelet/config.yaml`. This means that node local configuration must
> be applied either by flags in `/var/lib/kubelet/kubeadm-flags.env` or by
> manually updating the contents of `/var/lib/kubeket/config.yaml` after
> `kubeadm upgrade`, and then restarting the kubelet.

#### Applying kube-proxy configuration changes

##### Updating the `KubeProxyConfiguration`

During cluster creation and upgrade, kubeadm writes its `KubeProxyConfiguration`
in a ConfigMap in the `kube-system` called `kube-proxy`.

This ConfigMap is used by the `kube-proxy` DAemonSet in the `kube-system`
namespace.

To change a particular option in the `KubeProxyConfiguration`, we can edit the
ConfigMap with this command:

```bash
kubectl edit cm -n kube-system kube-proxy
```

The configuration is located under the `data.config.conf` key.

#### Reflecting the kube-proxy configuration changes

Once the `kube-proxy` ConfigMap is updated, we can restart all kube-proxy pods:

```bash
kubectl get pod -n kube-system | grep kube-proxy
```

Delete a pod with:

```bash
kubectl delete pod -n kube-system <pod-name>
```

> [!NOTE]
> Because kubeadm deploys kube-proxy as a DaemonSet, node specific configuration
> is unsupported.

### Applying CoreDNS configuration changes

#### Updating the CoreDNS Deployment and Service

kubeadm deploys CoreDNS as a Deployment called `codedns` and with a Service
`kube-dns`, both in the `kube-system` namespace.

To update any of the CodeDNS settings, we can edit the Deployment and Service
objects:

```bash
kubectl edit deployment -n kube-system coredns
kubectl edit service -n kube-system kube-dns
```

#### Reflecting the CoreDNS changes

Once the CoreDNS changes are applied, we can delete the CoreDNS pods. To obtain
the pod names:

```bash
kubectl get pod -n kube-system | grep coredns
```

Delete a pod with:

```bash
kubectl delete pod -n kube-system <pod-name>
```

New pods with the updated CoreDNS configuration will be created.

> [!NOTE]
> kubeadm does not allow CoreDNS configuration during cluster creation and
> upgrade. This means that if we execute `kubeadm ugprade apply`, our changes to
> the CoreDNS objects will be lost and must be reapplied.

## Perssting the reconfiguration

During the execution of `kubeadm upgrade` on managed node, kubeadm might
overwrite configuration that was applied after the cluster was created
(reconfiguration).

### Persisting Node object reconfiguration

kubeadm writes Labels, Taints, CRI socket and other information on the Node
object for a particular k8s node. To change any of the contents of this Node
object, we can use:

```bash
kubect edit no <node-name>
```

During `kubeadm upgrade` the contents of such Node might get overwritten. If we
would like to persist our modifications to the Node object after upgrade, we can
prepare a `kubectl patch` and appli it to the Node object:

```bash
kubectl patch no <node-name> --patch-file <patch-file>
```

#### Persisting control plane component reconfiguration

The main source of control plane configuration is the `ClusterConfiguration`
object stored in the cluster. To extend the static pod manifest configuration,
patches can be used.

These patch files must remain as files on the control plane nodes to ensure that
they can be used by the `kubeadm upgrade ... --patches <directory>`. 

If reconfiguration is done to the `ClusterConfiguration` and static pod
manifests on disk, the set of node specific patches must be updated accordingly.

#### Persisting kubelet reconfiguration

Any changes to the `KubeletConfiguration` stored on
`/var/lib/kubelet/config.yaml` will be overwritten on `kub3eadm upgrade` by
downloading the contents of the cluster wide `kubelet-config` ConfigMap. To
persist kubelet node speficic configuration either the file
`/var/lib/kubelet/config.yaml` has to be updated manually post-upgrade or the
file `var/lib/kubelet/kubeadm-flags.env` can include flags.

The kubelet flags override the associated `KubeletConfiguration` options, but
note that some of the flags are deprecated.

A kubelet restart will be required after changing 
`/var/lib/kubelet/config.yaml` or `/var/lib/kubelet/kubeadm-flags.env`.
