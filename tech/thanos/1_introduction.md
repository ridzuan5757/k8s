# Introduction

`thanos` is a set of components that can be composed into `high-availability`
metric system with unlimited storage capacity. It can be added seamlessly on top
existing `prometheus` deployment.

## Starting initial prometheus server

`thanos` is meant to scale and extend vanilla `prometheus`. This means that we
can gradually, without disruption, deploy `thanos` on top of `prometheus` setup.

Let's start this process by first spining up 3 prometheus servers. The real
advantage of `thanos` is when we need to scale out `prometheus` from a single
replica. Some reason for scale-out might be:
- Adding functional sharding because of metrics high cardinality.
- Need for high availability of `prometheus` for example: rooling upgrades.
- Aggregating queries from multiple clusters.

Imagine the following situation:
- We have one prometheus server in some `eu1` cluster.
- We have 2 replica `prometheus` server in some `us1` cluster that scapes the
  same target.

## Prometheus Configuration Files

Now, we will prepare configuration files for all prometheus instances.
