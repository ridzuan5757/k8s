# Creating minikube cluster

```bash
  minikube start
```

# Dashboard

```bash
minikube dashboard --url
```

# Creating a Deployment
`Pod` is a group of one or more `Containers`, tied together for purposes of
administration and networking. For this case, this `pod` has only one 
`container`. A k8s deployment checks on the health of the $pod$ and restarts
the $pod$'s container if it terminates. This is the recommended way to manage
creation and scaling of `pods`.
- `kubectl create` is used to create a deployment that manages a `pod`. The 
`pod` runs a `container` based on the provided `docker` image.

```bash
kubectl create deployment hello-node --image=registry.k8s.io/e2e-test-images/agnhost:2.39 -- /agnhost netexec --http-port=8080
```

- `kubectl create deployment hello-node` creates a new deployment named 
  "hello-node".
- `--image=registry.k8s.io/e2e-test-images/agnhost:2.39` specifies the docker
  image to use for the pods managed by the deployment. In this case, the image
  with tag "2.39" is pulled from the repository.
- `-- /agnhost netexec --http-port=8080` specifies the command to run inside the
  container. In this case, it is running the `netexec` from the `agnhost` image
  with argument of `http-port=8080`. This command is used to start simple HTTP
  server listening on port 8080.

- To view the deployment:
```bash
kubectl get deployments
```

- To view the pod:
```bash
kubectl get pods
```

- To view the cluster events
```bash
kubectl get events
```

- To view the `kubectl` configuration:
```bash
kubectl config view
```

- To view the application logs for container in a pod:
```bash
kubectl logs <pod_name>
```


# Creating Service
By default, a `pod` is only accessible by its internal IP address within the
`k8s` cluster. To make the `hello-node` container accessible from outside of
`k8s` vpn, we have to expose the `pod` as a `k8s` service.
- Use `kubectl expose` to expose the `pod` to public internet:
```bash
kubectl expose deployment hello-node --type=LoadBalancer --port=8080
```
    
- The `--type=LoadBalancer` flag indicates that we want to expose our service 
outside of the cluster. The application code inside the test image only listens
on TCP port 8080. If we used `kubectl expose` to different port, clients could 
not connect to other port.

- To view the created service:
```bash
kubectl get services
```

On cloud providers that support load balancers, an external IP address would be 
provisioned to access the service. On `minikube`, the `LoadBalancer` type makes 
the service accessible through the `minikube service` command.

- Access the application:
```bash
minikube service hello-node
```

# Cleaning up
```bash
kubectl delete service hello-node
kubectl delete deployment hello-node
minikube stop
minikube delete
```



