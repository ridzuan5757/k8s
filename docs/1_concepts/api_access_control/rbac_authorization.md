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

#### ClusterRole example

A ClusterRole can be used to grant the same permission as a Role. Because
ClusterRoles are cluster-scoped, we can also use them to grant access to:
- cluster-scoped resources like nodes
- non-resource endpoints like `/healthz`
- namespaced resources like Pods across all namespaces

For example, we can use a ClusterRole to allow a particular user to run `kubectl
get pods --all-namespaces`. Here is an example of a ClusterRole that can be used
to grant read access to secrets in any particular namespace, or across all
namespaces, depending on how it is bound.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    # namespace mitted since ClusterRoles are not namespaced
    name: secret-reader
rules:
- apiGroups: [""]
  # at the HTTP level, the name of the resource for accessing Secret objects is 
  # "secrets"
  resources: ["secrets"]
  verbs: ["get", "watch", "list"]
```

The nae of a Role or a ClusterRole object must be a valid path segment name.

### RoleBinding and ClusterRoleBinding

A role binding grants the permissions defined in a role to a user or set of
users. It holds a list of subjects (user, groups, or service accounts), and a
reference to the role being granted. A RoleBinding grants permissions within a
specific namespace whereas a ClusterRoleBinding grants that access cluster-wide.

A RoleBinding may reference any Role in the same namespace. Alternatively, a
RoleBinding can reference a ClusterRole and bind that ClusterRole to the
namespace of the RoleBinding. If we want to bind a ClusterRole to all the
namespaces in the cluster, we use a ClusterRoleBinding.

The name of a RoleBinding or ClusterRoleBinding object must be a valid path
segment name.

#### RoleBinding example

Here is an example of a RoleBinding that grants the "pod-reader" Role to the
user "jane" with "default" namespace. This allows "jane" to read pods in the
"default" namespace.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
    name: read-pods
    namespace: default
subjects:
# more than one subject can be specified
- kind: User
  name: jane
  apiGroup: rbac.authorization.k8s.io
roleRef:
    # roleRef specifies the binding to a Role / ClusterRole
    kind: Role
    # this must match the name of the Role or the ClusterRole that need to be
    # bind to
    name: pod-reader
    apiGroup: rbac.authorization.k8s.io
```

A RoleBinding can also reference a ClusterRole to grant the permisison defined
in that ClusterRole to resources inside the RoleBinding's namespace. This kind
of reference lets we define a set of common roles across the cluster, then reuse
them within multiple namespaces.

For instance, even though the following RoleBinding refers to a ClusterRole,
"dave" (the subject, case sensitive) will only be able to read Secrets in the
"development" namespace, because the RoleBinding's namespace (in its metadata)
is "development".

```yaml
apiVersion: rbac.authorization.k8s.io/v1
# this role binding allows "dave" to read secrets in the "development"
# namespace
# we need to already have a ClusterRole named "secret-reader"
kind: RoleBinding
metadata:
    name: read-secrets
    # the namespace of RoleBinding determines where the permissions are granted.
    # this only grants permissions within the "development" namespace.
    namespace: development
subjects:
- kind: User
  name: dave # name is case sensitive
  apiGroup: rbac.authorization.k8s.io
roleRef:
    kind: ClusterRole
    name: secret-reader
    apiGroup: rbac.authorization.k8s.io
```
### ClusterRoleBinding example

To grant permissions across a whole cluster, we can use ClusterRoleBinding. The
following ClusterRoleBinding allows any user in the group "manager" to read
secrets in any namespace.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
# this cluster role binding allows anyone in the "manager" group to read
# secrets in any namespace.
kind: ClusterRoleBinding
matadata:
    name: read-secrets-global
subjects:
- kind: Group
  # name is case sensitive
  name: manager
  apiGroup: rbac.authorization.k8s.io
roleRef:
    kind: ClusterRole
    name: secret-reader
    apiGroup: rbac.authorization.k8s.io
```

After we create a binding, we cannot change the Role or ClusterRole that it
refers to. If we try to change a binding's `roleRef`, we get a validation error.
If we do want to change to `roleRef` for a binding, we need to remove the
binding object and create a replacement. 

There are two reasons for this restriction:
- Making `roleRef` immutable allows granting someone `update` permission on an
  existing binding object, so that they can manage list of subjects, without
  being able to change the role that is granted to those subjects.
- A binding to a different role is a fundamentally different binding. Requiring
  a binding to be deleted/recreated in order to chnage the `roleRef` ensures the
  full list subjects in the binding is intended to be granted the new role as
  opposed to enabling or accidentally modifying only the roleRef without
  verifying all of the existing subjects should be given the new role's
  permissions.

The `kubectl auth reconcile` command-line utility creates or updates a manifest
file containing RBAC objects, and handles deleting and recreating binding
objects if required to change the role they refer to.


