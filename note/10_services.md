# Services

We have spun up pods and connected to them individually, but that is frankly not
super usefil if we want to distribute real traffic across those pods. That is
where the services come in.

`Services` provide stable endpoint for pods. They are an abstraction used to
provide a stable endpoint and load balance traffic across a groups of pods. The
service will always be available at a given URL, even if the pod is destroyed
and created.

## Creating a service

Parameters:
- `apiVersion: v1`
- `kind: service`
- `metadata/name: web-service` We could call it anything.
- `spec/selector/app: synergychat-web` this is how the service knows which pods
  to route traffic to.
- `spec/ports` An array of port objects. Minimum of one entry is required.
  - `protocol: TCP`
  - `port: 80` This is the port that the service will listen to.
  - `targetPort: 8080` This is the port that the pod are listening on. 


```yaml
apiVersion: v1
kind: Service
metadata:
    name: web-service
spec:
    selector:
        app: synergychat-web
    ports:
        - protocol: TCP
          port: 80
          targetPort: 8080
```

This creates a new service called `web-service` with a few properties:
- It listen on port 80 for incoming traffic.
- It forwards that traffic to pods on listening on their port 8080.
- Its controller will continuously scan for pods mathcing the `app:
  synergychat-web` label selector and automatically add them to its pool.


To create the service:


```bash
kubectl apply -f web-service.yaml
```

To port forwards the service's port to our local machine so we can test it out.


```bash
kubectl port-forward service/web-service 8080:80
```

Now the service should be accessible via `http://localhost:8080` and it is
better this time around because now our requests are being load-balanced across
3 pods.

# Service Types

To view the YAML file that describe `web-service`:

```bash
kubectl get svc web-service -o yaml
```

We should see a section that looks like this:

```yaml
spec:
    clusterIP: 10.96.213.234
    ...
    type: ClusterIP
```

We did not specify a service type however `ClusterIP` type is being assigned
here since it is the default service type. The `ClusterIP` is the IP address
that the service is bound to on the internal k8s network. Similarly as pods
having their own internal virtual IP address, the service also have their own
internal virtual IP address.

There are also other type of services:
- `NodePort`: Exposes the services on each node's IP at a static port.
- `LoadBalancer`: Creates an external load balancer in the current cloud
  environment (if supported, such as AWS, GCP, Azure) and assins a fixed,
  external IP to the service.
- `ExternalName`: Maps the service to the contents of the `externalName` field 
  (for example, to the hostname `api.foo.bar.example`). The mapping configures 
  our cluster's DNS server to return a `CNAME` record with that external 
  hostname value. No proxying of any kind is set uo.

The interesting thing about service types is that they typically build on top of
each other. 

#### `NodePort` 
This service is just a `ClusterIP` service with the added functionality of 
exposing the service on each node's IP at a static port (it still has an 
internal cluster IP).

#### `LoadBalancer`
This service is just a `NodePort` service with the added functionality of
creating an external load balancer in the current cloud environment (it still
has an internal cluster IP and node port).

#### `ExternalName`
This service functions as DNS-level redirect. We can use it to redirect traffic
from one service to another.

To identify which service that should be used, there are lot of things that need
to be considered. If we are working in a microservices environment where many
services are only meant to be accessed within the cluster, then `ClusterIP` is
going to be it. `NodePort` and `LoadBalancer` are used when we want to expose
the service to the outside world. `ExternalName` is primarily for DNS redirect.

## Ingress Service
`NodePort` and `LoadBalancer` services are used to expose services to the
outside world. However, in most cloud-based k8s environment, we will actually
use `Ingress` object to expose our services. The ingress object not only exposes
our service to the outside world, but also allows us to do thing like:
- Host multiple services on the same IP address
- Host multiple services on the same port (path-based routing)
- Terminate SSL
- Integrate directly with external DNS and load balancers


