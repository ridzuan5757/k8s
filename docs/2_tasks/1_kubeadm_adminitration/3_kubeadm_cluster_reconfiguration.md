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
