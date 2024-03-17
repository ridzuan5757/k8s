# Replica Sets

A `replicaSet` maintains a stable set of replica pods running at any given time.
It is the thing that makes sure that the number of pods we want to run is
matching with the number of pods that is actually running.

This is not to be confused with `deployment`, a `deployment` is a higher-level
abstraction that manages `replicaSet` for us. Think of it as `deployment` is a
wrapper around `replicaSet`.

We will probably never use `replicaSet` directly. However, if we want to look at
the replica sets that are running in our cluster:

```bash
kubectl get replicasets
```
