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


