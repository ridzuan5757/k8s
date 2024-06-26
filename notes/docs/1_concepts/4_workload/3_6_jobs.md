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

### Pod Failure Policy

> [!note] We can only configure a Pod failure policy for a Job isf we have the
> `JobPodFailurePolicy` feature gate enabled in the cluster. Additionally, it is
> recommended to enable the `PodDisruptionConditions` feature gate in order to
> be able to detect and handle Pod disruption conditions in the Pod failure
> policy. Both feature gates are available in K8s v.1.30.

A Pod failure policy, defined with the `.spec.podFailurePolicy` field, enables
the cluster to handle Pod failures based on the container exit codes and the Pod
conditions.

In some situations, we may want to have a better control when handling Pod
failures than the control provided by the Pod backoff failure policy, which is
based on the Job's `.spec.backoffLimit`. These are some examples of use cases:
- To optimize costs of running workloads by avoiding unnecessary Pod restarts,
  we can terminate a Job as soon as one of its Pods fails with an exit code
  indicating a software bug.
- To guarantee that the Job finishes even if there are disruptions, we can
  ignore Pod failures caused by disruptions such as preemption, API-initiated
  eviction or taint-based eviction so that they do not count towards the
  `.spec.backoffLimit` limit of retries.

We can configure a Pod failure policy, in the `.spec.podFailurePolicy` field, to
meet the above use cases. This policy can handle Pod failures based on the
container exit codes and the Pod conditions.

Here is a manifest for a Job that defines a `podFailurePolicy`:

```yaml
# job-pod-failure-policy-example.yaml

apiVersion: batch/v1
kind: Job
metadata:
    name: job-pod-failure-policy-example
spec:
    completions: 12
    parallelism: 3
    template:
        spec:
            restartPolicy: Never
            containers:
            - name: main
              image: docker.io/library/bash:5
              # example command simulating a bug which triggers the failjob
              # action
              command: ["bash"]
              args:
              - -c
              - echo "Hello world!" && sleep 5 && exit 42
    backoffLimit: 6
    podFailurePolicy:
        rules:
        - action: FailJob
          onExitCodes:
            # this container name is optional
            containerName: main
            operator: In    # one of In, NotIn
            values: [42]
        - action: Ignore    # one of Ignore, FailJob, Count
          onPodConditions:
          # indicates pod disruption
          - type: DisruptionTarget
```

In the example above, the first rule of the Pod failure policy specifies that
the Job should be marked failed if the `main` container fails with the 42 exit
code. The following are rules for the `main` container specifically:
- An exit code of `0` means that the container succeeded.
- An exit code of `42` means that the entire job failed.
- Any other exit code represents that the container failed, and hence the entire
  Pod. The Pod will be re-created if the total number of restarts is below
  `backoffLimit`. If the `backoffLimit` is reached the **entire Job** failed.

> [!NOTE]
> Because the Pod template specifies a `restartPolicy: Never`, the kubelet does
> not restart the `main` container in that particular Pod.

The second rule of the Pod failure policy, specifying the `Ignore` action for
failed Pods with condition `DisruptionTarget` excludes Pod disruptions from
being counted towards the `.spec.backoffLimit` limit of retries.

> [!NOTE]
> If the Job failed, either by the Pod failure policy or Pod backoff failure
> policy, and the Job is running multiple Pods, K8s terminates all the Pods in
> that Job that are still Pending or Running.

These are some requirements and semantics of the API:
- If we want to use a `.spec.podFailurePolicy` field for a Job, we must also
  define that Job is pod template with `.spec.restartPolicy` set to `Never`.
- The Pod failure policy rules that we have to specify under
  `spec.podFailurePolicy.rules` are evaluated in order. Once a rule matches a
  Pod failure, the remaining rules are ignored. When no rule matches the Pod
  failure, the default handling applies.
- We may want to restrict a rule to a specific container by specifying its name
  in `spec.podFailurePolicy.rules[*].onExitCodes.containerName`. When not
  specified the rule applies to all containers. When specified, it should match
  one the container or `initContainer` names in the Pod template.
- We may specify the action taken when a Pod failure policy is matched by
  `spec.podFailurePolicy.rules[*].action`. Possible values are:
    - `FailJob` use to indicate that the Pod's job should be marked as Failed
      and all running Pods should be terminated.
    - `Ignore` use to indicate that the counter towards the `.spec.backoffLimit`
      should not be incremented and a replacement Pod should be created.
    - `Count` use to indicate that the Pod should be handled in the default way.
      The counter towards the `.spec.backoffLimit` should be incremented.
    - `FailIndex` use this action along with backoff limit per index to avoid
      unnecessary retries within the index of a failed pod.

> [!NOTE]
> When we use a `podFailurePolicy`, the job controller only matches Pods in the
> `Failed` phase. Pods with a deletion timestamp that are not in a terminal
> phase `Failed` or `Succeeded` are considered still terminating. this implies
> that terminating pods retain a tracking finalizer until they reach a terminal
> phase. Kubelet transitions deleted pods to a terminal phase. This ensures that
> deleted pods have their finalizers removed by the Job controller.

> [!NOTE]
> When Pod failure policy is used, the Job controller recreates terminating Pods
> only once these Pods reach the terminal `Failed` phase. This behaviour is
> similar to `podReplacementPolicy: Failed`.

## Success Policy

> [!NOTE]
> We can only configure a success policy for an Indexed Job if we have the
> JobSuccessPolicy feature gate enabled in the cluster.

When creating an Indexed Job, we can define when a Job can be declared as
succeeded unsing a `.spec.successPolicy`, based on the pods that succeeded.

By default, a Job succeeds when the number of succeeded Pods equals
`.spec.completions`. These are some situations where we might want additional
control for declaring a Job succeeded:
- When running simulations with different parameters, we might not need all
  simulations to succeed for the overall Job to be successful.
- When following a leader-worker pattern, only the success of the leader
  determines the success or failure of a Job. Examples of this are frameworks
  like MPI and PyTorch.

We can configure a success policy, in the `.spec.successPolicy` field, to meet
the above use cases. This policy can handle Job success based on the succeeded
pods. After the Job meets the success policy, the Job controller terminates the
lingering Pods. A success policy is defined by rules. Each rule can take one of
the following forms:
- When we specify the `succeededIndexes` only, once all indexes specified in the
  `succeededIndexes` succeed, the job controller makrs the job as succeeded.
  The `succeededIndexes` must be a list of intervals between 0 and
  `.spec.completions-1`.
- When we specify the `succeededCount` only, once the number of succeeded
  indexes reaches the `succeededCount`, the Job controller marks the Job as
  succeeded. 
- When we specify both `succeededIndexes` and `succeededCount`, once the number
  of succeeded indexes from the subset of indexes specified in the
  `succeededIndexes` reaches the `succeededCount`, the Job controller marks the
  Job as succeded.

Note that when we specify multiple rules in the `.spec.successPolicy.rules`, the
Job controller evaluates the rules in order. Once the Job meets a rule, the Job
controller ignores remaining rules. Here is a manifest for a Job with
`successPolicy`:

```yaml
# job-success-policy.yaml

apiVersion: batch/v1
kind: Job
spec:
    parallelism: 10
    completions: 10
    completionMode: Indexed     # required for the success policy
    successPolicy:
        rules:
            - succeededIndexes: 0, 2-3
              succeededCount: 1
    template:
        spec:
            containers:
            - name: main
              image: python
              command:
              # provided that at least one of the Pods with 0, 2 and 3 indexes
              # has succeded, the overall job count as successful
                - python 3
                - -c
                - |
                    import os, sys
                    if os.environ.get("JOB_COMPLETION_INDEX") == "2":
                        sys.exit(0)
                    else:
                        sys.exit(1)
```

In the example above, both `succeededIndexes` and `succeededCount` have been
specified. Therefore, the job controller will mark the Job as succeeded and
terminate the lingering Pods when either of the specified indexes, 0, 2, or 3,
succeed. The Job that meets the success policy gets the `SuccessCriteriaMet`
condition. After the removal of the lingering Pods is issued, the Job gets the
`Complete` condition.

Note that the `succeededIndexes` is represented as intervals separated by a
hyphen. The number are liested in represented by the first and last element of
the series, separated by hyphen.

> [!NOTE]
> When we specify both a success policy and some terminating policies such as
> `.spec.backoffLimit` and `.spec.podFailurePolicy`, once the Job meets either
> policy, the Job controller respects the terminating policy and ignores the
> success policy.

## Job termination and cleanup

When a Job completes, no more Pods are created, but the Pods are usually not
deleted either. Keeping them around allows us to still view the logs of
completed pods to check for errors, warnings, or other diagnostic output. The
job object also remains after it is compelted so that we can view its status. It
is up to the user to delete old jobs after noting their status. Delete the job
with `kubectl` (`kubectl delete jobs/pi` or `kubectl delete -f job.yaml`). When
we delete job using `kubectl`, all the pods it created are deleted too.

By default, a Job will run uninterrupted unless a Pod fails 
`restartPolicy=Never` or a container exits in error `restartPolicy=OnFailure`,
at which point the Job defers to the `.spec.backoffLimit` described above. Once
`.spec.backoffLimit` has been reached the Job will be marked as failed and any
runnign Pods will be terminated.

Another way to terminate a Job is by setting an active deadline. Do this by
setting `.spec.activeDeadlineSeconds` field of the Job to a number of seconds.
The `activeDeadlineSeconds` applies to the duration of the job, no matter how
many pods are created. Once a Job reaches `activeDeadlineSeconds`, all of its
running Pods are terminated and the Job status will become `type: Failed` with
`reason: DeadlineExceeded`.

Note the a Job's `.spec.activeDeadlineSeconds` takes precendence over its
`.spec.backoffLimit`. Therefore, aJob that is retrying one or more failed Pods
will not deploy additional Pods once it reaches the time limit specified by
`activeDeadlineSeconds`, even if the `backoffLimit` is not yet reached. For
example:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
    name: pi-with-timeout
spec:
    backoffLimit: 5
    activeDeadlineSeconds: 100
    template:
        spec:
            containers:
            - name: pi
              image: perl:5.34.0
              command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
            restartPolicy: Never
```

Note that both the Job spec and the Pod template spec within the Job have an
`activeDeadlineSeconds` field. Ensure that we set this field at the proper
level.

Keep in mind that the `restartPolicy` applies to the Pod, and not to the Job
itself. There is no automatic Job restart once the Job status is `type: Failed`.
That is, the Job termination mechanisms activated with 
`.spec.activeDeadlineSeconds` and `.spec.backoffLimit` result in a permanent Job
failure that requires manual intervention to resolve.

## Clean up finished jobs automatically

Finished Jobs are usually no longer needed in the system. Keeping them around in
the system will put pressure on the API server. If the Jobs are managed directly
by a higher level controller, such as CronJobs, the Jobs can be cleaned up by
CronJobs based on the specified capacity-based cleanup policy.

### Time To Live TTL mechanism for finished Jobs

Another way to clean up finished Jobs either `Complete` or `Failed`
automatically is to use a TTL mechanism provided by a TTL controller for
finished resources, by specifying the `.spec.ttlSecondsAfterFinished` field of
the Job.

When the TTL controller cleans up the Job, it will delete the Job cascadingly,
i.e. delete its dependent objects, such as Pods, together with the Job. Note
that when the Job is deleted, its lifecycle guarantees, such as finalizers, will
be honored. For example:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi-with-ttl
spec:
  ttlSecondsAfterFinished: 100
  template:
    spec:
      containers:
      - name: pi
        image: perl:5.34.0
        command: ["perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
```

The Job `pi-with-ttl` will be eligible to be automatically deleted, 100 seconds
after it finishes. If the field is set to 0, the Job will be eligible to be
automatically deleted immediately after it finishes. If the field is unset, this
Job would not be cleaned up by the TTL controller after it fisnishes.

> [!NOTE]
> It is recommended to set `ttlSecondsAfterFinished` field because unmanaged
> jobs (Jobs that we created directly, and not indirectly through other workload
> APIs such as CronJob) have a default deletion policy of `orphanDependents`
> causing Pods created by an unmanaged Job to be left around after that Job is
> fully deleted.
>
> Even though the control plane eventually garbage collects the Pods from a
> deleted Job after they either fail or complete, sometimes those lingering pods
> may cause cluster performance degradation or in worst case cause the cluster
> to go offline due to this degradation.
>
> We can use LimitRanges and ResourceQuotas to place a cap on the amount of
> resources that a particular namespace can consume.

## Job Patterns

The Job object can be used to process a set of independent but related work
items. These might be emails to be sent, frames to be rendered, files to be
transcoded, ranges of keys in a NoSQL database to scan and so on.

In a complex system, there may be multiple different sets of work items. Here we
are just considering one set of work items that the user wants to manage
together - a batch job.

There are several different patterns for parallel computation, each with
strengths and weaknesses. The tradeoffs are:
- One Job object for each work item, versus a single Job object for all work
  items. One Job per work item creates some overhead for the user and for the
  system to manage large numbers of Job objects. A single Job for all work items
  is better for large numbers of items.
- Number of Pods created equals number of work items, versus each Pod can
  process multiple work items. When the number of Pods equals the number of work
  item, the Pods typically requires less modification to existing code and
  containers. Having each Pod process multiple work items is better for large
  numbers of items.
- Several approaches use a work queue. This requires running a queue service,
  and modifications to the existing program or container to make it use the work
  queue. Other approaches are easier to adapt to an existing containerised
  application.
- When the Job is associated with a headless Service, we can enable the Pods
  within a Job to communicate with each other to collaborate in a computation.

## Advanced Usage

### Suspending a Job

When a Job is created, the Job controller will immediately begin creating Pods
to satisfy the Job's requirements and will continue to do so until the Job is
complete. However, we may want to temporarily suspend a Job's execution and
resume it later, or starts Jobs in suspended state and have a custom controller
decide later when to start them.

To suspend a Job, we can update the `.spec.suspend` field of the Job to true;
later, when we want to resume it again, update it to false. Creating aJob with
`.spec.suspend` set to true will create it in the suspended state.

When a Job is resumed from suspension, its `.status.startTime` field will be
reset to the current time. This means that the `.spec.activeDeadlineSeconds`
timer will be stopped and reset when a Job is suspended and resumed.

When we suspend a Job, any running Pods that do not have a status of `Completed`
will be terminated with a SIGTERM signal. The Pod's graceful termination period
will be honored and the Pod must handle this signal in this period. This may
involve saving progress for later or undoing changes. Pods terminated this way
will not count towards the Job's `completions` count. An example Job definition
in the suspended state can be like so:

```bash
kubectl get job myjob -o yaml
```

```yaml
apiVersion: batch/v1
kind: Job
metadata:
    name: myjob
spec:
    suspend: true
    parallelism: 1
    completions: 5
    template:
        spec:
```

We can also toggle Job suspension by patching th eJob using the command line. To
suspend an active Job:

```bash
kubectl patch jpb/myjob --type=strategic --patch '{"spec":{"suspend":true}}'
```
Resume a suspended Job:

```bash
kubectl patch job/myjob --type=strategic --patch '{"spec":{"suspend":false}}'
```

The Job status can be used to determine if a Job is suspended or has been
suspended in the past:

```bash
kubectl get jobs/myjobs -o yaml
```

```yaml
apiVersion: batch/v1
kind: Job
# .metadata and .spec omitted
status:
  conditions:
  - lastProbeTime: "2021-02-05T13:14:33Z"
    lastTransitionTime: "2021-02-05T13:14:33Z"
    status: "True"
    type: Suspended
  startTime: "2021-02-05T13:13:48Z"
```

The Job condition of type "Suspended" with status "True" means the Job is
suspended; the `lastTransitionTime` field can be used to determine how long th
eJob has been suspended for. If the status of that condition is "False", then
the Job was previously suspended and is now running. If such a condition does
not exist in the Job's status, the Job has never been stopped. Events are also
created when the Job is suspended and resumed:

```bash
kubectl describe job/myjob
```

```bash
Name:           myjob
...
Events:
  Type    Reason            Age   From            Message
  ----    ------            ----  ----            -------
  Normal  SuccessfulCreate  12m   job-controller  Created pod: myjob-hlrpl
  Normal  SuccessfulDelete  11m   job-controller  Deleted pod: myjob-hlrpl
  Normal  Suspended         11m   job-controller  Job suspended
  Normal  SuccessfulCreate  3s    job-controller  Created pod: myjob-jvb44
  Normal  Resumed           3s    job-controller  Job resumed
```

The last four events, particularly the "Suspended" and "Resumed" events, are
directly a result of toggling the `.spec.suspend` field. In the time between
these two events, we see that no Pods were created, but Pod creation restarted
as soon as the Job was resumed.

## Mutable Scheduling Directives

In most cases, a parallel job will want the pods to run with constraints, like
all in the same zone, or all either on GPU model x or y but not a mix of both.

The suspend field is the first step towards achieving those semantics. Suspend
allows a custom queue controller to decide when a job should start. However,
once a job is unsuspended, a custom queue controller has no influence on where
the pods of a job will actually land.

This feature allows updating a Job's scheduling directives before it starts,
which gives custom queue controllers the ability to influence pod placement
while at the same time offloading actual pod-to-node assignbment to
kube-scheduler. This is allowed only for suspended Jobs that have never been
unsuspended before.

The fields in a Job's pod template that can be updated are node affinity, node
selector, tolerations, labels, annotations and shceduling gates.


