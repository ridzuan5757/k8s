# Service

In Kubernetes, a Service is a method for exposing a network application that is
running as one or more Pods in the cluster.

A key aim of SErvices in Kubernetes is that we do not need to modify existing
application to use an unfamiliar service discovery mechanism. We can run code in
Pods, whether this is a code designed for a cloud-native world, or an older app
we have containerized. We use Service to make that set of Pods available on the
network so that clients can interact with it.

If we use `Deployment` to run the app, that Deployment can create and destroy
Pods dynamically.. From one moment to the next, we do not know how many of those
Pods are working and healthy. We might not even know that whose healthy Pods are
named. 
