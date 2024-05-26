# Run Stateless Application Using Deployment

## Creating and exploring an nginx Deployment

We can run an application by creating `Deployment` object, and we can describe a
`Deployment` in a YAML file.

```bash
# nginx-deployment.yaml

apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2 # tells deployment to run 2 pods matching the template
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

Create the `Deployment`:

```bash
kubectl apply -f nginx-deployment.yaml
```

Display information about the `Deployment`:

```bash
kubectl describe deployment nginx-deployment
```

The output is similar to this:

```bash
Name:     nginx-deployment
Namespace:    default
CreationTimestamp:  Tue, 30 Aug 2016 18:11:37 -0700
Labels:     app=nginx
Annotations:    deployment.kubernetes.io/revision=1
Selector:   app=nginx
Replicas:   2 desired | 2 updated | 2 total | 2 available | 0 unavailable
StrategyType:   RollingUpdate
MinReadySeconds:  0
RollingUpdateStrategy:  1 max unavailable, 1 max surge
Pod Template:
  Labels:       app=nginx
  Containers:
    nginx:
    Image:              nginx:1.14.2
    Port:               80/TCP
    Environment:        <none>
    Mounts:             <none>
  Volumes:              <none>
Conditions:
  Type          Status  Reason
  ----          ------  ------
  Available     True    MinimumReplicasAvailable
  Progressing   True    NewReplicaSetAvailable
OldReplicaSets:   <none>
NewReplicaSet:    nginx-deployment-1771418926 (2/2 replicas created)
No events.
```

List the Pods created by the deployment:

```bash
kubectl get pods -l app=nginx
```

The output is similar to this:

```bash
NAME                                READY     STATUS    RESTARTS   AGE
nginx-deployment-1771418926-7o5ns   1/1       Running   0          16h
nginx-deployment-1771418926-r18az   1/1       Running   0          16h
```

Display information about the pod:

```bash
kubectl describe pod <pod-name>
```


