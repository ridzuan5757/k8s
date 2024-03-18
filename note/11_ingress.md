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

We will declare that:
- Any traffic to the `synchat.internal` domain name should be routed to the
  `web-service`.
- Any traffic to `synchatapi.internal` domain name should be routed to the
  `api-service`.


