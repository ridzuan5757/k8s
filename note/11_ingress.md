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
