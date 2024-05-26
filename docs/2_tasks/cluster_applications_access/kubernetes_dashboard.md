# Deploy and Access the k8s Dashboard

## Deploying the Dashboard UI

> [!NOTE]
> K8s Dashboard supports only Helm-based installation currently as it is faster
> and fives better control over all depndencies required by dashboard to run.

The Dashboad UI can be deployed using Helm charts:

```bash
helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard/
helm upgrade --install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard \
    --create-namespace \
    --namespace kubernetes-dashboard
```

## Accessing the Dashboard UI

> [!WARNING]
> Sample user created in the tutorial will have administrative privileges and is
> for educational purposes only.

To protect the cluster data, Dashboard deploys with minimal RBAC configuration
by default. Currently, Dashboard only supports logging in with a Bearer Token.
Token can be generated from the `ServiceAccount` user.

Create sample user `ServiceAccount` manifest in `kubernetes-dashboard` namespace:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  namespace: kubernetes-dashboard
```

In most cases after provisioning the cluster using `kops`, `kubeadm` or any
other popular tool, the `ClusterRole` `cluster-admin` already exists in the
cluster. We can use it to create only a `ClusterRoleBinding` for our
`ServiceAccount`. If it does not exist, then we need to create this role first
and grant required privileges manually.

Create sample user `ClusterRoleBinding` manifest:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: admin-user
  namespace: kubernetes-dashboard
```

Now we need to find the token that we can use to log in.

```bash
kubectl -n kubernetes-dashboard create token admin-user
```

It should print something like:

```bash
eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Ii
wia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlcm5ldGVzLWRhc2hib
2FyZCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi11c2Vy
LXRva2VuLXY1N253Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQ
ubmFtZSI6ImFkbWluLXVzZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYW
Njb3VudC51aWQiOiIwMzAzMjQzYy00MDQwLTRhNTgtOGE0Ny04NDllZTliYTc5YzEiLCJzdWIiOiJze
XN0ZW06c2VydmljZWFjY291bnQ6a3ViZXJuZXRlcy1kYXNoYm9hcmQ6YWRtaW4tdXNlciJ9.Z2JrQli
tASVwWbc-s6deLRFVk5DWD3P_vjUFXsqVSY10pbjFLG4njoZwh8p3tLxnX_VBsr7_6bwxhWSYChp9hw
xznemD5x5HLtjb16kI9Z7yFWLtohzkTwuFbqmQaMoget_nYcQBUC5fDmBHRfFvNKePh_vSSb2h_aYXa
8GV5AcfPQpY7r461itme1EXHQJqv-SN-zUnguDguCTjD80pFZ_CmnSE1z9QdMHPB8hoB4V68gtswR1V
La6mSYdgPwCHauuOobojALSaMc3RH7MmFUumAgguhqAkX3Omqd3rJbYOMRuMjhANqd08piDC3aIabIN
X6gP5-Tuuw2svnV6NYQ
```

We can also create a token with the secret which bound the service account and
the token will be saved in the `Secret`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: admin-user
  namespace: kubernetes-dashboard
  annotations:
    kubernetes.io/service-account.name: "admin-user"   
type: kubernetes.io/service-account-token
```

### Command line proxy

We can enable access to the Dashboard using `kubectl` command-line tool, by
running the following command:

```bash
kubectl proxy
```

`kubectl` will make the Dashboard available at:
http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy

> [!NOTE]
> The `kubeconfig` authentication method does not support external identity
> providers or X.509 certificate based authentication.
