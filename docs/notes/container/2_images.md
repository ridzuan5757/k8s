# Images

A container image represents binary data that encapsulates an application and
all its software dependencies. Contaienr images are executable software bundles
that can run standalone and that make very well defined assumptions about their
runtime environment.

We typically create container image of the application and push it to a registry
such as Dockerhub before referring to it in a pod.

Pod : Set of running containers in a cluster.

## Image names

Container images are usually given a name such as:
- `pause`
- `example/mycontainer`
- `kube-api server`

Image can also include a registry hostname. For example:
- `fictional.registry.example/imagename`

And possibly port number as well. For example:
- `fictional.registry.example:10443/imagename`

If we do not specify registry hostname, k8s assumes that we mean the Docker
public registry.

After the image name part, we can add a tag in a same way we would when using
commands like `docker` or `podman`. Tags lets us identify different versions of
the same series of images.

Image tags consists of lowercase and uppercase letters, digits, underscores `_`,
periods `.`, and dashes `-`. There are additional rules about where we can place
the separator characters mentioned inside an image tag. If we do not specify a
tag, k8s assumes we mean the tag `latest`.

## Updating images

When we first create a Deployment, StatefulSet, Pod or other object that include
a pod template, then by default the pull policy of all containers in that pod
will be set to `IfNotPresent` if is is not explicitly specified. This policy
causes the kubelet to skip pulling an image if it already exists.

### Image pull policy

The `imagePullPolicy` for a container and the tag of the image affect when the
kubelet attempts to pull the specified image. Some of the possible values are:

##### `IfNotPreset`

Image is pulled only if it is not already present locally.

##### `Always`

Every time the kubelet launches a container, the kubelet queries the container
image registry to resolve the name to an image digest. If the kubelet has a
container image with that exact digest cached locally, the kubelet uses its
cached image; otherwise, the kubelet pulls the image with the resolved digest,
and uses that image to launch the container.

##### `Never`

The kubelet does not try fetching the image. If the image is somehow already
present locally, the kubelet attempts to start the container, otherwise, startup
fails.



We should avoid using the `:latest` tag when deploying containers in production
as it is harder to trach which version of the image is running and more
difficult to rollback properly. Instead, specify a maningful tag such as
`v1.42.0` and / or a digest.

To make sure the pod always use the same version of a container image, we can
specify the image's disgest; replace `<image-name>:<tag>` with
`<image-name>@<disgest>` for example:
`image@sha256:45b23dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2`.

When using image tags, if the image registry were to change the code that the
tag on the image represents, we might end up with a mix of pods running the old
and new code. An image digest uniquely identifies specific version of the image,
so k8s runs the same code every time it starts a container with that image name
and digest specified. Specifying an image by digest fixes the code that we run
so that a change at the registry cannot lead to that mix of versions.

There are 3rd-party admission controllers that mutate pods and pod templates
when they are created, so that the running workload is defined based on the
image digest rather than tag. That might be useful if we want to make sure that
all of our workload is running the same code no matter what tag changes happen
at the registry.

### Default image pull policy

When we or a controller submit a new pod to the new API server, our cluster sets
the `imagePullPolicy` when specific conditions are met:
- If `imagePullPolicy` is omitted, but image digest for the container is 
  specified, the `imagePullPolicy` is automatically set to `IfNotPresent`.
- If `imagePullPolicy` is omitted and the tag for the container image is
  `:latest`, `imagePullPolicy` is automatically set to `Always`.
- If `imagePullPolicy`  is omitted and tag for the container image is not
  specified, the `imagePullPolicy` is automatically set to `Always`.
- If `imagePullPolicy` is omitted and tag for the container image is specified
  and the image tag is not `:latest`, the `imagePullPolicy` is automatically set
  to `IfNotPresent`.

The value of `imagePullPolicy` of the container is always set when the object is
first created, and is not updated if the image's tag or digest later changes.

For example, if we create a deployment with an image whose tag is not `:latest`,
and later update that deployment's image to a `:latest` tag, the `imagePullPolicy`
field will not change to `Always`. We must manually change the pull policy of
any object after its initial creation.

### Required image pull

If we woul like to always force a pull, we can do one of the following:
- Set the `imagePullPolicy` to `Always`
- Omit `imagePullPolicy` and use `:latest` tag for the image to use. k8s will
  set the policy to `Always` when we submit the pod.
- Omit the `imagePullPolicy` and the tag for the image to use. k8s will set the
  policy to `Always` when the pod is submitted.
- Enabled `AlwaysPullImage` admission controller.

### `ImagePullBackOff`

When a kubelet starts creating containers for a pod using a container runtime,
it might be possible the container is in waiting state because of 
`ImagePullBackOff`.

The status `ImagePullBackOff` means that a container could not start because k8s
could not pull a container image for reasons such as invalid image name, or
pulling from a private registry wiithout `imagePullSecret`. The `BackOff` part
indicates that k8s will keep trying to pull the image, with an increasing
back-off delay.

k8s raises the delay between each attempt until it reaches a compiled-in limit
which is 300 seconds / 5 minutes.

### Image pull per runtime class

k8s includes alpha support for performing image pulls based on the `RuntimeClass` 
of a pod. 

If we enable `RuntimeClassInImageCriApi` feature gate, the kubelet references
container images by a tuple of image name, runtime handler rather than just the
image name or digest.

The container runtime may adapt its behaviour based on the selected runtime
handler. Pulling images based on runtime class will be helpful for VM based
containers like windows hyperV containers.

## Serial and parallel image pulls

By default, kubelet pulls image serially, in other words, kubelet sends only one
image pull request to the image service at a time. Other image pull requests
have to wait until on being processed is complete.

Nodes make image pull decisions in isolation. Even when we use serialized image
pulls, two different nodes can pull the same image in parallel.

If we would like to enable parallel image pulls, we can set the field
`serializeImagePulls` to false in the kubelet configuration. With
`serializeImagePulls` set to false, image pull requests will be sent to the
image service immediately, and multiple images will be pulled at the same time.

When enabling parallel image pulls, please make sure the image service of the
container runtime can handle parallel image pulls.

The kubelet never pulls multiple images in parallel on behalf of one pod. For
example, if we have a pod that has an init container and an applicaiton
container, the image pulls for the two containers will not be parallelized.
However if we have 2 pods that use different images, the kubelet pulls the
images in paralellel on behalf of 2 different pods, when parallel image pulls is
enabled.

### Maximum parallel image pulls

When `serializeImagePulls` is set to false, the kubelet defaults to no limit on
the maximum number of images being pulled at the same time. If we would like to
limit the number of parallel image pulls, we can set the field
`maxParallelImagePulls` in kubelet configuration. With `maxParallelImagePulls`
set to n, only n images can be pulled at the same time, and any image pull
beyond n will have to wait until at least one ongoing image pull is complete.

Limiting the number parallel image pulls would prevent image pulling from
consuming too much network bandwidth or disk I/O, when parallel image pulling is
enabled.

We can set `maxParallelImagePulls` to a positive number that is greater than or
equal to 1. If we set `maxParallelImagePulls` to be greater than or equal to 2,
we must set the `serializeImagePulls` to false. The kubelet will fail to start
with invalid `maxParallelImagePulls` settings.

## Multi-architecture images with image indexes

As well as providing binary images, a container registry can also serve a
container image index. An image index can point to multiple  manifests for
architecture-specific versions of a container. The idea is that we can have a
name for an image for example `pause`, `example/mycointainer`, `kube-apiserver`
and allow different systems to fetch the right binary image for the machine
architecture they are using.

k8s itself typically names container images with a suffex `-$(ARCH)`. For
backward compatibility, please generate the older images with suffixes. The idea
is to generate say `pause` image which has the manifest for all the arch(es) and
say `pause-amd64` which is backwards compatible for older configurations or YAML
files which may have hard coded the images with suffixes.

## Private registry

Private registries may require keys to read images from them. Credentials can be
provided in serveral ways:
- Configuring nodes to authenticate to a private registry
    - All pods can read any configured private registries
    - Requires node configuration by cluster administrator
- Kubelet credential provider to dynamically fecth credentials for private
  registries. Kubelet can be configured to use credential provider exec plugin
  for the respective private registry.
- Pre-pulled images
    - All pods can use any image cached on a node
    - Requires root access to all nodes to set up
- Specifying `ImagePullSecrets` on a pod
    - Only pods which provide own keys can access the private registry
- Vendor specific or local extensions
    - If we are using a custom node configuration, we or the cloud provider can
      implement the mechanism for authenticating the node to the container
      registry

### Configuring nodes to authenticate to a private registry

Specific instructions for setting credentials depends on the container runtime
and registries chosen. We should refer to the solution's documentation for the
most accurate information.

#### Dockerhub

##### Preparation

We need to have k8s cluster, and the kubectl CLI must be configured to
communicate with the cluster. It is recommended to run this on a cluster with at
least two nodes that are not acting as control plane hosts. Docker CLI and a
Docker ID with known password would be essential.

###### Login

On the machine, authenticate with a registry in order to pull a private image.
Use the `docker` tool to log in to the DockerHub.

```bash
docker login
```

When prompted, enter the Docker ID, and then the crendential we want to use,
access token or the password for the Docker ID. The login process creates or
updates a `config.json` file that holds an authorization token.

```bash
cat ~/.docker/config.json
```

The output contains a section similar to this:

```json
{
    "auths": {
        "https://index.docker.io/v1/": {
            "auth": "c3R...zE2"
        }
    }
}
```
If we are using Docker crendial store, we won't see that `auth` enttry but a
`credStore` entry with the name of the store as value. In that case, we can
create a secret directly.

###### Create a secret based on existing credentials

A k8s cluster uses the secret of `kubernetes.io/dockerconfigjson` type to
authenticate with a container registry to pull a private image. If we alreay ran
`docker login`, we can copy the crendential into k8s:

```bash
kubectl create secret generic regcred \
    --from-file=.dockerconfigjson=~/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson
```

If we need more control, for example to set a namespace or a label on the new
secret, then the secret can be customized before stored. Be sure to:
- Set the name of the data item to `.dockerconfigjson`
- Base64 encode the docker configuration file and then paste that string,
  unbroken as the value for the field `data[".dockerconfigjson"]`
- Set `type` to `kubernetes.io/dockerconfigjson`

For example:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: myregistrykey
  namespace: awesomeapps
data:
  .dockerconfigjson: UmVhbGx5IHJlYWxseSByZWVlZWVlZWVlZWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGx5eXl5eXl5eXl5eXl5eXl5eXl5eSBsbGxsbGxsbGxsbGxsbG9vb29vb29vb29vb29vb29vb29vb29vb29vb25ubm5ubm5ubm5ubm5ubm5ubm5ubm5ubmdnZ2dnZ2dnZ2dnZ2dnZ2dnZ2cgYXV0aCBrZXlzCg==
type: kubernetes.io/dockerconfigjson
```

###### Troubleshooting

- `error: no objects passed to create`
    - The base64 encoded string is invalid
- `Secret "myregistrykey" is invalid: data[.dockerconfigjson]: invalid value ...`
    - The base64 encoded string in the data was successfully decoded, but could
      not be parsed as a `.docker/config.json` file.

###### Creating a Secret by providing crendentials on the CLI

Secret can be created using `kubectl` CLI:

```bash
kubectl create secret docker-registry regcred --docker-server=<your-registry-server> --docker-username=<your-name> --docker-password=<your-pword> --docker-email=<your-email>
```

Where:
- `<your-registry-server>` is the private container registry FQDN. Use
  `https://index.docker.io/v1/` for DockerHub.
- <your-name> is the registry username
- <your-pword> is the registry password
- <your-email> is the registry email

If this command is successful, a Secret called `regcred` will be created as the
registry credentials in the cluster.

Typing secrets on the command line may store them in the shell history
unprotected, and those secrets might also by visible to other users on the
machine while `kubectl` is running.

###### Inspecting `regcred` secret

To understand the contents of the `regcred` secret created, start by viewing the
secret in YAML format:

```bash
kubectl get secret regred --output=yaml
```

The output is similar to this:

```yaml
apiVersion: v1
kind: Secret
metadata:
  ...
  name: regcred
  ...
data:
  .dockerconfigjson: eyJodHRwczovL2luZGV4L ... J0QUl6RTIifX0=
type: kubernetes.io/dockerconfigjson
```

The value of the `.dockerconfigjson` field is a base64 representation of the
Docker credentials. To understand what is in the `.dockerconfigjson` field, we
can convert the secret data to a readable format:

```bash
kubectl get secret regcred --output="jsonpath={.data.\.dockerconfigjson}" \ 
    | base64 --decode
```

The output is similar to this:
```bash
{
    "auths":{
        "your.private.registry.example.com":{
            "username":"janedoe",
            "password":"xxxxxxxxxxx",
            "email":"jdoe@example.com",
            "auth":"c3R...zE2"
        }
    }
}
```

To understand what is in the `auth` field, convert th ebase64-encoded data to a
readable format:

```bash
echo "c3R...zE2" | base64 --decode
```

The output, username and password concatenated with `:` character, similar to
this:

```bash
janedoe:xxxxxxxxxxx
```

Notice that the secret data contains the authorization token similar to the
local `~/.docker/config.json` file. 

##### Creating pod with secret

Here is the manifest for an example pod that needs access to the docker
credentials in `regcred`:


```yaml
apiVersion: v1
kind: Pod
metadata:
  name: private-reg
spec:
  containers:
  - name: private-reg-container
    image: <your-private-image>
  imagePullSecrets:
  - name: regcred

```

To pull the image from the private registry, k8s need credentials. The 
`imagePullSecrets` field in the configuration file specifies that k8s should get
the credentials from a Secret named `regred`.

To use image pull secrets for a pod or a deployment, or other object that has a
pod template, we need to make sure that the appropiate secret does exist in the
right namespace. The namespace to use is the same namespace where the po dis
defined.

In the case the pod fails to start with the status `ImagePullBackOff`, view the
pod events:

```bash
kubectl describe pod private-reg
```

If we then see an event with the reason set to
`FailedToRetrieveImagePullSecret`, k8s can't find a secret with name `regred`.
If we specify that a pod needs image pull credentials, the kubelet checks that
it can access that secret before attempting to pull the image.

Make sure that the secret that is specified in the configuraiton file exists,
and that its name is applied properly.

```bash
Events:
  ...  Reason                           ...  Message
       ------                                -------
  ...  FailedToRetrieveImagePullSecret  ...  Unable to retrieve some image pull 
                                             secrets (<regcred>); attempting to 
                                             pull the image may not succeed.
```

### Kubelet credential provider for authenticated image pulls

This approach is especially suitable when kubelet needs to fetch registry
crendentials dynamically. Most commonly used for registries provided by cloud
providers where auth tokens are short-lived.

We can configure the kubelet to invoke a plugin binary to dynamically fetch
registry crendials for a container image. This is the most robust and versatile
way to fetch credentials for private registries, but also requires kubelet-level
configuration to enable.

### Interpretation of `config.json`

The interpretation of `config.json` varies between the original Docker
implementation and the k8s interpretation. In Docker, the `auths` keys can only
specify root URLs, whereas k8s allow glob URLs as well as prefix-matched paths.

The only limitation is that glob patterns (`*`) have to include the dot (`.`)
for each subdomain. The amount of matched subdomains has to be equal to the
amount of glob patterns (`*.`), for example:
- `*.kubernetes.io` will not match `kubernetes.io` but `abc.kubernetes.io`
- `*.*.kubernetes.io` will not match `abc.kubernetes.io` but
  `abc.def.kubernetes.io`
- `prefix.*.io` will match `prefix.kubernetes.io`
- `*-good.kubernetes.io` will match `prefix-good.kubernetes.io`

This would means that a `config.json` like this is valid:

```json
{
    "auths": {
        "my-registry.io/images": { "auth": "…" },
        "*.my-registry.io/images": { "auth": "…" }
    }
}
```

Image pull operations would now pass the credentials to the CRI container
runtime for every valid pattern. For example, the following container will match
successfully:
- `my-registry.io/images`
- `my-registry.io/images/my-image`
- `my-registry.io/images/another-image`
- `sub.my-registry.io/images/my0image`

But not:
- `a.sub.my-registry.io/images/my-image`
- `a.b.sub.my-registry.io/images/my-image`

The kubelet performs image pulls sequentially for every found credential. This
means, that multiple entries in `config.json` for different paths are possible
too.

```json
{
    "auths": {
        "my-registry.io/images": {
            "auth": "…"
        },
        "my-registry.io/images/subpath": {
            "auth": "…"
        }
    }
}
```

If not a container specifies an image `my-registry.io/images/subpath/my-image`
to be pulled, then the kubelet will try to download them from both
authentication sources if one of them fails.

### Pre-pulled images

This approach is suitable if we can control node configuration. It will not work
reliably if the cloud provider manages nodes and replaces them automatically.

By defauly, the kubelet tries to pull each image from the specified registry.
However, if the `imagePullPolicy` property of the container is set to
`IfNotPresent` or `Never`, then a local image is used preferentially or
exclusively, respectively.

If we want to rely on pre-pulled images as a substitute for registry
authentication, we must ensure all nodes in the cluster have the same pre-pulled
images.

This can be used to preload certain images for speed or as an alternative to
authenticating to a private registry.

All pods will have read access to any pre-pulled images.

### Specifying `imagePullSecrets` on pod

This is the recommended approach to run contaienrs based on images in private
registries.

k8s supports specifying container image registry keys on a pod. 
`imagePullSecrets` must all be in the same namespace as the pod. The referenced
secrets must be of type: 
- `kubernetes.io/dockercfg`
- `kubernetes.io/dockerconfigjson`

### Creating a secret with Docker config

In order to authenticate a private registry, th efollowing information would be
needed:
- username
- registry password
- client email address
- hostname

```bash
kubectl create secret docker-registry <name> \
  --docker-server=DOCKER_REGISTRY_SERVER \
  --docker-username=DOCKER_USER \
  --docker-password=DOCKER_PASSWORD \
  --docker-email=DOCKER_EMAIL
```

If we already have a docker credential file, then rather than using the above
command, we can just import the credentials file as k8s secrets.

This is particularly useful if we are using multiple private contaienr
registries as `kubectl create secret docker-registry` creates a secret that only
works with a single private registry.

Pods can only reference image pull secrets in their own namespace, so this
process needs to be done one time per namespace.

### Referring to an `imagePullSecrets` on a pod

Now we can create pods which reference that secret by adding an `imagePullSecrets` section to a pod definition. Each item in the `imagePullSecrets` array can only reference a secret in the same namespace.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: foo
  namespace: awesomeapps
spec:
  containers:
    - name: foo
      image: janedoe/awesomeapp:v1
  imagePullSecrets:
    - name: myregistrykey
```

This needs to be done for each pod that is using a private registry. However,
setting of this field can be automated by setting the `imagePullSecrets` in a
`ServiceAccount` resource.

We can use this in conjuction with a per-node `.docker/config.json`. The
credentials will be merged.

## Use cases

There are number of solutions for configuring private registries. Here are some
common use cases and suggested solutions.
- Cluster running only open source images. No need to hide images.
    - Use public image from the public registries.
        - No configuration required.
        - Some cloud providers automatically cache or mirror public images,
          which improves  availability and reduce time to pull images.
- Cluster running some proprietary images which should be hidden to those
  outside the company, but visible to all cluster users.
    - Use a hosted private registry.
        Manual configuration may be required on the nodes that need to access to
        private registry.
    - Alternatively, run internal private registry behind firewall with open
      read access.
        - No k8s configuration required.
    - Use a hosted container image registry service that controls image access.
        - It will work better on cluster autoscaling than manual mode
          configuration.
    - Or, on a cluster where changing the node configuraiton is inconvenient,
      use `imagePullSecrets`.
- Cluster with proprietary images, a few of which require stricter access
  control.
    - Ensure `AlwaysPullImage` admission controller is active. Otherwise, all
      pods potentiall have access to all images.
    - Move sensitive data into a `Secret` resource, instead of packaging it in
      the image.
- Multi-tenant cluster where each tenant needs own private registry.
    - Ensure `AlwaysPullImage` admission controller is active. Otherwise, all
      pods of all tenants potentially have access to all images.
    - Run a private registry with authorization required.
    - Generate registry credential for each tenant, put into secret, and
      populate secret to each tenant namespace.
    - The tenant adds that secret to `imagePullSecrets` of each namespace.

If access to multiple registries is needed, one secret for each registries can
be created.

## Legacy built-in kubelet credential provider

In older version of k8s, the kubelet had a direct integration with cloud
provider credentials. This give the ability  to dynamically fetch credentials
for image registries.

3 built in implementation of the kubelet credential provider integration:
- ACR 
- ECR
- GKE

k8s v1.26 through v1.29 do not include legacy mechanism, so we would need to
either:
- Configure a kubelet image credential provider on each node.
- Specify image pull credentials using `imagePullSecrets` and at least one
  `Secret`
