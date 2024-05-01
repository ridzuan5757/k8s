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
