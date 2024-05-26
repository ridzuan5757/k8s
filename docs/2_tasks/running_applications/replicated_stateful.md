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

```yaml
# headless service for stable DNS entries of StatefulSet nenber
apiVersion: v1
kind: Service
metadata:
    name: mysql
    labels:
        app: mysql
        app.kubernetes.io/name: mysql
spec:
    ports:
    - name: mysql
      port: 3306
    clusterIP: None
    selector:
        app: mysql

---

# client service for connecting to any MySQL instance for reads
# for writesm we must instead connect to the primary: mysql-0.mysql.
apiVersion: v1
kind: Service
metadata:
    name: mysql-read
    labels:
        app: mysql
        app.kubernetes.io/name: mysql
        readonly: "true"
spec:
    ports:
    - name: mysql
      port: 3306
    selector:
        app: mysql
```

```bash
kubectl apply -f service.yaml
```

The headless `Service` provides a home for the DNS entries that the `StatefulSet` controllers creates for each Pod that is part of the set. Because the headless 
`Service` is named `mysql`, the Pods are accessible by resolving
`<pod-name>.mysql` from within any other Pod in the same k8s cluster and
namespace.

The client `Service`, called `mysql-read`, is a normal `Service` with its own
cluster IP that distributes connections across all MySQL Pods that report being
`Ready`. The set of potential endpoints includes the primary MySQL server and
all replicas.

Note that only read queries can use the load-balanced client `Service`. Because
there is only one primary MySQL server, clients should connect directly to the
primary MySQL Pod through its DNS entries withing the headless `Service` to
execute writes.

