# Creating a Cluster

```bash
kind create cluster
```

This will bootstrap a cluster using preuilt node image. Prebuilt images are
hosted at `kindest/node`. To specify another image use the `--image` flag:

```bash
kind create cluster --image=...
```

Using different image allows us to change the k8s version of the created
cluster.
