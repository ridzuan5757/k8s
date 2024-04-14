# User namespaces

A `namespace` isolates the user running inside the container from the one in the
host.

A process running as root in a container can run as a different non-root user in
the host; in other words , the process has full privileges for operations inside
the user namespace, but is unpriveleged for operations outside the namespace.

This feature can be userd to reduce the damage a compromised container can do to
the host or other pods in the same node. There are several security
vulnerabilities rated either **HIGH** or **CRITICAL** that were not exploritable
when user namespace is active.

## Note on linux system

This is Linux-only feature and support is needed in Linux for idmap mounts on
the filesystems used. This means:
- On the node, the filesystem we use for `/var/lib/kubelet/pods`, or the custom
  directory we configure for this, needs `idmap` mount support.
- All the filesystems used in the pod's volumes must support `idmap` mounts.

In practice, this means we need at least Linux 6.3, as `tmpfs` started
supporting `idmap` mounts in that version. This usually needed as several k8s
features use `tmpfs` (the service account token that is mounted by default uses
a `tmpfs`, `Secrets` use a `tmpfs`, etc).

Some popular filesystems that support `idmap` mounts in Linux 6.3 are:
- `btrfs`
- `ext4`
- `xfs`
- `fat`
- `tmpfs`
- `overlayfs`

In addition, support is needed in the container runtime to use this feature with
k8s pods:
- CRI-O: version 1.25 and later supports user namespace for containers.

## Introduction

User namespaces is a Linux feature that alows to map users in the container to
different users in the host. Furthermore, the capabilities granted to a pod in a
user namespace are valid only in the namespace and void outside of it.

A pod can opt-in to use user namespapces by setting the `pod.spec.hostUsers`
field to `false`.

The kubelet will pick host UIDs/GIDs a pod is mapped to, and will do so in a way
to guarantee that no two pods on the same node use the same mapping.

The `runAsUser`, `runAsGroup`, `fsGroup`, etc. fields in the `pod.spec` always
refer to the user inside the container.

The valid UIDs/GIDs when this feature is eabled is the range 0-65536. This
applies to files and processes (`runAsUser`, `runAsGroup`, etc.).

Files using a UID/GID outside this range will be seen as belonging to the
overflow ID, usually 65534 (configured in `/proc/sys/kernel/overflowuid` 
and `/proc/sys/kernel/overflowguid`). However, it is not possible to modify
those files, even by running as the 65534 user/group.

Most applications that need to run as root but do not access other host
namespaces or resources, should continue to run fine without any changes needed
if user namespace is activated.

## Unserstanding user namespaces for pods

Several container runtimes with their default configuration such as Docker
Engine, containerd, CRI-O use Linux namespaces for isolation. Other technologies
exist and can be used with those runtimes too (for example Kata Containers uses
VMs instead of Linux namespaces). This page is applicable for container runtimes
using Linux namespaces for isolation.

When creating a pod, by default, several new namespaces are used for isolation:
- A network namespace to isolate the network of the container.
- A PID namespace to isolate the view of processes, etc.
- If a user namespace is used, this will isolate the users in the container
  from the users in the node.

This means containers can run as root and be mapped to a non-root user on the
host. Inside the container, the process will think it is running as root (and
therefore tools like `apt`, `yum`, etc work fine), while in reality the process
does not have privileges on the host.
- This can be verified by checking which user the container process is by
  executing `ps aux` form the host.
- The user `ps` shows is not the same as the user you see if you execute inside
  the command `id`.

This abstraction limits what can happen, for example, if the container manages
to escape to the host. Given that the container is running as a non-privileged
user on the host, it is limited what it can do to the host.

Furthermore, as users on each pod will be mapped to different non-overlapping
users in the host, it is limited what they can do to other pods too.

Capabilities granted to a pod are also limited to the pod user namespace and
mostly invalid out of it, some are even compeltely void. Here are 2 examples:
- `CAP_SYS_MODULE` does not have any effect if granted to a pod user user
  namespaces, the pod is not able to load kernel modules.
- `CAP_SYS_ADMIN`is limited to the pod's user namespace and invalid otside of
  it.

Wihtout using a user namespace a container running as root, in the case of a
container breakout, has root privilesges on the node. And if some capability
were granted to the containerm the capabilities are valid on the host too. None
of this is true when we use user namespaces.

## Setting up node to support user namespaces

It is recommended that the host's files and host's processes use UIDs/GUIDs in
the range of 0-65535.

The kubelet will assign UIDs/GUIds higher than that to pods. Therefore, to
guarantee as much isolation as possible, the UIDs/GIDs used by the host's files
and host's processes should be in range 0-65535.

This recommendation is important to mitigate the imapct of CVE likes
CVE-2021-25741, where a pod can potentially read arbitrary files in the hosts.
If the UIDs/GIDs of the pod and the host do not overlap, it is limited what a
pod would be able to do: the pod UID/GID would not match the host's file
owner/group.

# Integration with pod security admission checks

For Linux pods that enable user namespaces, k8s relaxes the application of pod
security standards in a controlled way. This behaviour can be controlled by the
feature gate `UserNamespacesPodSecurityStandards`, which allows an early opt-in
for end users. Admins have to ensure that user namespaces are enabled by all
nodes within the cluster if using the feature gate.

If we enable the associated feature gate and create a pod that uses user
namespaces, the following fields would not be constrained even in contexts that
enfore the Baseile or Restricted pod security standard.

This beaviour does not present a security concern because `root` inside a pod
with user namespaces actually refers to the user inside the container, that is
never mapped to a privileged user on the host.

Here is the list of the field that are not checks for pods in those
curcumstances:

```yaml
spec:
    securityContext:
        runAsNonRoot:
        runAsUser:
    containers:
        - name: <container_1>
          securityContext:
            runAsNonRoot:
            runAsUser:
        - name: <container_2>
          securityContext:
            runAsNonRoot:
            runAsUser:
    initContainers:
        - name: <container_1>
          securityContext:
            runAsNonRoot:
            runAsUser:
        - name: <container_2>
          securityContext:
            runAsNonRoot:
            runAsUser:
    ephemeralContainers:
        - name: <container_1>
          securityContext:
            runAsNonRoot:
            runAsUser:
        - name: <container_2>
          securityContext:
            runAsNonRoot:
            runAsUser:
```

## Limitations

When using a user namespace for the pod, it is disallowed to use other host
napespaces. In particular, if we set `hostUsers: false` then we are not allowed
to set any of:
- `hostNetwork: true`
- `hostIPC: true`
- `hostPID: true`
