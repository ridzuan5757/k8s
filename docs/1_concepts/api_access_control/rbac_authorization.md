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

Allow reading the resources `"nodes"` in the core group because a Node is
cluster-scoped, this must be in ClusterRole bound with a ClusterRoleBinding to
be effective:

```yaml
rules:
- apiGroups: [""]
  # at the HTTP level, the name of the resource for accessing node objects is
  # "nodes"
  resources: ["nodes"]
  verbs: ["get", "list", "watch"]
```

Allow GET and POST requests to the non-resource endpoint `/healthz` and all
subpaths must be in a ClusterRole bound with a ClusterRoleBinding to be
effective:

```yaml
rules:
  # '*' in a nonResourceURLs is a suffix glob match
- nonResourceURLs: ["/healthz", "/healthz/*"]
  verbs: ["get", "post"]
```

### Referring to subjects

A RoleBinding or ClusterRoleBinding binds a role to subjects. Subjects can be
groups, users, or ServiceAccounts.

k8s represents usernames as strings. These can be: plain names, such as 
"alice", email-style names, like "bob@example.com"; or numeric user IDs
represented as string. It is up to us as a cluster admin to configure the
authentication modules so that authentication produces usernames in the format
we want.

> [!CAUTION]
> The prefix `system:` is reserved for k8s system use, so we should ensure that
> we do not have users or groups with names that start with `system:` by
> accident. Other than this special prefix, the RBAC authorization system does
> not require any format for usernames.

In k8s, Authenticator modules provide group information. Groups, like users, are
represented as strings, and that string has no format requirements, other than
that the prefix `system:` is reserved.

ServiceAccounts have names prefixed with `system:serviceaccount:`, and belong to
groups that have names prefixed with `system:serviceaccounts:`.

> [!NOTE]
> - `system:serviceaccount:` (singular) is the prefix for service account
>   usernames.
> - `system:serviceaccounts:` (plural) is the prefix for service account groups.

### RoleBinding examples

The following examples are `RoleBinding` excerpts that only show the subjects
section. For a user named `alice@example.com`:

```yaml
subjects:
- kind: User
  name: "alice@example.com"
  apiGroup: rbac.authorization.k8s.io
```

For a group named `frontend-admins`:

```yaml
subjects:
- kind: Group
  name: "frontend-admins"
  apiGroup: rbac.authorization.k8s.io
```

For the default service account in the `kube-system` namespace:

```yaml
subjects:
- kind: ServiceAccount
  name: default
  namespace: kube-system
```

For all service accounts in the `qa` namespace:

```yaml
subjects:
- kind: Group
  name: system:serviceaccounts:qa
  apiGroup: rbac.authorization.k8s.io
```

For all service accounts in any namespace:

```yaml
subjects:
- kind: Group
  name: system:serviceaccounts
  apiGroup: rbac.authorization.k8s.io
```

For all authenticated users:

```yaml
subjects:
- kind: Group
  name: system:authenticated
  apiGroup: rbac.authorization.k8s.io
```

For all unauthenticated users:

```yaml
subjects:
- kind: Group
  name: system:unauthenticated
  apiGroup: rbac.authorization.k8s.io
```

For all users:

```yaml
subjects:
- kind: Group
  name: system:unauthenticated
  apiGroup: rbac.authorization.k8s.io
- kind: Group
  name: system:authenticated
  apiGroup: rbac.authorization.k8s.io
```

## Default roles and role bindings 

API servers create a set of default ClusterRole and ClusterRoleBinding objects.
Mange of these are `system:` prefixed, which indicates that the resource is
directly managed by the cluster control plane. All of the default ClusterRole
and ClusterRoleBinding are labeled with
`kubernetes.io/bootstrapping=rbac-defaullts`.

> [!Caution]
Take care when modifying ClusterRole and ClusterRoleBinding with names that have
`system:` prefix. Modifications to these resources can result in non-functional
clusters.

### Auto-reconciliation

At each startup, the API server updates default cluster roles with any missing
permissions, and updates default cluster role bindings with any missing
subjects. This allows the cluster to repair accidental modifications, and helps
to keep roles and role binding to `false`. Be aware that missing default
permissions and subjects can result in non-functional clusters.

Auto-reconciliation is enabled by default if the RBAC authorizer is active.

### API discovery roles

Default role bindings authorize unauthenticated and authenticated users to read
API information that is deemed safe to be publicly accessible (including
CustomResourceDefinitions). To disable anonymous unauthenticated access, add
`--anonymous-auth=false` to the API server configuration.

To view the configuration of these roles via `kubectl`:

```bash
kubectl get clusterroles system:discovery -o yaml
```

> [!NOTE]
> If we edit that ClusterRole, the changes will be overwritten on API server
> restart via auto-reconciliation. To avoid that overwriting, either do not
> manually edit the role or disable auto-reconciliation.

### k8s RBAC API discovery roles

Default ClusterRole: `system:basic-user`
Default ClusterRoleBinding: `system:authenticated` group
- Allows a user read-only access to basic information about themselves. Prior to
  version v1.14, this role was also bound to `system:unauthenticated` by
  default.

Default ClusterRole: `system:discovery`
Default ClusterRoleBinding: `system:authenticated` group
Allow read-only access API discovery endpoints needed to discover and negotiate
an API level. Prior to v1.14, this role was also bound to
`system:unauthenticated` by default.

Default ClusterRole: `system:public-info-viewer`
Default ClusterRoleBinding: `system:authenticated` and `system:unauthenticated`
group.
Allow read-only access to non-sensitive information about the cluster.

### User-facing roles

Some of the default ClusterRoles are not `system:` prefixed. These are inteded
to be user-facing roles. They include super-user roles(`cluster-admin`), roles
intended to be granted cluster-wide using ClusterRoleBindings, and roles
intended to be granted within particular namespaces using RoleBindings (`admin`,
`edit`, `view`).

User-facing ClusterRoles use ClusterRole aggregation to allow admins to include
rules for custom resources on these ClusterRoles. To add rules to the `admin`,
`edit`, or `view` roles, create a ClusterRole with one or more of the following
labels:

```yaml
metadata:
    labels:
        rbac.authorization.k8s.io/aggregate-to-admin: "true"
        rbac.authorization.k8s.io/aggregate-to-edit: "true"
        rbac.authorization.k8s.io/aggregate-to-view: "true"
```

Default ClusterRole: `cluster-admin`
Default ClusterRoleBinding: `system:masters` group
- Allows super-user access to perform any action on any resource. When used in a
  `ClusterRoleBinding`, it gives full control over every resource in the cluster
  and in all namespaces. When used in a `RoleBinding`, it gives full control
  over every resource in the role binding's namespace, including the namespace
  itself.

Default ClusterRole: `admin`
Default ClusterRoleBinding: None
- Allows admin access, intended to be granted within a namespace using a
  `RoleBinding`. If used in `RoleBinding`, allows read/write access to most
  resources in a namespace, including the ability to create roles and role
  bindings within the namespace. This role does not allow write access to
  resoure quota or to the namespace itself. This role also does not allow write
  access to EndpointSlices or Endpoints in clusters created using k8s v1.22+.

Default ClusterRole: `edit`
Default ClusterRoleBinding: None
- Allows read/write access to most objects in a namespace. This role does not
  allow viewing or modifying roles or role bindings. However, this role allows
  accessing Secrets and running Pods as any ServiceAccount in the namespace, so
  it can be used to gain the API access levels of any ServiceAccount in the
  namespace.
- This role also does not allow write access to EndpointSlices or Endpoints in
  clusters created using k8s v1.22+.

Default ClusterRole: `view`
Default ClusterRoleBinding: None
- Allows read-only access to see most objects in a namespace. It does not allow
  viewing roles or role bindings. This role does not allow viewing Secrets,
  since reading the contents of Secrets enables access to ServiceAccount
  credentials in the namespace, which would allow API access as any
  ServiceAccount in the namespace. (A form of privelege escalation).
