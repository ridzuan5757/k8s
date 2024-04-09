# Container Lifecycle Hooks

Analogous to many programming language frameworks that have component lifecycle
hooks such as Angular, k8s provides contaienrs with lifecycle hooks. The hooks
enable contaienrs to be aware of events in their management lifecycle and run
code implemented in the handler when the corresponding lifecycle hook is
executed.

## Container hooks

There are 2 hooks that are exposed to Contaienrs:

### `PostStart`
This hook is executed immediately after a container is created. However, there
is no guarantee that the hook will execute before the container `ENTRYPOINT`. No
parameters are passed to the handler.

### `PreStop`
This hook is called immediately before a container is terminated due to an API
request or management event such as a liveness/startup probe failure,
preemption, resource contention and others.

A call to the `PreStop` hook fails if the container is already in terminated or
complated state and the hook must complete before the TERM signal to stop the
container can be sent.

The pod's termination grace period countdown begins before the `PreStop` hook is
executed, so regardless of the outcode of the handler, the container will
eventually terminate within the pod's termination grace period. No paramters are
passed to the handler.

## Hook handler implementation

Containers can access a hook by implementing and registering a handler for that
hook. There are 3 types of hook handlers that can be implemented for Containers:
- Exec - execues a specific command, such as `pre-stop.sh`, inside the cgroups
  and namespace of the cotnainer. Resources consumed by the command are counted
  agains the Container.
- HTTP - executes an HTTP request against a specific endpoint on the container.
- Sleep - pauses the container for a specified duraiton. The sleep action is
  available when the feature gate `PodLifecycleSleepAction` is enabled.

## Hook handler executions

When a container lifecycle management hook is called, the k8s management system
secutes the handler according to the hook action, `httpGet`, `tcpSocket` and
`sleep` are executed by the kubelet process, and `exec` is executed in the
container.

Hook handler calls are synchronous within the context of the pod containing the
container. This means that for a `PostStart` hook, the container `ENTRYPOINT`
and hook fire asynchronously. However, if the hook takes too long to run or
hangs, the Container cannot reach a `running` state.

`PreStop` hooks are not executed asynchronously from the signal to stop the
Container, the hook must complete its execution before the TERM signal can be
sent. If a `PreStop` hook hangs during execution, the pod's phase will be
`Terminating` and remain there until the pod is killed after its
`terminationGracePeriodSeconds` expires.

This grace period applies to the totak time it takes for both the `PreStop` hook
to execute and for the container to stop normally. If, for example 
`terminationGracePeriodSeconds` is 60, and the hook takes 55 seconds to
complate, and the Container takes 10 seconds to stop normally after receiving
the signal, then the container will be killed before it can stop normally since
`terminationGracePeriodSeconds` is less than the total time it takes for these 2
things to happen.

If either a `PostStart` or `PreStop` hook failed, it kills the container.

Users should make their hook handle as lightweight as possible. There are cases
however, when long running commands make sense, such as when saving state prior
stopping a container.
