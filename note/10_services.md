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
