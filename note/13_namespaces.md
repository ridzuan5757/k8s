# Namespaces

`Namespaces` are a way to isolate cluster resources into groups. They are a bit
like directories on our computer, but instead of containing files, they contain
k8s objects. As for the example, every resource in k8s has a name, and some of
it would include:
- synergychat-api-configmap
- api-service
- api-deployment
- web-deployment
- ...

We can only use a name once. It is a unique identifier. That is how `kubectl
apply` knows when it should create a new resource and when it should update an
existing one. Namespaces allow us to use the same name for different resources,
as long as they are in different namespaces.

To check the nampespaces:

```bash
kubectl get namespaces
```

or:

```bash
kubectl get ns
```

# Moving Namespaces

Up until this point, we have been working in the `default` namespace. When using
`kubectl` commands, we can specify the namespace with the `--namespace` or `-n`
flag. If we did not do this, it will use `default` namespace.

```bash
wagslane@MacBook-Pro courses % kubectl get pod
NAME                                  READY   STATUS    RESTARTS   AGE
synergychat-api-646c6fd585-dk5db      1/1     Running   0          28m
synergychat-crawler-cd4947995-tcrkn   3/3     Running   0          39m
synergychat-web-846d86c444-d9c8q      1/1     Running   0          28m
synergychat-web-846d86c444-sk6n4      1/1     Running   0          28m
synergychat-web-846d86c444-w2pqg      1/1     Running   0          28m
```

vs

```bash
wagslane@MacBook-Pro courses % kubectl -n kube-system get pod
NAME                               READY   STATUS    RESTARTS     AGE
coredns-5d78c9869d-jwcbr           1/1     Running   0            4d
etcd-minikube                      1/1     Running   0            4d
kube-apiserver-minikube            1/1     Running   0            4d
kube-controller-manager-minikube   1/1     Running   0            4d
kube-proxy-j2ssm                   1/1     Running   0            4d
kube-scheduler-minikube            1/1     Running   0            4d
storage-provisioner                1/1     Running   1 (4d ago)   4d
```

The `kube-system` namespace is where all the core k8s components live, it is
created automatically when we install k8s.

## Making a new namespace

To create a namespace:

```bash
kubectl create ns <namespace>
```

Verify that is has been created:

```bash
kubectl get ns
```

To move the resources to the namespace, add:

```yaml
metadata:
    namespace: <namespace>
```

section to each of the resources that need to be moved and apply them.
Interestingly, we should see that the resources are "created" instead of
"updated". that is because they are now in a new namespace, and the unique
identifier of a resource in l8s is the combination of its name and its
namespace.

Make sure that the resources are now redeployed in the new namespace:

```bash
kubectl -n <namespace> get pods
kubectl -n <namespace> get svc
kubectl -n <namespace> get configmaps
```

Then go delete the old resource in the `default` namespace:

```bash
kubectl delete deployment <deployment-name>
kubectl delete service <service-name>
kubectl delete configmap <configmap-name>
```

# Intra-cluster DNS

The front-end of synergychat communicates with the `api` application via an
external ingress:

Domain name `http://synchatapi.internal` -> ingress -> service -> pod

It;s now time to connect the `crawler` and `api` applications. The `api` needs
to be able to make HTTP request directly to the `crawler` so that it can get the
latest data to power the "stats" slash command.

`front-end` -> `api` -> `crawler`

The HTTP communication between the `api` and the `crawler` is strictly internal
to the cluster, there is no need for an external domain name or ingress. That
makes it simpler, faster and more secure.

## Slash command

With the tunnel is open (`minikub tunnel -c`) open `http://synchat.internal/` in
our browser. When we are attempting to use `/stat` in the chat, we should see a
response that says:

```
crawler-bot: Crawler worker not configured
```

That is because the `api` does not know how to communicate with the `crawler`
yet.

## DNS

k8s automatically creates DNS entries for each service that can be used to route
HTTP traffic between services. The format is:

```bash
<service_name>.<namespace>.svc.cluster.local
```

#### DSN for Services and Pods

k8s creates DNS records for services and pods allowing us to contact services
with consistent DNS names instead of IP addresses. k8s publishes info about pods
and services which is used to program dns. Kubelet configures pods' DNS so that
running containers can lookup services by name, rather than IP.

Service defined in the cluster are assigned DNS names. By default, a client
pod's DNS search list includes the pod's own namespace and the cluster's default
domain.

#### Namespaces of services

A DNS query may return different results based on the namespace of the pod
making it. DNS queries that do not specify a namespace are limited to the pod's
namespace. Access services in other namespaces by specifying it in the DNS
query.

For example:
Consider a pod in `test` namespace. A `data` service is in the `prod` namespace.
- A query for `data` will returns no result because uses the pod's `test`
  namespace.
- A query for `data.prod` returns the intended result, because it specifies the
  namespace.

DNS queries may be expanded using the pod's `/etc/resolv/config`. Kubelet
configures this file for each pod. For example, a query for just `data` may be
expanded to `data.test.svc.cluster.local`. The values of the `search` option are
used to expand queries.

```conf
### /etc/resolv.conf
nameserver 10.32.0.10
search <namespace>.svc.cluster.local svc.cluster.local cluster.local
options ndots:5
```

In summary, a pod in `test` namespace can successfully resolve either `data.prod` or `data.prod.svc.cluster.local`.

Resource that get DNS records:
- Services
- Pods


