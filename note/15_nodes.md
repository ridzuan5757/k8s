# Nodes

We have talked about in a production environment, we will have multiple nodes in
our cluster and we have been using single node cluster with minikube. The nice
thing about k8s is that almost everything that we do is abstracted away from the
underlying infrastructure with the `kubectl` CLI.

## Deploying to production

### GKE, EKS, AKS
These are all managed k8s services, offered by cloud providers. They are all
pretty similar. GKE appear to be the most feature-rich out of the thress. GKE
also has auto-pilot mode that makes it so that we do not have to worry about
managing nodes at all.

The nice thing about a managed offering is that it can be configured to handle
autoscaling at the node level. This means that we can set up our cluster
automatically add and remove nodes based on the load on the cluster.

### Manual

We can also set up our own cluster manually. One of the way is by having custom
scripts that configure a cluster on top of standard EC2 instances. Then have our
own autoscaling scripts that add and remove nodes based on the load of the
cluster. It is also pissble to do the same thing on physical machines.
