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
