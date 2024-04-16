# Configuring a cgroup driver

The container runtimes says that `systemd` driver is recommended for kubeadm
based setups instead of the kubelet's default `cgroupfs` driver, because kubeadm
manages the kubelet as a systemd service.

## Configuring the kubelet cgroup driver

kubeadm allows us to pass a `KubeletConfiguration` structure during `kubeadm
init`. This `KubeletConfiguration can include the `cgroupDriver` field which
controls the cgroup driver of the kubelet.

> [!NOTE]
> In v1.22 and later, if the user does not set the `cgroupDriver` field until
> `KubeletConfiguration`, kubeadm defaults it to `systemd`.
>
> In k8s v1.28, we can enable automatic detection of the cgroup driver as an
> alpha feature.

A minimal example of configuring the field explicitly:

```yaml
# kubeadm-config.yaml
kind: ClusterConfiguration
apiVersion: kubeadm.k8s.io/v1beta3
kubernetesVersion: v1.21.0
---
kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
cgroupDriver: systemd
```

Such configuration file can then be passed to the kubeadm command:

```bash
kubeadm init --config kubeadm-config.yaml
```

> [!NOTE]
> kubeadm uses the same `KubeletConfiguration` for all nodes in the cluster. The
> `KubeletConfiguration` is stored in a ConfigMap object under the `kube-system`
> namespace.
> 
> Executing the subcommands `init`, `join`, and `upgrade` would result in
> kubeadm writing the `KubeletConfiguration` as a file under
> `/var/lib/kubelet/config.yaml` and passing it to the local node kubelet.

## Using the `cgroupfs` driver

To use `cgroupfs` and to prevent `kubeadm upgrade` from modifying the 
`KubeletConfiguration` cgroup driver on existing setups, we must be explicit 
about its value. This applies to a case where we do not wish future versins of
kubeadm to apply the `systemd` driver by default.

Refer the documentation of the selected container runtime if we want to
configure a container runtime to use the `cgroupfs` driver.

## Migrating to the `systemd` driver

To change the cgroup driver of an existing kubeadm cluster from `cgroupfs` to
`systemd` in-place, a similar procedure to a kubelet upgrade is required. This
must include both steps outlined below.

> [!NOTE]
> Alternatively, it is possible to replace the old nodes in the cluster with new
> one that uses the `systemd` driver. This requires executing only the first
> step below before joining new nodes and ensuring the workloads can safely move
> to the new nodes before deleting the old nodes.

### Modify the kubelet ConfigMap
- Call `kubectl edit cm kubelet-config -n kube-system`.
- Either modify the existing `cgrouupDriver` value or add a new field that looks
  like this:

  ```bash
  cgroupDriver: systemd
  ```

