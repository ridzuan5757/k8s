# Deployment
Once we have a running k8s cluster, we can deploy containerized apps on top of
it. To do so, we create a k8s deployment. Once we have created a deployment,
the k8s control plane schedules the application instances included in that
deployment to run on individual nodes in that cluster.

Once the application instances are created, a k8s deployment controller
continuously monitor those instances. Id the node hosting ann instances go down
or is deleted, the deployment controller replaces the instance with an instance
on another node in the cluster. This provides a self-healing mechanism to
address machine failure or maintenance.

# Deploying app on k8s
We can create a manage a deployment by using k8s cli `kubectl`. `kubectl` use
Kubernetes API to interact with the cluster. 

When we create a deployment, we will need to specify the container image for our
application and the number of replicas that we want to run. We can change that
information later by updating the deployment.

For this example, we will use a `hello-node` application packaged in Docker
container that use NGINX to echo back all of the request.

## Creating deployment
Let's deploy the app on `k8s` using `kubectl create deployment`. We need to
provide the deployment name and app image location.

```bash
kubectl create deployment kubernetes-bootcamp --image=gcr.io/google-samples/kubernetes-bootcamp:v1
```

This command:
- searched for a suitable node where an instance of the application could be
  run. (We have only 1 available node.)
- scheduled the application to run on that node.
- configured the cluster to reshedule the instance on a new node when needed.

To list our deployments:
```bash
kubectl get deployments
```

We see there is 1 deployment running a single instance of our app. The instance
is running inside a container in our node.

## Viewing the app
`pods` that are running inside k8s are running on a private isolated network. By
default they are visible from other pods and service within the same k8s
cluster, but not outside the network. When we use `kubectl`, we are interacting
through an API endpoint to communicate with our application.

The `kubectl proxy` command can create a proxy that forward communications into
the cluster-wide private network.

### `kubectl proxy` vs `minikube expose`
These are two different approaches to make services running inside a k8s cluster
accessible from outside.

`kubectl proxy`:
- Creates a proxy server between our local machine and k8s API server.
- Allows to access k8s services via HTTP.
- Proxy creates a secure tunnel from local machine to the k8s cluster, so we can
  access k8s resources without exposing them directly to the internet.
  where we need services to be accessible externally.

`minikube expose`
- Tools that lets we run k8s locally.
- When we expose service in minikube, we are making it accessible from outside
  the minikube cluster.
- Provides various ways to expose services such as NodePort, LoadBalancer and
  Ingress.
- NodePort exposes a service on a port on each node in the cluster.
- LoadBalancer provisions external load balancer in environments that suppors
  it.
- Ingress exposes HTTP and HTTPS routes from outside of the clouster to services
  within the cluster.

