# Run a Single-Instance Stateful Application

Single-instance stateful application in k8s using `PersistentVolume` and a
`Deployment`. The application is MySQL.

## Objectives
- Create a `PersistentVolume` referencing a disk in the environment.
- Create a MySQL `Deployment`.
- Expose MySQL to other pods in the cluster at a known DNS name.

## Deployment

We can run a stateful application by creating ak8s `Deployment` and connecting
it to an existing `PersistentVolume` using `PersistentVolumeClaim`. For example,
this YAML file describes a `Deployment` that runs MySQL and references the
`PersistentVolumeClaim`. The file defines a volume mount for `/var/lib/mysql`,
and then creates a `PersistentVolumeClaim` that looks for a 2G volume. This
claim is satisfied by any existing volume that meets the requirements, or by a
dynamic provisioner.

The password is defined in the `config.yaml` and this is insecure. Use k8s
`Secrets` for a secure solution.

```yaml
# mysql-pv.yaml

apiVersion: v1
kind: PersistentVolume
metadata:
    name: mysql-pv-volume
    labels:
        type: local
spec:
    storageClassName: manual
    capacity:
        storage: 2Gi
    accessModes:
        - ReadWriteOnce
    hostPath:
        path: /mnt/data

---

apiVersion: v1
kind: PersistentVolumeClaim
metadata:
    name: mysql-pv-claim
spec:
    storageClassName: manual
    accessModes:
        - ReadWriteOnce
    resources:
        requests:
            storage: 2Gi
```
```yaml
# mysql-deployment.yaml

apiVersion: v1
kind: Service
metadata:
    name: mysql
spec:
    ports:
        - port: 3306
    selector:
        app: mysql
    clusterIP: None

---

apiVersion: apps/v1
kind: Deployment
metadata:
    name: mysql
spec:
    selector:
        matchLabels:
            app: mysql
    strategy:
        type: Recreate
    template:
        metadata:
            labels:
                app: mysql
        spec:
            containers:
            - image: mysql:5.6
              name: mysql
              env:
              - name: MYSQL_ROOT_PASSWORD
                value: password
              ports:
              - containerPort: 3306
                name: mysql
            volumeMounts:
            - name: mysql-persistent-storage
              mountPath: /var/lib/mysql
        volumes:
        - name: mysql-persistent-storage
          persistentVolumeClaim:
              claimName: mysql-pv-claim
```

Deploy the `PersistentVolume` and `PersistentVolumeClaim`:

```bash
kubectl apply -f mysql-pv.yaml
```

Deploy the `Deployment`:

```bash
kubectl apply -f mysql-deployment.yaml
```

Display information about the `Deployment`:

```bash
kubectl describe deployment mysql
```

The output is similar to this:

```bash
Name:                 mysql
Namespace:            default
CreationTimestamp:    Tue, 01 Nov 2016 11:18:45 -0700
Labels:               app=mysql
Annotations:          deployment.kubernetes.io/revision=1
Selector:             app=mysql
Replicas:             1 desired | 1 updated | 1 total | 0 available | 1 unavailable
StrategyType:         Recreate
MinReadySeconds:      0
Pod Template:
  Labels:       app=mysql
  Containers:
    mysql:
    Image:      mysql:5.6
    Port:       3306/TCP
    Environment:
      MYSQL_ROOT_PASSWORD:      password
    Mounts:
      /var/lib/mysql from mysql-persistent-storage (rw)
  Volumes:
    mysql-persistent-storage:
    Type:       PersistentVolumeClaim (a reference to a PersistentVolumeClaim in the same namespace)
    ClaimName:  mysql-pv-claim
    ReadOnly:   false
Conditions:
  Type          Status  Reason
  ----          ------  ------
  Available     False   MinimumReplicasUnavailable
  Progressing   True    ReplicaSetUpdated
OldReplicaSets:       <none>
NewReplicaSet:        mysql-63082529 (1/1 replicas created)
Events:
  FirstSeen    LastSeen    Count    From                SubobjectPath    Type        Reason            Message
  ---------    --------    -----    ----                -------------    --------    ------            -------
  33s          33s         1        {deployment-controller }             Normal      ScalingReplicaSet Scaled up replica set mysql-63082529 to 1
```

List the pods created by the `Deployment`:

```bash
kubectl get pods -l app=mysql
```

The output is similar to this:

```bash
NAME                   READY     STATUS    RESTARTS   AGE
mysql-63082529-2z3ki   1/1       Running   0          3m
```

Inspect the `PersistentVolumeClaim`:

```bash
kubectl describe pvc mysql-pv-claim
```

The output is similar to this:

```bash
Name:         mysql-pv-claim
Namespace:    default
StorageClass:
Status:       Bound
Volume:       mysql-pv-volume
Labels:       <none>
Annotations:    pv.kubernetes.io/bind-completed=yes
                pv.kubernetes.io/bound-by-controller=yes
Capacity:     20Gi
Access Modes: RWO
Events:       <none>
```

## Accessing the MySQL instance

The preceding YAML file creates a `Service` that allows other Pods in the cluster
to access the database. The `Service` option `clusterIP: None` lets the `Service` DNS name resolve directly to the Pod's IP address. This is optimal when we have only one Pod behind a Service and we do not intend to increase the number of Pods.

Run a MySql client to connect to the server:

```bash
kubectl run -it --rm --image=mysql:5.6 --restart=Never mysql-client -- mysql -h mysql -ppassword
```

This command creates a new `Pod` in the cluster running a MySQL client and
connects it to the server through the `Service`. If it connects, we know that
the stateful MySQL is up and running.

## Updating

The image or any other part of the `Deployment` can be updated as usual with the
`kubectl apply` command. Here are some precautions specific to stateful apps:
- Do not scale the app. This setup is for single-instance apps only. The
  underlying `PersistentVolume` can only be mounted to one Pod.
- Use `strategy: Recreate` in the `Deployment` manifest. This instructs k8s to
  not use rolling updates. Rolling updates will not work, as we cannot have more
  than one Pod running at a time. The `Recreate` strategy will stop the first
  pod before creating a new one with the updated configuration.

## Deleting a `Deployment`

```bash
kubectl delete deployment,svc mysql
kubectl delete pvc mysql-pv-claim
kubectl delete pv mysql-pv-volume
```

If `PersistentVolume` is manually provisioned, we also need to manually delete
it, as well as release the underlying resource. If we used a dynamic
provisioner, it automatically deletes `PersistentVolume` when it see that we
deleted the `PersistentVolumeClaim`. Some dynamic provisioners such as those for
Elastic Block Store (EBS-AWS) and Persistent Disk (PD-GCP) also release the
underlying resource upon deleting the `PersistentVolume`.

