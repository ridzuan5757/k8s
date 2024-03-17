# `ConfigMaps`

There are several ways to manage environment variables in k8s. One of the most
common ways is to use `ConfigMaps`. `ConfigMaps` allow us to decouple our
configuration from our container images, which is important because we do not
want to have to rebuild our images every time we want to change a configuration
value.

In `Dockerfile`, we can set environment variables like this:

```bash
ENV PORT=3000
```

The trouble is, that means that everyone using that image will have to use port
3000. It also means that if we want to change the port, we have to rebuild the 
image.

In k8s, we can use YAML file as `ConfigMap` to modify the environment variable
of the container.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
    name: synergychat-api-configmap
data:
    API_PORT: 8080
```

We can then apply the `ConfigMap` using `kubectl apply`:

```bash
kubectl apply -f api-configmap.yaml
```

In order to validate that the config map was created, we can use the following
command:

```bash
kubectl get configmaps
```

# Applying the `ConfigMap`

Now that we have a `ConfigMap`, we need to connect it to our deployment. We need
to link the deployment config file with the ConfigMap config file. In the
deployment file, we have to add the following info:

```yaml
containers:
    env:
        - name: API_PORT
          valueFrom:
            configMapKeyRef:
                name: synergychat-api-configmap
                key: API_PORT
```

This tells k8s to set the `API_PORT` environment variable to the value of the
`API_PORT` key in the `synergychat-api-configmap` `ConfigMap`. We then apply the
deployment file:

```bash
kubectl apply -f api-deployment.yaml
```

Once it is applied, we should be able to take a look at the pods and see that a
new API pod has been deployed without crashing. To verify that it is working,
port forward the pod's to our local machine:

```bash
kubectl port-forward <pod_name> 8080:8080
curl http://localhost:8080
```

# ConfigMap Security Issue

`ConfigMap` is a great way to manage innocent environment variables in k8s.
Things such as:
- Ports
- URLs of other services
- Feature flags
- Settings that change between environment, such as `DEBUG` mode

However, they are not cryptographically secure. `ConfigMap` are not encrypted,
and they can be accessed by anyone with access to the cluster. If we need to
store sensitive information, we should use k8s Secrets or a third-party
solution.
