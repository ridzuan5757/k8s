# Disruption

Pods do not disappear until someone be a person or controller destroys them, or
there is unavoidable hardware or system software error.

## Involuntary disruption

We call these unavoidable cases involuntary disruption to an application. For
examples:
- A hardware failure of the physical machine backing the node.
- Cluster administrator deletes VM instance by mistake.
- Cloud provider or hypervisor failure making VM disappear.
- Kernel panic.
- The node disappear from the cluster due to cluster network partition.
- Eviction of a pod due to the node being out-of-resources.

Except for the out-of-resources condition, all these conditions should be
familiar for most users as they are not k8s specific.

## Voluntary disruption

These cinlude both actions initiated by the application owner and those
initiated by a cluster administrator. Typical application owner actions include:
- Deleting the deployment or other controller that manages the pod.
- Updating a deployment's pod template causing a restart.
- Directly deleting a pod.

Cluster administrator actions include:
- Draining a node for repair or upgrade.
- Draining a node from a cluster to scale the cluster down.
- Removing a pod from a node to permit soemthing else to fit on that node.

These actions might be taken directly by the cluster administrator, or by
automation run by the cluster administrator, or by the cluster hosting provider.

## Dealing with disruptions

Here are some ways to mitigate involuntary disruptions:
- Ensure the pod requests the resources it needs.
- Replicate the application if we need higher availability.
- For even higher availability when running replicated applications, spread
  applications across rachs using anti-affinity or across zones using a
  multi-zone cluster.

The frequency of voluntary disruptions varies. On a basic k8s cluster, there are
no automated voluntary disruptions (only user-triggered ones). However, the
cluster administrator or hosting provider may run some additional services which
casue voluntary disruptions.

For example, rolling out node software updates can cause voluntary disruptions.
Also, some implementation of cluster (node) autoscaling may cause voluntary
disruptions to defragment and compact nodes. 

Cluster administrator or hosting provider should have documented what level of
voluntary disruptions, if any, to expect. Certain configuration options, such as
using `PriorityClass` in the pod spec can cause voluntary and involuntary
disruptions.

## Pod disruption budgets

k8s offers features to hep us run highly available applications even when we
introduce frequent voluntary disruptions.

As an application owner, we can create a `PodDisruptionBudget` PDB for each
application. A PDB limits the number of pods of a replicated application that
are down simultaneously from voluntary disruptions.

For example, a quorum-based application would like to ensure that the number of
replicas running is never brought below the number needed for a quorum. A web
front end might want to ensure that the number of replicas serving load never
falls below a certain percentage of the total.

Cluster managers and hosting providers should use tools which respect PDB by
calling the eviction API instead of directly deleting pods or deployments.

For example, the `kubectl drain` subcommand lets us mark a node as going out of
service. When we run `kubectl drain`, the tool tries to evict all of the pods on
the node we are taking out of service. The eviction request that `kubectl`
submits on our behalf may be temporarily rejected, so the tool periodically
reties all failed requests until all pods on the target node are terminated, or
until a configurable timeout is reached.

A PDB specifies the number of replicas that an application can tolerate having,
relative to how many it is intended to have. For example, a deployment which has
`.spec.replicas: 5` is supposed to have 5 pods at any given time. If its PDB
allows for there to be 4 at a time, then the eviction API will allow voluntary
disruption of one but not 2 pods at a time.

The group of pods that comprise the application is specified using a label
selector, the same as the one used by the application's controller (deployment,
stateful-set, etc).

The "intended" number of pods is computed from the `.spec.replicas` of the
workload resource that is managing those pods. The control plane discovers the
owning workload resource by examining the `.metadata.ownerReferences` of the
pod.

Involuntary disruptions cannot be prevented by PDBs; however they do count
against the budget.

Pods which are deleted or unavailable due to a rolling upgrade to an application
do count against the disruption budget, but workload resources such as
`Deployment` and `StatufulSet` are not limited by PDBs when doing rolling
upgrades. Instead, the handling of failures during application updates is
configured in the spec for the specific workload resource.

It is recommended to set `AlwaysAllow` Unhealthy Pod Eviction Policy to the PDB
to support eviction of misbehaving applications during a node drain. The default
behaviour is to wait for the application pods to become healthy before the drain
can proceed.

When a pod is evicted using the eviction API. it is gracefully terminated,
honoring the `terminationGracePeriodSeconds` setting in its PodSpec.

## `PodDisruptionBudget` example

Consider a cluster with 3 nodes:
- `node-1`
- `node-2`
- `node-3`

The cluster is running serveral application. One of the has 3 replicas initially
called:
- `pod-a`
- `pod-b`
- `pod-c`

There s also another unrelated pod without a PDB called `pod-x`. Initially, the
pods are laid out as follows:

|`node-1`|`node-2`|`node-3`|
|---|---|---|
|`pod-a` available|`pod-b` available|`pod-c` available|
|`pod-x` available| | |
 
All 3 pods are part of a deployment and they are collectively have a PDB which
requires there be at least 2 of 3 pods to be available at all times.

For example, assume the cluster administrator wants to reboot into a new kernel
version to fix a bug in the kernel. The cluster administrator first tries to
drain `node-1` using `kubectl drain` command. That tool tries to evict `pod-a`
and `pod-x`. This succeeds immediately. Both pods go into the `terminating`
state at the same time. This puts the cluster in this state:

|`node-1`|`node-2`|`node-3`|
|---|---|---|
|`pod-a` terminating|`pod-b` available|`pod-c` available|
|`pod-x` terminating| | |

The deployment notices that one of the pods is terminating, so it creates a
replacement called `pod-d`. Since `node-1` is cordoned, it lands on another
node. Something has also created `pod-y` as a replacement for `pod-x`.

> **Note**: For a `StatefulSet`, `pod-a`, which would be called something like
> `pod-0`, would need to terminate completely before its replacement, which is
> also called `pod-0` but has different UID, could be created. Otherwise, the
> example applies to a `StatefulSet` as well.

Now the cluster is in this state:

|`node-1` draining|`node-2`|`node-3`|
|---|---|---|
|`pod-a` terminating|`pod-b` available|`pod-c` available|
|`pod-x` terminating|`pod-d` starting|`pod-y`|

At some point, the pods terminate, and the cluster looks like this:

|`node-1` drained|`node-2`|`node-3`|
|---|---|---|
| |`pod-b` available|`pod-c` available|
| |`pod-d` starting|`pod-y`|


 
 
 



