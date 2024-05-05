# RBAC Authorization

Role-based access control RBAC is a method of regulating access to computer or
network resources based on the roles of individual users within the
organization.

RBAC authorization uses the `rbac.authorization.k8s.io` API Group to drive
authorization decisions, allowing us to dynamically configure policies through
k8s API.

To enable RBAC, start the API server with the `--authorization-mode` flag set to
a comma-separated list that includes `RBAC`. For example:

```bash
kube-apiserver --authorization-mode=Example,RBAC --other-options --more-options
```

## API objects

The RBAC API declares four kinds of k8s objects:
- Role
- ClusterRole
- RoleBinding
- ClusterRoleBinding

We can describe or amend the RBAC objects using tools such as `kubectl`, just
like any other k8s object.

> [!CAUTION]
> These objects, by design, impose access restrictions. If we are making changes
> to a cluster as we learn, check privelege escalation prevention and
> bootstrapping to understand how those restriction can prevent us from making
> some changes.

### Role and ClusterRole

An RBAC Role or ClusterRole contains rules that represent a set of permissions.
Permissions are purely additive. There are no "deny" rules.

A Role always sets permissions within a particular namespace; when we create a
Role, we have to specify the namespace it belongs in.

ClusterRoles have several uses. We can use a ClusterRole to:
- Define permissions on namespaced resources and be granted access within
  individual namespaces.
- Define permissions on namespaced resources and be granted access across all
  namespaces.
- Define permissions on cluster-scoped resources.

If we want to define a role within a namespace, use a Role; if we want to define
a role cluster-wide, use a ClusterRole.

#### Role example

Here is an example Role in the default namespace that can be used to grant read
access to Pods:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
    namespace: default
    name: pod-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
```


