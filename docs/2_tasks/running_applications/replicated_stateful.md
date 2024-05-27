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

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
    name: mysql
spec:
    selector:
        matchLabels:
            app: mysql
            app.kubernetes.io/name: mysql
    serviceName: mysql
    replicas: 3
    template:
        metadata:
            labels:
                app: mysql
                app.kubernetes.io/name: mysql
        spec:
            initContainers:
            - name: init-mysql
              image: mysql:5.7
              command:
              - bash
              - -c
              - |
                set -ex
                
                # generage mysql server id from pod ordinal index
                [[ $HOSTNAME =~ -([0-9]+)$ ]] || exit 1
                ordinal=${BASH_REMATCH[1]}
                echo [mysqld] > /mnt/conf.d/server-id.cnf
                
                # add an offset to avoid reserved server-id=0 value
                echo server-id=$((100 + ordinal)) >> /mnt/conf.d/server-id.cnf

                # copy appropiate conf.d files from config-map to empty dir
                if [[ $ordinal -eq 0]]; then
                    cp /mnt/config-map/primary.cnf /mnt/conf.d/
                else
                    cp /mnt/config-map/replica.cnf /mnt/conf.d/
                fi
              volumeMounts:
              - name: conf
                mountPath: /mnt/conf.d
              - name: config-map
                mountPath: /mnt/config-map
            
            - name: clone-mysql
              image: gcr.io/google-samples/xtrabackup:1.0
              command:
              - bash
              - -c
              - |
                set -ex
                
                # skip the clone if data already exists
                [[ -d /var/lib/mysql/mysql ]] && exit 0

                # skip the clone on primary (ordinal index 0)
                [[ `hostname` =~ -([0-9]+)$ ]] || exit 1
                ordinal=${BASH_REMATCH[1]}
                [[ $ORDINAL -eq 0 ]] && exit 0

                # clone data from previous peer
                ncat --recv-only mysql-$(($ordinal-1)).mysql 3307 | xbstream -x
                -C /var/lib/mysql

                # prepare the backup
                xtrabackup --prepare --target-der=/var/lib/mysql

              volumeMounts:
              - name: data
                mountPath: /var/lib/mysql
                subPath: mysql
              - name: conf
                mountPath: /etc/mysql/conf.d

            containers:
            - name: mysql
              image: mysql:5.7
              env:
              - name: MYSQL_ALLOW_EMPTY_PASSWORD
                value: "1"
              ports:
              - name: mysql
                containerPort: 3306
              volumeMounts:
              - name: data
                mountPath: /var/lib/mysql
                subPath: mysql
              - name: conf
                mountPath: /etc/mysql/conf.d
              resources:
                requests:
                    cpu: 500m
                    memory: 1Gi
              livenessProbe:
                exec:
                    command: ["mysqladmin", "ping"]
                initialDelaySeconds: 30
                periodSeconds: 10
                timeoutSeconds: 5
              readinessProbe:
                exec:
                    # check we can execute queries over TCP
                    # skip networking is off
                    command:
                    - mysql
                    - -h
                    - 127.0.0.1
                    - -e
                    - SELECT 1
                initialDelaySeconds: 5
                periodSeconds: 2
                timeoutSeconds: 1
            - name: xtrabackup
              image: gcr.io/google-samples/xtrabackup:1.0
              ports:
              - name: xtrabackup
                containerPort: 3307
              command:
              - bash
              - -c
              - |
                set -ex
                cd /var/lib/mysql

                # determine binlog position of cloned data, if any
                if [[ -f xtrabackup_slave_info && "x$(<xtrabackup_slave_info)" != "x" ]]; then
                    
                    # xtrabackup already generated a partial "CHANGE MASTER TO"
                    # query
                    # because we are cloning from existing replica
                    # (Need to remote the tailing semicolon)
                    cat xtrabackup_slave_info | sed -E 's/;$/g' > change_master_to.sql.in
                    
                    # ignore xtrabackup_binlog_info in this case
                    # it is useless
                    rm -f xtrabackup_slave_info xtrabackup_binlog_info

                elif [[ -f xtrabackup_binlog_info ]]; then
                    
                    # we are cloning directly from primary
                    # parse binlog position
                    [[ `cat xtrabackup_binlog_info` =~ ^(.*?)[[:space:]]+(.*?)$ ]] || exit 1
                    rm -f xtrabackup_binlog_info xtrabackup_slave_info
                    echo "CHANGE MASTER TO MASTER_LOG_FILE='${BASH_REMATCH[1]}',\
                        MASTER_LOG_POS=${BASH_REMATCH[2]}" > change_master_to.sql.in
                fi
                
                # check if we need to complete a clone by starting replication
                if [[ -f change_master_to.sql.in ]]; then
                    echo "Waiting for mysqld to be ready"
                    until mysql -h 127.0.0.1 -e "SELECT 1"; do sleep 1; done

                    echo "Initializing replication from clone position"
                    mysql -h 127.0.0.1 \
                          -e "$(<change_master_to.sql.in), \
                                MASTER_HOST='mysql-0.mysql', \
                                MASTER_USER='root', \
                                MASTER_PASSWORD='', \
                                MASTER_CONNECT_RETRY=10; \
                            START SLAVE;" || exit 1

                    # in case of container restart, attempt this at most once
                    mv change_master_to.sql.in change_master_to.sql.orig 
                fi

                # start a server to send backups when requested by peers
                exec ncat --listen --keep-open --send-only --max-conns=1 3307 \
                    -c "xtrabackup --backup --slave-info --stream=xbstram 
                    --host=127.0.0.1 --user=root"
            volumeMounts:
            - name: data
              mountPath: /var/lib/mysql
              subPath: mysql
            - name: conf
              mountPath: /etc/mysql/conf.d
            resources:
                requests:
                    cpu: 100m
                    memory: 100Mi
            volumes:
            - name: conf
              emptyDir: {}
            - name: config-map
              configMap:
                name: mysql
    volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes:
            - ReadWriteOnce
        requests:
            storage: 10Gi
```

```bash
kubectl apply -f deployment.yaml
```

You can watch the startup progress by running:

```bash
kubectl get pods -l app=mysql --watch
```

After a while, we should see all 3 Pods become `Running`:

```bash
NAME      READY     STATUS    RESTARTS   AGE
mysql-0   2/2       Running   0          2m
mysql-1   2/2       Running   0          1m
mysql-2   2/2       Running   0          1m
```

> [!NOTE]
> If we do not see any progress, make sure we have a dynamic `PersistentVolume`
> provisioner enabled, as mentioned in the prerequisites.
