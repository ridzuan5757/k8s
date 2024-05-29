# Service

Service is a method for exposing a network application that is running as one or
more Pods in the cluster.

The aim of Services in k8s is that we do not need to modify the existing
application to use an unfamiliar service discovery mechanism. We can run code in
Pods, whether this is a code designed for a cloud native world, or an older app
we have containerized. We use a Service to make that set of Pods available on
the network so that clients can interact with it.

If we use a Deployment to run the app, that Deployment can create and destroy
Pods dynamically. From one moment to the next, we do not know how many of those
Pods are working and healthy, we might not even know what those healthy Pods are
named K8s Pods are created and destroyed to match the desired state of the
cluster. Pods are ephemeral resources as we should not expect that an individual
Pod is reliable and durable.

This leads to a problem:
If some set of Pods called "backends" provides functionality to other Pods
called "frontends" inside the cluster, how do the frontends find out and keep
track of which IP address to connect to, so that the frontend can use the
backend part of the worklod?

## Services in k8s


