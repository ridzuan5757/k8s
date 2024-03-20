# Storage in k8s

By default, containers running in pods on k8s have access to the filesystem, but
there are some big limitations to this. Even though we are saving the file to
the filesystem, it will not persists once the pod destroyed and recreated. In
other word, the filesystem is ephemeral as the pods.

This has to do with the philosophy behind k8s and even containers in general:
when we spin up a new one, it should always be a blank state, which makes
reproducing and debugging much easier since we don't have to maintain the state
consistency.

# Ephemeral Volumes
On-disk files in a container are ephemeral in nature. This presents some
problems for applications that want to save long-lived data across restarts. For
example user data in a database.

The k8s volume abstraction solves two primary problem:
- Data persistence.
- Data sharing across containers.

As it turns out, there are lot of different types of volumes in 8s. Some are
even ephemeral as well, just like a container's standard filesystem. the primary
reason for using an ephemeral volume is to share data between containers in a
pod.

### Containers scaling

Consider a service that is continuously doing crawling job and exposes the
information that it finds via a JSON API. The data is then made available via
slash commands in other applicatioons. Assumes that the crawler is pretty slow
by default. Each instance only crawls 1 JSON object every 30 seconds.

We can speed it up by increasing the number of concurrent crawlers. The trouble
with scaling up beyond one instance is that each crawler currently stores its
data in memory. We need all pods to share the same data so they can each add
their findings to the same database.

We can try update the crawler deployment to use a **volume** that will be shared
across all containers in the crawler pod and scale up the number of containers
in the pod.

In the crawler deployment file, add `volumes` section to `spec/template/spec`:

```yaml
spec:
    template:
        spec:
            volumes:
                - name: cache-volume
                  emptyDir{}
```

Add a new `volumeMounts` section to the container entry. This will mount the
volume we just created at the `/cache` path.

```yaml
spec:
    template:
        spec:
            containers:
                - name: synergychat-crawler-1
                  image: bootdotdev/synergychat-crawler:latest
                  envFrom:
                    - configMapRef:
                        name: synergychat-crawler-configmap
                  volumeMounts:
                    - name: cache-volume
                      mountPath: /cache
```

Duplicate the crawler containers, says 3 and update each of their individual
name.

```yaml
spec:
    template:
        spec:
            containers:
                - name: synergychat-crawler-1
                  image: bootdotdev/synergychat-crawler:latest
                  envFrom:
                    - configMapRef:
                        name: synergychat-crawler-configmap
                  volumeMounts:
                    - name: cache-volume
                      mountPath: /cache
                - name: synergychat-crawler-2
                  image: bootdotdev/synergychat-crawler:latest
                  envFrom:
                    - configMapRef:
                        name: synergychat-crawler-configmap
                  volumeMounts:
                    - name: cache-volume
                      mountPath: /cache
                - name: synergychat-crawler-3
                  image: bootdotdev/synergychat-crawler:latest
                  envFrom:
                    - configMapRef:
                        name: synergychat-crawler-configmap
                  volumeMounts:
                    - name: cache-volume
                      mountPath: /cache
```

Now all the containers in the pod will share the same volume at `/cache`. It's
just empty directory, but the crawler will use it to store its data.

Add a `CRAWLER_DB_PATH` environment variable to the crawler's `ConfigMap`. Set
it to `cache/db`. The crawler will use a directory called `db` inside the volume
to store its data.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: synergychat-crawler-configmap
data:
  CRAWLER_PORT: "8080"
  CRAWLER_KEYWORDS: love,hat,joy,sadness,anger,disgust,fear,surprise
  CRAWLER_DB_PATH: /cache/db
```

Apply the new `ConfigMap` and `Deployment` and use `kubectl get pod` to see the
status of the new pod.

```bash
kubectl apply -f crawler-configmap.yaml
kubectl apply -f crawler-deployment.yaml
```

We should notice that there is a proble with the pod. only 1 out of the 3
containers is in `ready` state. Use the `logs` command to get the logs for all 3
containers:

```bash
kubectl logs <crawler-pod-name> --all-containers
```

We should see something like:

```bash
listen tcp :8080: bind: address already in use
```

Because pods share the same network namespace, they cannot all bind to the same
port. We can remedy this by binding each container to a different port. `8080`
is the only one that will be exposed via the service. We will be using `8081`
and `8082` for the second and third crawler containers respectively.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: synergychat-crawler-configmap
data:
  CRAWLER_PORT: "8080"
  CRAWLER_KEYWORDS: love,hat,joy,sadness,anger,disgust,fear,surprise
  CRAWLER_DB_PATH: /cache/db
  CRAWLER_PORT_2: "8081"
  CRAWLER_PORT_3: "8082"
```

Change the second and third containers to map `CRAWLER_PORT_2 -> CRAWLER_PORT`.
`CRAWLER_PORT_3 -> CRAWLER_PORT` respectively. Since we have to use `env`
instead of `envFrom` this time, we also have to continue exposing the
`CRAWLER_KEYWORDS` and `CRAWLER_DB_PATH`

```yaml
template:
    spec:
      containers:
        - name: synergychat-crawler-1
          image: bootdotdev/synergychat-crawler:latest
          envFrom:
            - configMapRef:
                name: synergychat-crawler-configmap
          volumeMounts:
            - name: cache-volume
              mountPath: /cache
        - name: synergychat-crawler-2
          image: bootdotdev/synergychat-crawler:latest
          volumeMounts:
            - name: cache-volume
              mountPath: /cache
          env:
            - name: CRAWLER_PORT
              valueFrom:
                configMapKeyRef:
                  name: synergychat-crawler-configmap
                  key: CRAWLER_PORT_2
            - name: CRAWLER_KEYWORDS
              valueFrom:
                configMapKeyRef:
                  name: synergychat-crawler-configmap
                  key: CRAWLER_KEYWORDS
            - name: CRAWLER_DB_PATH
              valueFrom:
                configMapKeyRef:
                  name: synergychat-crawler-configmap
                  key: CRAWLER_DB_PATH
        - name: synergychat-crawler-3
          image: bootdotdev/synergychat-crawler:latest
          volumeMounts:
            - name: cache-volume
              mountPath: /cache
          env:
            - name: CRAWLER_PORT
              valueFrom:
                configMapKeyRef:
                  name: synergychat-crawler-configmap
                  key: CRAWLER_PORT_3
            - name: CRAWLER_KEYWORDS
              valueFrom:
                configMapKeyRef:
                  name: synergychat-crawler-configmap
                  key: CRAWLER_KEYWORDS
            - name: CRAWLER_DB_PATH
              valueFrom:
                configMapKeyRef:
                  name: synergychat-crawler-configmap
                  key: CRAWLER_DB_PATH


```
# Containers in pods

After all of the 3 crawlers has been deployed, when we run `kubectl get pods`,
we should see something like this:

```bash
synergychat-api-6c7944b5c4-rp2k4      1/1     Running   0          160m
synergychat-crawler-cd4947995-ftqg4   3/3     Running   0          151m
synergychat-web-846d86c444-2m6x7      1/1     Running   0          21h
synergychat-web-846d86c444-gxztt      1/1     Running   0          21h
synergychat-web-846d86c444-s88rz      1/1     Running   0          21h
```

It is important to remember that while it is common for a pod to run just a
single container, multiple containers can run in a single pod. This is useful
when we have containers that need to share resources. In other words, we can
scale up the instance of an application either at the contaioner level or at the
pod level.

# Persistence

All the volumes we have worked with so far have been ephemeral, meaning when the
associated pod is deleted the volume is deleted as well. This is fine for some
use cases, but for most CRUD applications we want to eprsist data even if the
pod is deleted.

If we think about it, it is not even just when pods are explicitly deleted with
`kubectl` that we need to worry about data loss. Pods can be deleted for
several reasons:
- The node they are running could fail.
- A new version of the image was published.
- A new node was added to the cluster and the pod was rescheduled.

## Persistent Volumes

Instead of simply adding a volume to a deployment, a persistent volume is a
cluster-level resource that is created separately from the pod and then attached
to the pod. It is similar to a `ConfigMap` in a way. Persistence volume can be
created statically or dynamically.
- Static persistent volume are created manually by the cluster admin.
- Dynamic persistent volume are created automatically when pod requests a volume
  that does not exist yet.

Generally speaking, and especially in the cloud-native world, we want to use
dynamic persistent volume. It is less work and more flexible.

## Persistent Volume Claims

A persistent volume claim is a request for a persistent volume. When using
dynamic provisioning, a persistent volume claims will automatically creates a
persistent volume if one does not exist that matches the claim. The persistent
volume claim is then attached to a pod, just like a volume would be.

#### Implementation

Create a new file called `api-pvc.yaml` and add the following:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
    name: synergy-chat-api
spec:
    accessMode:
        - ReadWriteOnce
    resources:
        requests:
            storage: 1Gi
```

This creates a new persistent volume claim called `synergy-chat-api` with a few
properties that can be read from and written to bey multiple pods ad the same
time. It also requests 1GB of storage.

After applying the persistent volume claim, we can check the persistent volume
and the persistent volume created based on the claim using the following
command:

```bash
kubectl get pvc
kubectl get pv
```

We whould see that a new persistent volume was created automatically. Similarly,
deleting the persistent volume claim will delete the persistent volume created
as well.
