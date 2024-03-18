# Ingress

An `ingress` resource exposes services to the outside world and is used often in
production environments.

> An Ingress may be configured to give Services externally-reachable URLs, load
> balance traffic, terminate SSL/TLS, and offer name-based virtual hosting. An
> Ingress controller is responsible for fulfilling the Ingress, usually with a
> load balancer, though it may also configure our edge router or additional
> frontends to help handle the traffic.
>
> An Ingres does not expose arbirary ports or protocols. Exposing services other
> than HTTP and HTTPS to the internet typically uses a service of type NodePort
> or LoadBalancer.

Think of Ingress as a load balancer that lives outside the cluster and routes
traffic through the ingress to a service.

## Setting up ingress

To work with an ingress first we need to enable it in minikube:

```bash
minikube addons enable ingress
```

Next, create a new file. We will call it `app-ingress.yaml` because it will be
an ingress for the entire synergychat application, not just a specific service.

## Ingress metadata

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
    annotations:
        nginx.ingress.kubernetes.io/rewrite-target: /
```

The `annotations` section is where we can add extra configuration for our
ingress. In this case, we are telling it to rewrite the target URL to `/` so
that it will work with our web app.


## Rules

`spec/rules` section is where we define the routing rules for our ingress. We 
will declare that:
- Any traffic to the `synchat.internal` domain name should be routed to the
  `web-service`.
- Any traffic to `synchatapi.internal` domain name should be routed to the
  `api-service`.

```yaml
spec:
  rules:
    - host: synchat.internal
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: web-service
                port:
                  number: 80
    - host: synchatapi.internal
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-service
                port:
                  number: 80
```

This says that any traffic to the `synchat.internal` domain should be routed to
the `web-service` and any traffic to `synchatapi.internal` should be routed to
`api-service`.

# DNS
Now that we have configured the ingress to route the domains:
- `synchat.internal` to the `web-service`
- `synchatapi.internal` to the `api-service`

We need to configure our local machine to resolve those domains to the ingress
load balancer. We would not be setting up global DNS so that anyone on the
internet can access our app. We will just be configuring our local machine to
resolve those domains to the ingress load balancer.

There is a file called `/etc/hosts` on our local machine that is used to resolve
domain names to IP addresses. We can add entries to that file to resolve our
domains to the ingress load balancer.

```conf
127.0.0.1   synchat.internal
127.0.0.1   synchatapi.internal
```

#### WSL DNS configuration

For WSL users, we also need to add the entries above to the Windows host file
located in:

```bash
C:\Windows\System32\drivers\etc\hosts
```

The WSL oath to this file is:

```bash
/mnt/c/Windows/System32/drivers/etc/hosts
```

To verify that it is working:

```bash
ping synchat.internal
```

# Tunnel

In production, once we have an ingress configured and have pointed our domain
name to it(and perhaps its load balancer), we can access our application from
anywhere in the world. The trouble is, we are using Minikube and the cluster is
running on our local machine, and not only that, it is running in an isolated
virtual machine. Minikube however, capable to forward the ingress to our local
machine.

To open a tunnel to our cluster:

```bash
minikube tunnel -c
```

The tunnel should expose the ingress controller's load balancer to our local
machine on `localhost:80`, which we mapped to the `synchat.internal` and
`synchatapi.internal`.

While the tunnel is open, both `http://synchat.internal` and
`http://synchatapi.internal` is accessible via browser and will return 200.

# Ingress Types

On top of all of our resources, we have the following information:

```yaml
apiVersion: v1
```

This is the API version of the resource, and because those resources are core to
k8s, they are in the standard `v1` API group. However, ingress is not a core k8s
resource, it is an extension of sorts, this is why the configuration for
`apiVersion` is different:

```yaml
apiVersion: networking.k8s.io/v1
```

We can think of the `networking.k8s.io` aPI grouo as a core extension. It is not
third-party, but it is not part of the core k8s API either.


