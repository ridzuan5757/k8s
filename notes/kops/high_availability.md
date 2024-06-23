# High Availability

For testing purposes, `k8s` works just fine with single master. However, when
the master become unavailable, for example due to upgrade or instance failure,
the `k8s` API will not be available. Pods and services that are running in the
cluster continue to operate as long as they do not depend on interacting with
the API, but operations such as:

- Adding nodes
- Scaling pods
- Replacing terminated pods

will not work. Running `kubectl` will also not work.

`kops` runs each master in a dedicated autoscaling groups and stores data on EBS
volumes. That wau, if a master note is terminated the ASG will launch a new
master instance with the master'svolume. Because of the dedicated EBS volumes,
each master is bound to a fixed Availability Zones. If the availability zones
becomes unavailable, the amster instance in the zone will also become
unavailable.

For production use, we therefore want to run `k8s` in HA setup with multiple
masters. With multiple master nodes,w e will able both to do graceful upgrades
and we will able to survive AZ failures.

Very few regions offer less than 3 AZs. In this case, running multiple masters
in the same AZ is an option. If the AZ with multiple masters becomes unavailable
we will still have downtime with this cofniguration. But regular changes to
master nodes such as upgrades will be graceful and without downtime.

If we already have a single master cluster, we would like to convert to a
multi-master cluster.

Note that running clusters spanning several AZs is more expensive than running
cluster spanning one or two AZs. This happens not only becuase of the master EC2
cost, but also because we have to pay for cross-AZ traffic. Depending on the
workload we may therefore also want to consider running worker nodes only in two
AZs. As long as the application do not rely on quorum, we will still have AZ
fault tolerance.

# Creating HA cluster

## Example 1 : Public Topology

The simplest way to get started with a HA cluster is to run `kops create cluster` as shown below.

```bash
kops create cluster \
    --node-count 3 \
    --zones ap-southeast-1a, ap-southeast-1b, ap-southeast-1c \
    --master-zones ap-southeast-1a, ap-southeast-1b, ap-southeast-1c \
    proactivemonitoring.silentmode.com
```

The `--master-zones` flag lists the zones we want the masters to run in. By
defaults, `kops` will create one master per AZ. Since the `k8s` etcd cluster
runs on the master nodes, we have to specify an odd number of zones in order to
obtain quorum.

## Example 2: Private Topology

Using private network topology:

```bash
kops create cluster \
    --node-count 3 \
    --zones us-west-2a,us-west-2b,us-west-2c \
    --master-zones us-west-2a,us-west-2b,us-west-2c \
    --topology private \
    --networking <provider> \
    ${NAME}
```

Note that the default networking provider `kubenet` does not support private
topology.

## Example 3: Multiple masters in the same AZ

If necessary, for example in regions with less than 3 AZs, we can launch
multiple masters in the same AZ.

```bash
kops create cluster \
    --node-count 3 \
    --master-count 3 \
    --zones ap-southeast-1 \
    --master-zones ap-southeast-1 \
    pmsm.k8s.local
```
