# Custom Resources

Custom resources are extensions of the K8s API. This page discusses when to add
a custom resource to the K8s cluster and when to use a standalone service. It
describes the two methods for adding custom resources and how to choose between
them.

### `resource`

A resource is an endpoint in the k8s API that stores a collection of API objects
of a certain kind; for example the built-in pods resource contains a collection
of Pod objects.

### `custom resource`

Extension of the k8s API that is not necessarily avaialbe in a default k8s
installation. It represents a customization of a particular k8s installation.
However, many core k8s functions are now built using custom resources, making
k8s more modular.

Custom resources can appear and disappear in a running cluster through dynamic
registration, and cluster admins can update custom resources independently of
the cluster itself. once a custom resource is installed, users can create and
access its object using `kubectl` just as they do for built-in resources like
Pods.

## Custom Controllers

On their own, custom resources let us store and retrieve structured data. When
we combine a custom resource with a custom controller, custom resources provide
a true declarative API.

The k8s declarative API enforces a separation of responsibilities. We declare
the desired state of the resource. The k8s controller keeps the current state of
k8s objects in sync with the declared desired state. This is contrast to an
imperative API, where we instruct a server what to do.

We can deploy and update a custom controller on a running cluster, independently
of the cluster's lifecycle. Custom controllers can work with any kind of
resource, but they are especially effective when combined with custom resources.

The operator pattern combines custom resources and custom controllers. We can
use custom controllers to encode domain knowledge for specific applications into
an extentions of the k8s API.

When creating a new API, consider whether to aggregate the API with the k8s
cluster APIs or let the API standalone.

### Consider API aggregation if:
- The API is declarative.
- We want the new types to be readable and writable using `kubectl`.
- We want to view the new types in a K8s UI, such as dashboard, alognside
  built-in types.
- We are developing new API.
- We are willing to accept the format restriction that K8s puts on REST resource
  paths, such as API Groups and Namespaces.
- Reources are naturally scoped to a cluster or namespaces of a cluster.
- We want to resue K8s API support features.

### Prefer standalone API if:
- The API does not fit the declarative model.
- `kubectl` support is not required.
- k8s UI support is not required.
- We already have a program that serves the API and works well.
- We need to have specific REST path to be compatible with an already defined
  REST API.
- Cluster or namespace scoped resources are a poor fit, you need control over
  the specifics of resources paths.

### Declarative APIs

In a declarative API, typically:
- The API consists of a relatively small number of relatively small objects
  (resources).
- The objects define configuration of applications of infrastructure.
- The object are updated relatively unfrequently.
- Humans often need to read and write the objects.
- The main operations on the objects are CRUD in nature.
- Transactions across objects are not required: the aPI represents a desired
  state, not an exact state.

Imperative APIs are not declarative. Signs that the API might not be declarative
include:
- The client says "do this", and then gets a synchronous response back when it
  is done.
- The client says "do this", and then gets an operation ID back, and has to
  check a separate Operation object to determine completion of the request.
- We talk about Remote Procedure Calls RPCs.
- Directly storing large amount of data. For example, a few kB per object or
  thousands of objects.
- High bandwidth access (10s of requests per second sustained) needed.
- Store end-user data (such as images, PII, etc.) or other large-scale data
  processed by applications.
- The natural operations on objects are not CRUD in nature.
- The API is not easily modeled as objects.
- We chose to represent pending operations with an operation ID or an operation
  object.

## ConfigMap vs Custom Resource

Use a ConfigMap if any of the following apply:
- There is an existing, well-documented configuration file format, such as
  `mysql.cnf` or `pom.xml`.
- We want to put the entire configuration into one key of a ConfigMap.
- The main use of the configuration file is for a program running in a Pod or
  environment variablein a pod, rather than the k8s API.
- We want to perform rolling updates via Deployment, etc., when the file is
  updated.

> [!NOTE]
> Use a Secret for sensitive data, which is similar to a ConfigMap but more
> secure.

Use a custom resource CRD or Aggregated API if most of the following apply:
- We want to use k8s client libraries and CLIs to create and update the new
  resource.
- We want top-level support form `kubectl`. For example, `kubectl get my-object
  object-name`.
- We want to write automation that handles updates to the object.
- We want to use k8s API conventions like `.spec`, `.status`, and `.metadata`.
- We want the object to be an abstraction over a collection of controlled
  resources, or a summarization of other resources.

## Adding custom resources

k8s provides two ways to add custom resources to the cluster:
- CRDs are simple and can be created without any programming.
- API Aggregation requires programming, but allows more control over API
  behaviour like how data is stored and conversion between API versions.

k8s provides these two options to meet the needs of different users, so that
neither ease of use nor flexibility is compromised.

Aggregated APIs are subordinate API servers that sit behind the primary API
server, which acts as a proxy. This arrangement is called API aggregation AA. To
users, the k8s API appears extended.

CRDs allow users to create new types of resources without adding another API
server. We do not need to understand API aggregation to use CRDs.

Regardless of how they are installed, the new resources are referred to as
Custom Resources to distinguih them from built-in k8s resources like Pods.

> [!NOTE]
> Avoid using a Custom Resource as data storage for application, end user, or
> monitoring data. Architecture design that store application data within the
> k8s API typically represent a design that is too closely coupled.
>
> Architecturally, cloud native application architectures favor loose coupling
> between components. If part of the workload requires a backing service for its
> routine operation, run that backing service as a component or consume it as an
> external service. This way, the workload does not rely on k8s API for its
> normal operation.

## CustomResourceDefinition

The CRD API resource allows us to define custom resources. Defining a CRD object
creates a new custome resource with a name and schema that we specify, The k8s
API serves and handles the storage of the custom resource. The name of a CRD
object must be a valid DNS subdomain name. 

This frees us from writing our own API server to handle csutom resource, but the
generic nature of the implementation means that we have less flexibility than
with API aggregation.

## API server aggregation

Usually, each resource in k8s API requires code that handles REST requests and
manages persistent storage of objects. The main k8s API server handles built-in
resources like Pods and Services, and can also generically handle custom
resources through CRDs.

The aggregation layer wllows us to provide specialized implementatuins for
custom resources by writing and deploying our own API server. The main API
server delegates requests to the API server for the custom resources that we
handle, making them available to all of its clients.

## Choosing a Method for Adding Custom Resource

CRDs are easier to use. Aggregated APIs are more flexible. Choose the method
that best meets your needs. 

Typically, CRDs are a good fit if:
- We have handful of fields.
- We are using the resource within our own company, or as part of a small
  open-serouce project (as opposed to a commercial product)

### Ease of use

#### CRDs
- Do not require programming. Users can choose any language for a CRD
  controller.
- No additional service to run; CRDs are handled by API server.
- No ongoing support once the CRD is created. Any bug fixes are picked up as
  part of normal k8s Master upgrades.
- No need to handle multiple versions of the API. For example, when we control
  the client for this resource, we can upgrade it in sync with API.

#### Aggregated API
- Require programming and building binary and image.
- An additional service to create and could fail.
- May need to periodically pickup bug fixes from upstream and rebuild and update
  the Aggregated API server.
- We need to handle multiple versions of the API. For example, when developing
  an extension to share with the world.

### Advanced features and flexibility

Aggregated APIs offer more advanced API features and customization of other
features; for example, the storage layer.

#### Feature

##### Validation and Defaulting

Help users prevent errors and allow us to evolve the API independently to our
clients. These features are most useful when there are many clients who cannot
all update at the same time.

###### CRDs
Most validation can be specified in the CRD using OpenAPI v3.0 validation.
CRDValidationRatcheting feature gate allows failing validations specified using
OpenAPI also can be ignored if the failing part of the resource was unchanged.
Any other validations supported by addition of a Validating Webhook.

Defaulting is possible via OpenAPI v3.0 validation `default` keyword or via a
mutating webhook. This will not be run when reading from etcd for old object.

###### Aggregation APIs

Arbitration validation checks.

##### Multi-versioning

Allows serving the same object through two API versions. Can help ease API
changes like renaming fields. Less important if we control our client versions.
Available for both CRDs and Aggregation APIs.

##### Custom Business Logic

Perform arbitrary checks or actions when creating, reading, updating or deleting
an object. Possible for CRDs via webhooks.

##### Scale Subresource

Allow systems like HorizontalPodAutoscaler and PodDisruptionBudget interact with
new resource. Possible for both CRDs and Aggregation API. For CRDs, WebHooks
would be required.

##### Status Subresource

Allows fine-grained access control where user writes the spec section and the
controller writes the status section. Allows incrementing object Generation on
custom resource data mutation (requires separate spec and status sections in the
resource). Possible for both CRDs and Aggregation API.

##### Other Subresources

Add operations other than CRUD, such as `logs` or `exec`. Only available for Aggregation API.

##### Strategic Merge Patch

The new endpoints support `PATCH` request with `Content-Type:
application/strategic-merge-patch+json`. Useful for updating objects that may be
modified both locally, and by the server. Only available via Aggregation API.

##### Protocol Buffer

The new resource supports clients that want to use Protocol Buffers. Only
available via Aggregation API.

##### OpenAPI Schema

Is there an OpenAPI (swagger) schema for the types that can be dynamically
fetched from the server? Is the user protected from misspelling field names by
ensuring only allowed fields are set? Are types enforced(in other words, do not
put an `int` in a `string` field?). This feature available for both CRDs and
Aggregation API.

### Common features

When we cerate a custom resource, either via a CRD or an AA, we get many
features for the API, compared to implementing it outside the k8s platform.

#### CRUD
The new endpoints support CRUD basic operations via HTTP and kubectl.

#### Watch
The new endpoints support k8s watch operations via HTTP.

#### Discovery
Clients like kubectl and dashboard automatically offer list, display, and field
edit operations on the resources.

#### json-path
The new endpoints support PATCH request with `Content-Type:
application/json-path+json`.

#### merge-path
The new endpoints support PATCH request with `Content-Type:
application/merge-patch+json`.

#### HTTPS
The new endpoints uses HTTPS.

#### Built-in Authentication
Access to the extension uses the core API server (aggregation layer) for
authentication.

#### Built-in Authorization
Access to the extension can reuse the authorization used by the core API server;
for example, RBAC.

#### Finalizers
Block deletion of extension resources until external cleanup happens.

#### Admission Webhooks
Set default values and validate extension resources during any
create/update/delete operation.

#### UI/CLI Display
Kubectl, dashboard can display extension resources.

#### Unset versus Empty
Clients can distinguish unset fields from zero-valued fields.

#### Client Libraries Generationk8s provides generic client libraries, as well
as tools to generate type-specific client libraries.

#### Labels and annotations
Common metadata across objects that tools know how to edit for core and custom
resources.

## Preparing to Install Custom Resource

There are several points to be aware of before adding a custom resource to the
cluster.

### Third party code and new points of failure

While creating a CRD does not automatically add any new points of failure (for
example, by causing third party code to run on the API server), packages (for
example, Charts) or other installation bundles often include CRDs as well as
Deployment of third-party code that implements the business logic for a new
custom resource.

Installing an Aggregated API server always involves running a new Deployment.

#### Storage

Custom resources consume storage space in the same way that ConfigMaps do.
Creating too many custom resources may overload the API server's storage space.

Aggregated API servers may use the same storage as the main API server, in which
case the same warning applies.

#### Authentication, authorization, and auditing

CRDs always use the same authentication, authorization, and audit logging as the
built-in resources of the API server.

If we use RBAC for authorization, most RBAC roles will not grant access to the
new resources (except the cluster-admin role or any role created with wildcard
rules).

Aggregated API servers may or may not use the same authentication,
authorization, and auditing as the primary API server.

## Accessing a custom resource

k8s client libraries can be used to access custom resources. Not all client
libraries support custom resources. The Go and Python client libraries do.

When custom resource is added, we can access it using:
- kubectl
- k8s dynamic client
- REST client
- client generated using k8s client generation tools.

## Custom resource field selectors

Field selectors let clients select custom resources based on the value of one or
more resource fields.

All custom resources support the `metadata.name` and `metadata.namespace` field
selectors.

Field declared in a CustomResourceDefinition may also be used with field
selectors when included in the `spec.versions[*].selectableFields` field of the
CustomResourceDefinition.

### Selectable fields for custom resources

We need to enable the `CustomResourceFeildSelectors` feature gate to use this
behaviour, which then applies to all CustomResourceDefinition in the cluster.

The `spec.versions[*].selectableFields` field of a CustomResourceDefinition may
be used to declare which other fields in a custom resource may be used in field
selectors. The following example adds the `.spec.color` and `.spec.size` fields
as selectable fields.

```yaml
#shirt-resource-definition.yaml

apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
    name: shirts.stable.example.com
spec:
    group: stable.example.com
    scope: Namespaced
    names:
        plural: shirts
        singular: shirt
        kind: Shirt
    versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
            type: object
            properties:
                spec:
                    type: object
                    properties:
                        color:
                            type: string
                        size:
                            type: string
    selectableFields:
    - jsonPath: .spec.color
    - jsonPath: .spec.size
    additionalPrinterColumns:
    - jsonPath: .spec.color
      name: Color
      type: string
    - jsonPath: .spec.size
      name: Size
      type: string
```

Field selectors can then be used to get only resources with a `color` of `blue`:

```bash
kubectl get shirts.stable.example.com --field-selector spec.color=Blue
```

The output should be:

```bash
NAME       COLOR  SIZE
example1   blue   S
example2   blue   M
```
