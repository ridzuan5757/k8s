# Ephemeral Containers

Ephemeral containers is a special type of containers that runs temporarily in an
existing pod to accomplish user-initiated actions such as troubleshooting. We
can use ephemeral containers to inspect services rather than to build
applications.

## Concepts

Pods are the fundamental building block of k8s applications. Siince pods are
intended to be disposable and replaceable, we cannot add a container to a pod
once it has been created. Instead, we usually delete and replace pods in a
controlled fashion using workload resource deployments.

Sometimes it is necessary to inspect the state of existing pod, however, for
example to troubleshoot a hard-to-reproduce bug. In these case we can run an
ephemeral container in an existing pod to inspect its state and run arbitrary
commands.

Ephemeral containers differ from other containers in that they lack guarantees
for resources or execution and they will never be automatically restarted, so
they are not appropriate for building applications. Ephemeral containers are
described using the same `ContainerSpec` as regular containers, but many fields
are incompatible and disallowed for ephemeral containers.
- Ephemeral containers may not have ports, so these fields are disallowed:
    - `ports`
    - `livenessProbe`
    - `readinessProbe`
- Pod resource allocation are immutable, so the following setting is disallowed:
    - `resources`

Ephemeral containters are created using a special `ephemeralcontainers` handler
in the API rather than by addng them directly to `.pod.spec`, so it is not
possible to add an ephemeral container using `kubectl edit`.

Like a regular containers, we may not change or remove an ephemeral container
after it has been added to a pod.

> Ephemeral containers are not supported by static pods.

## Use cases

Ephemeral containers are useful for interactive troubleshooting when `kubectl
exec` is insufficient because a container has crashed or a container image does
not include debugging utilities.

In particular, distroless images enable us to deploy minimal container images
that reduce attack surface and exposure to bugs and vulnerabilities. Since
distroless images do not include a shell or any debugging utilities, it is
difficult to troubleshoot distroless images using `kubectl exec` alone.

When using ephemeral containers, it is helpful to enable process namespace
sharing so we can view processes in other containers.
