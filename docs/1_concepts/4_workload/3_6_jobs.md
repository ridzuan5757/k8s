# Jobs 

A Job creates one or more Pods and will continue to retry execution of the Pods
until a specified number of them successfully terminate. As pods successfully
complete, the Job tracks the successful completions. When a specified number of
successful completions is reached, the task (i.e Job) is complete. Deleting
aJonb will clean up the Pods it created. Suspending a Job will delete its active
Pods until the Job is resumed again.

A simple case is to create one Job object in order ot reliable run one Pod to
completion. The Job object will start a new Pod if the first Pod falls or is
deleted for example due to a node hardware failure or a node reboot.

## Running an example Job

```yaml
# job.yaml

apiVersion: batch/v1
kind: Job
metadata:
    name: pi
spec:
    template:
        spec:
            containers:
            - name: pi
              image: perl:latest
              command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
            restartPolicy: Never
    backoffLimit: 4
```

Running the manifest with command `kubectl apply -f job.yaml`, the output is
similar to this:


```yaml
job.batch/pi created 
```

Checking the status of the Job with `kubectl get job pi -o yaml` will return:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  annotations: batch.kubernetes.io/job-tracking 
  creationTimestamp: "2022-11-10T17:53:53Z"
  generation: 1
  labels:
    batch.kubernetes.io/controller-uid: 863452e6-270d-420e-9b94-53a54146c223
    batch.kubernetes.io/job-name: pi
  name: pi
  namespace: default
  resourceVersion: "4751"
  uid: 204fb678-040b-497f-9266-35ffa8716d14
spec:
  backoffLimit: 4
  completionMode: NonIndexed
  completions: 1
  parallelism: 1
  selector:
    matchLabels:
      batch.kubernetes.io/controller-uid: 863452e6-270d-420e-9b94-53a54146c223
  suspend: false
  template:
    metadata:
      creationTimestamp: null
      labels:
        batch.kubernetes.io/controller-uid: 863452e6-270d-420e-9b94-53a54146c223
        batch.kubernetes.io/job-name: pi
    spec:
      containers:
      - command:
        - perl
        - -Mbignum=bpi
        - -wle
        - print bpi(2000)
        image: perl:5.34.0
        imagePullPolicy: IfNotPresent
        name: pi
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Never
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
status:
  active: 1
  ready: 0
  startTime: "2022-11-10T17:53:57Z"
  uncountedTerminatedPods: {}
```

To view completed Pods of a Job, use `kubectl get pods`. To list all the Pods
that belong to a Job in a machine readable form, we can use command like this:

```bash
pods=$(kubectl get pods --selector=batch.kubernetes.io/job-name=pi --output=jsonpath='{.items[*].metadata.name}')
echo $pods
```

The output will be similar to this:

```bash
pi-5rwd7
```

Here, the selector is the same as the selector for the Job. the
`--output=jsonpath` option specifies an expression with the name from each Pod
in the returned list. Viewing the standard output of one of the pods:

```bash
kubectl logs $pods
```

Another way to view the logs of a Job:

```bash
kubectl logs jobs/pi
```

The output is similar to this:

```bash
3.14159
```

## Writing a Job Spec

As with all other k8s config, a Job needs the following field:

```yaml
apiVersion:
kind:
metadata:
spec:
```

When the control plane creates a new Pods for a Job, the `.metadata.name` of the
Job is part of the basis for naming those Pods. The name of a Job must be a
valid DNS subdomain value, but this can produce unexpected results for the Pod
hostnames. For best compatibility, the name should follow more restrictive rules
for a DNS label. Even when the name is a DNS subdomain, the name must be no
longer than 63 characters.

### Job Labels

Job labels will have `batch.kubernetes.io/` prefix for `job-name` and
`controller-uid`.

### Pod Template

The `.spec.template` is the only required field of the `.spec`. The
`.spec.template` is a pod template. It has exactly the same schema as a Pod
except it is nested and does not have `apiVersion` and `kind`.

A Pod template in a Job must specify appropriate labels and an appropriate
restart policy. Only a `RestartPolicy` equal to `Never` or `OnFailure` is
allowed.

### Pod Selector

The `.spec.selector` field is optional. In almost all cases we should not
specify it. 


### Parallel execution for Jobs

There are three main types of task suitable to run as a Job:
- Non-parallel Jobs
    - Normally, only one Pod is started, unless the Pod fails.
    - The Job is complete as soon as its Pod terminates sucessfully.
- Parallel Jobs with a fixed completion count:
    - Specify a non-zero positive value for `.spec.completions`.
    - The Job represents the overall task, and is complete when there are
      `.spec.completions` successful Pods.
    - When using `.spec.completionMode="Indexed"`, each Pod gets a different
      index in the range 0 to `.spec.completions-1`.
- Parallel Jobs with a work queue:
    - Do not specify `.spec.completions` and it will default to
      `.spec.parallelism`.
    - The Pods must coordinate amongst themselves or an external service to
      determine what each should work on. For example, a Pod might fetch a batch
      of up to N items from the work queue.
    - Each Pod is independently capable of determining whether or not all its
      peers are done, and thus that the entire Job is done.
    - When any Pod from the Job terminates with success, no new Pods are
      created.
    - Once at least one Pod has terminated with success and all Pods are
      terminated, then the Job is completed with success.
    - Once any Pod has exited with success, no other Pod should still be doing
      any work for this task or writing any output. They should all be in the
      process of exiting.

For a non-parallel Job, we can leave both `.spec.completions` and
`.spec.parallelism` unset. When both are unset, both are defaulted to 1.


For a fixed completion count Job, we should set `.spec.completions` to the
number of completions needed. We can set `.spec.parallelism`, or leave it unset
and it will default to 1.

For a work queue Job, we must leave `.spec.completions` unset, and set
`.spec.parallelism` to a non-negative integer.

### Controlling parallelism

The requested parallelism `.spec.parallelism` can be set to any non=negative
value. If it is unspecified, it defaults to 1. If it is specified as 0, then the
Job is effectively paused until it is increased.

Actual parallelism (number of pods running at any instant) may be more or less
than requested parallelism, for a variety of reasons:
- For fixed completion count Jobs, the actual number of pods running in parallel
  will not exceed the number of remaining completions. Higher values of
  `.spec.parallelism` are effectively ignored.
- For work queue Jobs, no new Pods are started after any Pod has succeeded -
  remaining Pods are allowed to complete, however.
- If the Job Controller has no had time to react.
- If th eJob Controller failed to create Pods for any reasons (lack of
  `ResourceQuota`, lack of permission, etc.), then there may be fewer pods than
  requested.
- The Job controller may throttle new Pod creation due to excessive previous pod
  failures in the same Job.
- When a Pod is gracefully shut down, it takes time to stop.

### Completion Mode

Jobs with fixed completion count - that is, jobs that have non-null
`.spec.compltions` can have a completion mode that is specifed in
`.spec.completionMode`:
- `NonIndexed` (default) : The Job is considered complete when there have been
  `.spec.completions` successfully completed Pods. In other words, each Pod
  completion is homologous to each other. Note that Jobs that have null
  `.spec.completions` are implicitly `NonIndexed`.
- `Indexed` : The Pods of a Job get an associated completion index from 0 to
  `.spec.completions-1`. The index is available through four mechanism:
    - The Pod annotation `batch.kubernetes.io/job-completion-index`.
    - The Pod label `batch.kubernetes.io/job-completion-index`. Note that the
      feature gate `PodIndexLabel` must be enabled to use this label, and it is
      enabled by default.
    - As part of the Pod hostname, following the pattern `$(job-name)-$(index)`.
      When we use an Indexed Job in combination with a Service, Pods within the
      Job can use the deterministic hostnames to address each other via DNS.
    - From the containerized task, in the environment variable
      `JOB_COMPLETION_INDEX`.

The Job is considered complete when there is one successfully completed Pod for
each index.

> [!NOTE]
> Although rare, more than one Pod could be started for the same index due to
> various reasons such as node failures, kubelet restarts or Pod evictions. In
> this case, only the first Pod that compltes successfully will count towards
> the completion count and update status of the Job. The other Pods that are
> running or completed for the same index will be deleted by the Job controller
> once they are detected.

## Handling Pod and Container Failures

A container in a Pod may fail for a number of reasons, such as because the
process in it exited with a non-zero exit code, or the container was killed for
exceeding memory limit, etc. If this happens, and the
`.spec.tempalate.spec.restartPolicy = "OnFailure"`, then the Pod stays on the
node, but the container is re-run. Therefore, the program needs to handle case
when it is restarted locally, or else specify
`.spec.template.spec.restartPolicy = Never`.

An entire Pod can also fail, for number of reasons, such as when the Pod is
kicked off the node (node is upgraded, rebooted, deleted, etc.), or if a
container of the Pod fails and the 
`.spec.template.spec.restartPolicty = "Never"`. When a Pod fails then the Job
controller starts a new Pod. This means that the application needs to handle the
case when it is restarted in a new pod. In particular, it needs to handle
temporary files, locks, incomplete output and the like caused by previous runs.

By default, each pod failure is counted towards the `.spec.backoffLimit` limit.
However, we can customize handling of pod failures by setting the Jobs' pod
failure policy.

When the feature gate `PodDisruptionConditions` and `JobPodFailurePolicy` are
both enabled, and the `spec.podFailurePolicy` field is set, the Job controller
does not consider a terminating Pod (a pod that has a
`.metadata.deletionTimestamp` field set) as a failure until that Pod is terminal
(its `.status.phase` is either `Failed` or `Succeeded`). However, the Job
controller creates a replacement Pod as soon as the termination becomes
apparent. Once the pod terminates, the Job controller evaluates `.backoffLimit`
and `.podFailurePolicy` for the relevant Job, taking this now-terminated Pod
into consideration.

If either of these requirements is not satisfied, the Job controller counts a
temrinating Pod as an immediate failure, even if that Pod later terminates with
`phase: "Succeeded"`.

### Pod Backoff Failure Policy

There are situations where we want to fail a Job after some amount of reties due
to logical error in configuration etc. To do so, set `.spec.backoffLimit` to
specify the number of reties before considering a Job as failed. The back-off
limit is set by default to 6. Failed Pods associated with the Job are recreated
by the Job controller with an exponential back-off delay (10s, 20s, 40s, ...)
capped at six minutes.

The number of retries is calculated in two ways:
- The number of Pods with `.status.phase = "Failed`.
- When using `restartPolicy = "OnFailure"`, the number of reties in all the
  containers of Pods with `.status.phase` equal to `Pending` or `Running`.

If either of the calculations reaches the `.spec.backoffLimit`, the Job is
considered failed.

> [!NOTE]
> If the job has `restartPolicy = "OnFailure"`, keep in mind that the pod
> running the Job will be terminated once the job backoff limit has beenr
> eached. This can make debugging the Job's executable more difficult. Setting
> the `restartPolicy = "Never"` when debugging the Job or using a logging system
> is recommended to ensure output from failed Jobs is not lost inadvertently.

### Backoff Limit per Index

> [!NOTE]
> We can only configure the backoff limit per index for an Indexed Job, if we
> have the `JobBackoffLimitPerIndex` feature gate enabled in your cluster.

When we run an indexed Job, we can choose to handle retries for pod failures
independently for each index. To do so, set the `.spec.backoffLimitPerIndex` to
specify the maximumal number of pod failures per index.

When the per-index backoff limit is exceeded for an index, K8s considers the
index as failed and adds it to the `.status.failedIndexes` field. When the
number of failed indexes exceeds the `maxFailedIndexes` field, the Job
controller triggers termination of all remaining running Pods for the Job. Once
all pods are terminated, the entire Job is marked failed by the Job controller,
by setting the Failed condition in the Job status.

Here is an example manifest for a Job that defines a `backoffLimitPerIndex`:

```yaml
# job-backoff-limit-per-index-example.yaml


apiVersion: batch/v1
kind: Job
metadata:
  name: job-backoff-limit-per-index-example
spec:
  completions: 10
  parallelism: 3
  completionMode: Indexed  # required for the feature
  backoffLimitPerIndex: 1  # maximal number of failures per index
  maxFailedIndexes: 5      # maximal number of failed indexes before terminating the Job execution
  template:
    spec:
      restartPolicy: Never # required for the feature
      containers:
      - name: example
        image: python
        command:           # The jobs fails as there is at least one failed index
                           # (all even indexes fail in here), yet all indexes
                           # are executed as maxFailedIndexes is not exceeded.
        - python3
        - -c
        - |
          import os, sys
          print("Hello world")
          if int(os.environ.get("JOB_COMPLETION_INDEX")) % 2 == 0:
            sys.exit(1)          
```

In the example above, the Job controller allows for one restart for each of the
indexes. When the total number of failed indexes exceeds 5, then the entire Job
is terminated. Once the job is finished, the Job status looks as follows:

```bash
kubectl get -o yaml job job-backoff-limit-per-index-example
```

```bash
status:
    completedIndexes: 1,3,5,7,9
    failedIndexes: 0,2,4,6,8
    succeeded: 5          # 1 succeeded pod for each of 5 succeeded indexes
    failed: 10            # 2 failed pods (1 retry) for each of 5 failed indexes
    conditions:
    - message: Job has failed indexes
      reason: FailedIndexes
      status: "True"
      type: Failed
```

Additionally, we may want to use the per-index backoff along with a pod failure
policy. When using per-index backoff, there is a new `FailIndex` action
available which allows us to avoid unnecessary retries within an index.
