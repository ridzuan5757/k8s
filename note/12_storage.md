# Storage in k8s

By default, containers running in pods on k8s have access to the filesystem, but
there are some big limitations to this. Even though we are saving the file to
the filesystem, it will not persists once the pod destroyed and recreated. In
other word, the filesystem is ephemeral as the pods.

This has to do with the philosophy behind k8s and even containers in general:
when we spin up a new one, it should always be a blank state, which makes
reproducing and debugging much easier since we don't have to maintain the state
consistency.
