# Control groups

Control groups is a group of Linux processes with optional resource isolation,
accounting and limits.

On Linux, control groups constrain resources that are allocated to processes.
The kubelet and the underlying container runtime need to interface with croups
to enfore resource management for pods and containers which includes cpu /
memory / requests and limits for containerized workloads.

There are 2 versions of croups in Linux:
- cgroup v1
- cgroup v2

## `cgroup v2`

cgroup v2 is the next version of the Linux `cgroup` API. cgroup v2 provides a
unified control system with enhanced resource management capabilities.

cgroup v2 offers several improvements over cgroup v1, such as the following:
- Single unified hierarchy design in API.
- Safer sub-tree delegation to containers.
- Newer feature like pressure stall information.
- Enhanced resource allocation management and isolation across multiple
  resources:
    - Unified accounting for different types of memory allocations (network
      memory, kernel memory, etc).
    - Accounting for non-immediate changes such as page cache write backs.

Some k8s features exclusively use cgroup v2 for enhanced resource management and
isolation. For example, the MemoryQoS feature improves memory QoS and relies on
cgroup v2 primitives.

## Using cgroup v2

The recommended way to use cgroup v2 is to use Linux distribution that enables
and use cgroup v2 by default. 

### Requirements

cgroup v2 has the following requirements:
- OS distribution enables cgroup v2.
- Linux kernel version 5.8 or later.
- Container runtime supports for cgroup v2:
    - containerd v1.4 and later.
    - cri-o v1.20 and later.
- The kubelet and container runtime are configured to use the systemd cgroup
  driver.

### Linux distribution cgroup v2 support.

- Container Optimized OS (since M97)
- Ubuntu (since 21.10, 22.04+ recommended)
- Debian GNU/Linux (since Debian 11 bullseye)
- Fedora (since 31)
- Arch Linux (since April 2021)
- RHEL and RHEL-like distributions (since 9)

We can also enable cgroup v2 manually on our Linux distribution by modifying the
kernel cmdline boot arguments. If the distribution use GRUB,
`systemd.unified_cgroup_hierarchy=1` should be added in `GRUB_CMDLINE_LINUX`
under `/etc/default/grub`, followed by `sudo update-grub`. However, the
recommended approach is to use distribution that already enables cgroup v2 by
default.

### Migrating to cgroup v2

To migrate to cgroup v2, ensure that we meet the requirements, then upgrade to a
kernel version that enables cgroup v2 by default.

The kubelet automatically detects that the OS is running on cgroup v2 and
performs accordingly with no additional configuration required.

There should not be any noticeable difference in the user experience when
switching to cgroup v2, unless users are accessing the cgroup file system
directly, either on the node or from within the containers.

cgroup v2 uses a different API than cgroup v1, so if there are any applications
that directly access the cgroup file system, they need to be updated to newer
versions that support v2. For example:
- Some third party monitoring and security agenst may depend on the cgroup
  filesystem. Update these agents the versions that support cgroup v2.
- If we are running cAdvisor as a stand-alone DaemonSet for monitoring pods and
  containers, update it to version 0.43.0 or later.
- If we are deploying Java applications, prefer to use versions which fully
  support cgroup v2:
    - OpenJDK / HotSpot: jdk8u372, 11.0.16, 15 and later
    - IBM Semeru Runtimes: 8.0.382.0, 11.0.20.0, 17.0.8.0, and later
    - IBM Java: 8.0.8.6 and later
- If we are using the uber automaxprocs package, make sure the version use is
  v1.51.1 or higher.

## Identifying the cgroup version on Linux nodes

The cgroup version depends on the Linux distribution being used and the default
cgroup version configured on the OS. To check which cgrouo version of the
distribution:

```bash
stat -fc %T /sys/fs/cgroup/
```

For cgroup v2, the output is `cgroup2fs` while for cgroup v1, the output is
`tmpfs`.
