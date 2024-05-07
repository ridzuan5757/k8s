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

## Referring to resources

In the k8s API, most resources are represented and accessed using a string
representation of their object name, such as `pods` for a Pod. RBAC refers to
resources using exactly the same name that appears in the URL for the relevant
API endpoint. Some k8s API involve a subresource, such as the logs for a Pod. A
request for a Pod's logs looks like:

```bash
GET /api/v1/namespaces/{namespace}/pods/{name}/log
```

In this case, `pods` is the namespaced resource for Pod resources, and `log` is
a subresource of `pods`. To represent this in an RBAC role, use a slash `/` to
delimit resource and subresource. To allow a subject to read `pods` and also
access the `log` subresource for each Pods, we write:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
    namespace: default
    name: pod-and-pod-logs-reader
rules:
- apiGroups: [""]
  resources: ["pods", "pods/log"]
  verbs: ["get", "list"]
```

We can also refer to resources by name for certain requests through the
`resourceNames` list. When specified, requests can be restricted to individual
instances of a resource. Here is an example that restrcits its subject to only
`get` or `update` a ConfigMap named `my-configmap`:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
    namespace: default
    name: configmap-updater
rules:
- apiGroups: [""]
  # at the HTTP level, the name of the resource for accessing ConfigMap objects
  # is configmaos
  resources: ["configmaps"]
  resourceNames: ["my-configmap"]
  verbs: ["update", "get"]
```

> [!NOTE]
> We cannot restrict `create` or `deletecollection` requests by their resource
> name. For `create` this limitation is because the name of the new object may
> not be known at authorization time. If we restrict `list` or `watch` by
> resourceName, clients must include a `metadata.name` field selector in their
> `list` or `watch` request that matches the specified resourceName in order to
> be authorized. For example:
>
> ```bash
>   kubectl get configmaps --field-selector=metadata.name=my-configmap
> ```

Rather than referring to individual `resources`, `apiGroups`, and `verbs`, we
can use the wildcard `*` as symbol to refer to all such objects. For
`nonResourceURLs`, we can use the wildcard `*` as suffix glob match. For
`resourceNames`, an empty set means that everything is allowed.

Here is an example that allows access to perform any current and future action
on all current and future resources in the `example.com` API group. This is
similar to the built-in `cluster-admin` role.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
    namespace: default
    name: example.com-superuser
rules:
- apiGroups: ["example.com"]
  resources: ["*"]
  verbs: ["*"]
```

> [!CAUTION]
> Using wildcards in resource and verb entries could result in overly permissive
> access being granted to sensitive resources. For instance, if a new resource
> is added, or a new subresource is added, or a new custom verb is checked, the
> wildcard entry automatically grants access, which may be undesirable. The
> principle of least privelege should be employed, using specific resources and
> verbs to ensure only the permissions required for the workload to function
> correctly are applied.

### Aggregated ClusterRoles

We can aggregate several ClusterRoles into one combined ClusterRole. A
controller, running as part of the cluster control plane, watches for
ClusterRole objects with an `aggregationRule` set. The `aggregationRule` defines
a label selector that the controller uses to match other ClusterRole objects
that should be combined into the `rules` field of this one.

> [!CAUTION]
> The control plane overwirtes any values that we manually specify in the
> `rules` field of an aggregate ClusterRole. If we want to change or add rules,
> do so in the `ClusterRole` objects that are selected by the `aggregationRule`.

Here is an example of aggregated ClusterRole:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: monitoring
aggregationRule:
    clusterRolesSelectors:
    - matchLabels:
        rbac.example.com/aggregate-to-monitoring: "true"
# the control plane automatically fills in the rules
rules: []
```

If we create a new ClusterRole that matches the label selector of an existing
aggregated ClusterRole, that change triggers adding the new rules into the
aggregated ClusterRole. Here is an example that adds rules to the "monitoring"
ClusterRole, by creating another ClusterRole labeled
`rbac.example.com/aggregate-to-monitoring: true`.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: monitoring-endpoints
    labels:
        rbac.example.com/aggregate-to-monitoring: "true"
# when we create the monitoring-endpoints ClusterRole, the rules below will be
# added to the monitoring ClusterRole.
rules:
- apiGroups: [""]
  resources: ["services", "endpointslices", "pods"]
  verbs: ["get", "list", "watch"]
```

The default user-facing roles use ClusterRole aggregation. This lets us, as a
cluster administrator, include rules for custom resources, such as those served
by CustomResourceDefinitions or aggregated API servers, to extend the default
roles.

For example: the following ClusterRoles let the "admin" and "edit" default roles
manage the custom resource named CronTab, whereas the "view" role can perform
only read actions on the CronTab resources. We can assume that CronTab objects
are named "crontabs" in URLs as seen by the API server.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: aggregate-cron-tabs-edit
    labels:
        # add these permission to the "admin" and "edit" default roles.
        rbac.authorization.k8s.io/aggregate-to-admin: "true"
        rbac.authorization.k8s.io/aggregate-to-edit: "true"
rules:
- apiGroups: ["stable.example.com"]
  resources: ["crontabs"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

---

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
    name: aggrefate-cron-tabs-view
    labels:
        # add these permissions to the "view" default role
        rbac.authorization.k8s.io/aggregate-to-view: "true"
rules:
- apiGroups: ["stable.example.com"]
  resources: ["crontabs"]
  verbs: ["get", "list", "watch"]
```

### Role examples

The following examples are excerpts from Role or ClusterRole objects, showing
only `rules` section.

To adlow reading "pods" resources in the core API Group:

```yaml
rules:
- apiGroups: [""]
  # at the HTTP level, the name of the resource for access Pod object is pods
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
```

To allow reading/writing Deployment at the HTTP level (object with "deployments" in the resource part of their URL) in the "apps" aPI groups:

```yaml
rules:
- apiGroups: ["apps"]
  # at the HTTP level, the name of the resource for accessing Deployment object
  # is deployment
  resources: ["deployments"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
```

Allow reading Pods in the core API group, as well as reading or writing job
resources in the `"batch"` API group:

```yaml
rules:
- apiGroups: [""]
  # at the HTTP level, the name of the resource for accessing Pod objects is
  # "pods"
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["batch"]
  # at the HTTP level, the name of the resource for accessing Job object is
  # "jobs"
  resources: ["jobs"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
```

Allow reading ConfigMap named "my-config" (must be bound with a RoleBinding to
limit to a single ConfigMap in a single namespace):

```yaml
rules:
- apiGroups: [""]
  # at the HTTP level, the name of the resource for accessing ConfigMap object
  # is "configmaps"
  resources: ["configmaps"]
  resourceNames: ["my-config"]
  verbs: ["get"]
```
