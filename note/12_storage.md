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
