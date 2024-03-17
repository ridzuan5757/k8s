# Minikube

`Minikube` is a tool that allows us to run single-node k8s cluster on local
machine. In production, we would use cluster of servers, probably in the cloud.

# Creating minikube cluster / Running a Minikube

We will be using Kubernetes with Docker, which is the most common way to use
k8s. Ensure Docker daemon is running before starting `Minikube`.

```bash
  minikube start --extra-config
  "apiserver.cors-allowed-origins=["http://boot.dev"]"
```

The extra configuration is just so we can hit our cluster from Boot.dev. We
should see message like "kubectl is now configured to use "minikube" cluster and 
"default" namespace by default".


# Dashboard

This will generate a dashboard application for our cluster. We can use this
dashboard to view and manage our cluster.

```bash
minikube dashboard --url
```

## Creating a Deployment

`kubectl create deployment` command will create a `deployment` for us. 2
parameters would be needed:
- Name of the deployment. (This can be anything, it is used to identify the
  deployment)
- The ID of the Docker image we want to deploy. (Full URL if we were not hosting
  the image on Docker Hub).

```bash
kubectl create deployment synergychat-web
--image=bootdotdev/synergychat-web:latest
```

This command will deploy container build from the docker image to our local k8s
cluster.

## Viewing deployments

```bash
kubectl get deployments
```

## Accessing the web page

By default, the resources inside of k8s run on private isolated network. They
are visible to other source within cluster, but not to the outside world. In
order to access the application from the local network, we need to perform some
port forwarding. First we have to run:

```bash
kubectl get pods
```

We should see soemthing like:

```bash
NAME                                   READY   STATUS    RESTARTS   AGE
synergychat-web-679cbcc6cd-cq6vx       1/1     Running   0          20m
```
Next we run:

```bash
kubectl port-forward <pod_name> 8080:8080
```

The application is now accessible at `http://localhost:8080`.
