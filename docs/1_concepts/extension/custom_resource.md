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
- We wnat to put the entire configuration into one key of a ConfigMap.
-
