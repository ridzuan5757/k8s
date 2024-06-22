# Init Containers

"Init Containers" is a specialized containers that run before app containers in
a pod. Init containers can contain utilities or setup scrips not present in an
app image.

We can specify init containers in the pod specification alongside the
`containers` array which desribes the app containers.

In k8s, a "Sidecar Containers" is a container that starts before the main
application and continues to run. This is slightly different with "Init
containers" that run to completion (terminated after the processes completed) 
during pod initialization.

## Concept

A pod can have multiple containers running apps within it, but it can also have
one or more init containers, which are run before the app containers are
started.

Init containers are exactly like regular containers, except:
- Init containers always run to completion.
- Each init container must complete successfully before the next one starts.

If a pod's init container fails, the kubelet repeatedly restarts that init
container until it successds. However, if the pod has a `restartPolicy` of
`Never`, and an init container fails during startup of that pod, k8s treats the
overall pod as failed.

To specify an init container for a pod, add the `initContainers` field into the
pod specification, as an array of `contaienr` items similar to the app
`contaienrs` field and its content,

The status of the init containers is returned in `.status.initContainerStatuses`
field as an array of the container statuses similar to the
`.status.containerStatuses` field.

### Differences from regular containers

Init containers support all the fields and features of app containers,
including:
- resource limits
- volumes
- security settings

However, the resource requests and limits for an init container are handled
differently.

Regular init contaienrs (including sidecar containers) do not support these
field:
- `lifecycle`
- `livenessProbe`
- `readinessProbe`
- `startupProbe`

Init containers must run to completion before the pod can be ready; sidecar
containers continue running during a pod's lifetime, and do support some probes.

If we specify multiple init containers for a pod, kubelet runs each init
container sequentially. Each init container must succeed before the next can
run. When all of the init containers have run to completion, kubelet initializes
the application containers for the pod and runs them as usual.

### Differences from sidecar containers

Init containers run and complete their tasks before the main application
container starts. Unlike sidecar containers, init containers are not
continuously running alongside the main containers.

Init containers run to completion sequentially, and the main container does not
start until all the init containers have successfully completed.

Init containers do not support `lifecycle`, `livenessProbe`, `readinessProbe`,
or `startupProbe` whereas sidecar containers support all these probes to control
their lifecycle.

Init containers share the same resources such as CPU, memory and network with
the main application containers but do not interact directly with them. They
can, however, use shared volumes for data exchange.

## Using init containers

Because init containers have separate images from app containers, they have some
advantages for startup related code:
- Init containers can contain utilities or custom code for setup that are not
  present in an app image. For example, there is no need to make an image `FROM`
  another image just to use a tool like `sed`, `awk`, `python`, or `dig` during
  setup.
- The application image builder and deployer roles can work independently
  without the need to jointly build a single app image.
- Init containers can run with a different view of the filesystem than app
  containers in the same pod. Consequently, they can be given access to Secrets
  that app containers cannot access.
- Because init containers run to completion before any app containers start,
  init containers offer a mechanism to block or delay app container startup
  until a set of preconditions are met. Once preconditions are met, all of the
  app containers in a pod can start in parallel.
- Init containers can securely run utilities or custom code that would otherwise
  make an app container image less secure. By keeping unnecessary tools separate,
  we can limit the attack surface of the app container image.

### Use cases

- Wait for service to be created, sing a shell one-line command:

```bash
for i in {1..100}; 
do sleep 1; 
if nslookup myservice; 
then exit 0; 
fi; 
done; 
exit 1
```

- Register this pod with a remote server from the downward API with a command:

```bash
curl -X POST http://$MANAGEMENT_SERVICE_HOST:$MANAGEMENT_SERVICE_PORT/register -d \
    'instance=$(<POD_NAME>)&ip=$(<POD_IP>)'
```

- Wait for some time before starting the app container with command:

```bash
sleep 60
```

- Clone a Git repository into a Volume
- Place values into a configuration file and run a template tool to dynamically
  generate a configuration file for the main app container. For example, place
  the `POD_IP` value in a configuration and generate the main app configuration
  file using Jinja.

### Init containers in use

This example defines a simple pod that has 2 init containers.
- The first waits for `myservice`
- The second waits for `mydb`

Once both containers complete, the pod runs the app container from its `spec`
section.

```yaml
### myapp.yaml

apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app.kubernetes.io/name: MyApp
spec:
  containers:
  - name: myapp-container
    image: busybox:1.28
    command: 
        - 'sh'
        - '-c'
        - 'echo The app is running! && sleep 3600'
  initContainers:
  - name: init-myservice
    image: busybox:1.28
    command:
        - 'sh'
        - '-c'
        - "until nslookup myservice.$(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace).svc.cluster.local; do echo waiting for myservice; sleep 2; done"
  - name: init-mydb
    image: busybox:1.28
    command: 
        - 'sh'
        - '-c'
        - "until nslookup mydb.$(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace).svc.cluster.local; do echo waiting for mydb; sleep 2; done"
```
We can start this pod by running:

```bash
kubectl apply -f myapp.yaml
```

The output is similar to this:

```bash
pod/myapp-pod created
```

And check on its status with:

```bash
kubectl get -f myapp.yaml
```

The output is similar to this:

```bash
NAME        READY     STATUS     RESTARTS   AGE
myapp-pod   0/1       Init:0/2   0          6m
```

or for more details:


```bash
kubectl describe -f myapp.yaml
```

The output is similar to this:

```bash
Name:          myapp-pod
Namespace:     default
[...]
Labels:        app.kubernetes.io/name=MyApp
Status:        Pending
[...]
Init Containers:
  init-myservice:
[...]
    State:         Running
[...]
  init-mydb:
[...]
    State:         Waiting
      Reason:      PodInitializing
    Ready:         False
[...]
Containers:
  myapp-container:
[...]
    State:         Waiting
      Reason:      PodInitializing
    Ready:         False
[...]
Events:
  FirstSeen    LastSeen    Count    From                      SubObjectPath                           Type          Reason        Message
  ---------    --------    -----    ----                      -------------                           --------      ------        -------
  16s          16s         1        {default-scheduler }                                              Normal        Scheduled     Successfully assigned myapp-pod to 172.17.4.201
  16s          16s         1        {kubelet 172.17.4.201}    spec.initContainers{init-myservice}     Normal        Pulling       pulling image "busybox"
  13s          13s         1        {kubelet 172.17.4.201}    spec.initContainers{init-myservice}     Normal        Pulled        Successfully pulled image "busybox"
  13s          13s         1        {kubelet 172.17.4.201}    spec.initContainers{init-myservice}     Normal        Created       Created container init-myservice
  13s          13s         1        {kubelet 172.17.4.201}    spec.initContainers{init-myservice}     Normal        Started       Started container init-myservice
```

To see the logs for the init containers in this pod, run:

```bash
kubectl logs myapp-pod -c init-myservice # Inspect the first init container
kubectl logs myapp-pod -c init-mydb      # Inspect the second init container
```

At this point, those init containers will be waiting to discover services named
`mydb` and `myservice`. Here is the configuration that we can use to make the
service appear:

```yaml
### services.yaml
---
apiVersion: v1
kind: Service
metadata:
  name: myservice
spec:
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9376
---
apiVersion: v1
kind: Service
metadata:
  name: mydb
spec:
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9377
```

To create the `mydb` and `myservice` services:

```bash
kubectl apply -f services.yaml
```

The output is similar to this:

```bash
service/myservice created
service/mydb created
```

We will see that those init contaienrs complete, and the `myapp-pod` pod moves
into the running state:

```bash
kubectl get -f myapp.yaml
```

The output is similar to this:

```bash
NAME        READY     STATUS    RESTARTS   AGE
myapp-pod   1/1       Running   0          9m
```

## Detailed behaviour

During pod startup, the kubelet delays running init containers until the
networking and storage are ready. Then the kubelet runs the pod's init
containers in the order they appear in the pod's spec.

Each init container must exit successfully before the next container starts. If
a container fails to start due to the runtime or exits with failure, it is
retried according to the pod `restartPolicy`. However, if the pod `restartPolicy` 
is set to `Always`, the init containers use `restartPolicy` of `OnFailure`.

A pod cannot be `Ready` until all init container have succeeded. The ports on an
init container are not aggregated under a Service. A pod that is initializing is
in the `Pending` state but should have a condition `initialized` set to false.

If a pod restarts, or is restarted, all init containers must execute again.

Changes to the init container spec are limited to the container image field.
Altering an init container image field is equivalent to restarting the pod.

Because init containers can be restarted, retried, or re-executed, init
container code should be indempotent. In particular, code that writes to files
on `EmptyDirs` should be prepared for the possibility that an output file
already exists.

Init containers have all of the fields of an app container. However, k8s
prohibits `readinessProbe` from being used because init containers cannot define
readiness distinct from completion. This is enforced during validation.

Use `activeDeadlineSeconds` on the pod to prevent init containers from failing
forever. The active deadline includes init contaienrs. However it is reommended
to use `activeDeadlineSeconds` pnly if teams deploy their application as a Job,
because `activeDeadlineSeconds` has an effect even after init container
finished. The pod which is already running correctly would be killed by 
`activeDeadlineSeconds` if we set this value.

The name of each app and init container in a pod must be unique, a validation
error is thrown for any container sharing a name with another.

### Resource sharing within containers

Given the oder of execution of init, sidecar and app containers, the following
rules for resource usage apply:
- The highest of any particular resource request or limit defined on all init
  containers is the effective init request/limit. If any resource has no
  resource limit specified this is considered as the highest limit.
- The pod's effective request/limit for a resource is the higher of:
    - The sum of all app containers request/limit for a resource
    - The effective init request/limit for a resource
- Scheduling is done based on effective request/limit, which means init
  containers can reserve resources for initialization that are not used during
  the life of the pod.
- The quality of service QoS tier fo the pod's effective QoS tier is the QoS
  tier for init containers and app cotnainers alike.

Quota and limits are applied based on the effective pod request and limit. Pod
level control groups (cgroups) are based on the effective pod request and limit,
the same as the scheduler.

### Pod restarts reasons 

A pod can restart, causing re-execution of init contaienrs, for the following
reasons:
- The pod infrastructure container is restarted. This is uncommon and would have
  to be done by someone with root access to nodes.
- All containers in a pod are terminated while `restartPolicy` is set to `Always`, 
  forcing a restart, and the limit container completion record has been lost due
  to garbage collection.

The pod will not be restarted when the init container is changed, or the init
container completion record has been lost due to garbage collection. This
applies to k8s v1.20 and later.

