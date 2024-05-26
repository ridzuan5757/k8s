# Replicated Stateful Application

This application is a replicated MySQL database. The example topology has a
single primary server and multiple replicas, using asynchronous row-based
replication.

> [!NOTE]
> **This is not a production configuration.** MySQL settings remain on insecure
> defaults to keep the focus on general patterns for running stateful
> applications in k8s.

## Objectives
- Deploy a replicated MySQL topology with a `StatefulSet`.
- Send MYSQL client traffic.
- Observer resistance to downtime.
- Scale the `StatefulSet` up and down.

## Deploy MySQL

The deployment consists of a `ConfigMap`, two `Services` and a `StatefulSet`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
    name: mysql
    labels:
        app: mysql
        app.kubernetes.io/name: mysql
data:
    primary.cnf: |
        # apply this config only on primary
        [mysqld]
        log-bin
    replica.cnf: |
        # apply this config only on replica
        [mysqld]
        super-read-only
```

```bash
kubectl apply -f configmap.yaml
```

This `ConfigMap` provides `my.cnf` overrides that let us independently control
configuration on the primary MySQL server and its replicas. In this case, we
want the primary server to be able to serve replication logs to replicas and we
want replicas to reject any writes that do not come via replication.

There is nothing special about the `ConfigMap` itself that causes different
portions to apply to different `Pods`. Each `Pod` decides which portion to look
at as it is initializing, based on information provided by the `StatefulSet`
controller.


