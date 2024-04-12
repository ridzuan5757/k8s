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


