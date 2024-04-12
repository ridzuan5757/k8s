# Sidecar Containers

Sidecar contaienrs are the secondary containers that run along with the main
application container within the same pod These contaienrs are used to enhance
or to extend the functionality of the main application container by providing
additional services, or functionality such as logging, monitoring, security, or
data synchronization, without directly altering the primary application code.

## Enabling sidecar containers

Enabled by defualt with k8s v1.29, a feature gate named `SidecarContainers`
allows us to specify `restartPolicy` for containers listed in a pod's
`initContainers` field. 

These restartable sidecar containers are independent with other init containers
and main application container within the same pod. These can be started,
stopped or restarted without affecting the main application container and other
init containers.

## Sidecar containers and pod lifecycle

If an init container is created with its `restartPolicy` to `Always`, it will
start and remain running during the entire life cycle of the pod. This can be
helpful for running supporting services separated from the main application
containers.

If a `readinessProbe` is specified for this init container, its result will be
used to determine the `ready` state of the pod.

Since these containers are defined as init containers, they benefit from the
same ordering and sequential guarantees as other init containers, allowing them
to be mixed with other init containers into complex pod initialization flows.

Compared to regular init containers, sidecars defined within `initContainers`
continue to run after they have started. This is important when there is more
than one entry inside `.spec.initContainers` for a pod.

After a sidecar-style init container is running (the kubelet has set the
`started` status for that init container to true), the kubelet then starts the
next init container from the ordered `.spec.initContainers` list. That status
either becomes tre because there is a process running in the container and no
startup robe defined, or as result of its `startupProbe` succeeding.

Here is an eample of a deployment with 2 containers, one of which is a sidecar
container:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  labels:
    app: myapp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
        - name: myapp
          image: alpine:latest
          command:
            - 'sh'
            - '-c'
            - 'while true; do echo "logging" >> /opt/logs.txt; sleep 1; done'
          volumeMounts:
            - name: data
              mountPath: /opt
      initContainers:
        - name: logshipper
          image: alpine:latest
          restartPolicy: Always
          command: 
            - 'sh'
            - '-c',
            - 'tail -F /opt/logs.txt'
          volumeMounts:
            - name: data
              mountPath: /opt
      volumes:
        - name: data
          emptyDir: {}
```

This feature is also useful for running jobs with sidecars, as the sidecar
container will not prevent the job from completing after the main container has
finished.

Here is an example of a job with 2 containers, one of which is a sidecar:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: myjob
spec:
  template:
    spec:
      containers:
        - name: myjob
          image: alpine:latest
          command: 
            - 'sh'
            - '-c'
            - 'echo "logging" > /opt/logs.txt'
          volumeMounts:
            - name: data
              mountPath: /opt
      initContainers:
        - name: logshipper
          image: alpine:latest
          restartPolicy: Always
          command: 
            - 'sh'
            - '-c', 
            - 'tail -F /opt/logs.txt'
          volumeMounts:
            - name: data
              mountPath: /opt
      restartPolicy: Never
      volumes:
        - name: data
          emptyDir: {}
```

## Difference from regular containers

Sidecar containers run alongside regular containers in the same pod. However,
they do not execute the primary application logic; instead, they provide
supporting functionality to the main application.

Sidecar containers have their own independent lifecycles. They can be started,
stopped, and restarted independently of regular containers. this means we can
update, scale or maintain sidecar cotnainers without affecting the primary
application.

Sidecar containers share the same network and storage namespaces with the
primary container. this co-location allows them to interact closely and share
resources.

## Difference from init containers

Sidecar containers work alongside the main container, extending its
functionality and providing additional services.

Sidecar containers run concurrently with the main application container. They
are active throughout the lifecycle of the pod and can be started and stopped
independently of the main container. Unlike init containers, sidecar containers
support probes to control their lifecycle.

These containers can interact directly with the main application containers,
sharing the same network namespace, filesystem, and environment variables. They
work closely together to provide additional functionality.

## 
