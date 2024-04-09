# Controllers

A control loop is a non-terminating loop that regulates the state of the system.
In k8s, controllers are control loops that watch the state of the cluster, then
make or request changes where needed. Each controller tries to move the current
cluster state closer to the desired state.

## Controller Pattern

A controller tracks at least on k8s resource type. These objects have a spec
field that represents the desired state. The controllers for that resource are
responsible for making the current state some closer to that desired state.

The controller might carry the action out itself; more commonly, in k8s , a
controller will send messages to the API server that have useful side effects.

### Control via API server

The `Job` controller is an example of a k8s built-in controller. Built-in
controllers manage state by interacting with the cluster API server.

Job is a k8s resource that runs a Pod, or perhaps serveral Pods, to carry out
task and then stop. Once scheduled, Pod object become part of the desired state
of the kubelet.

When the Job controller sees a new task it make sure that, somewhere in our
cluster, the kubelets on a set of nodes are running the right number of pods to
get the work done. The job controller tells the API server to create or remove
pods. Other components in the control plane act on the new information (there
are new pods to schedule and run), and eventually the work is done.

After we create a new job, the desired state is for that job to be completed.
The job controller makes the current state for that job be nearer to our desired
state - creating pods that do the work we wanted for that job, so that the job
is closer to completion.

Controllers also update the objects that configure them. For example, once the
work is done for a job, the job controller updates that the job object to mark
it as `Finished`.

### Direct control

In contrast with job, some controllers need to make changes to things outside of
out cluster. For example, if we use a control loop to make sure there are enough
nodes in our cluster, then that controller needs something outside the current
cluster to set up new nodes when needed.

Controllers that interact with external state find their desired state from the
API server, then communicate directly with an external system to bring the
current state closer in line. There actually is a controller that horizontally
scales the nodes in our cluster.

The important point here is that the controller makes some changes to bring
about our desired state, and then reports the current state back to our
cluster's API server. Other control loops can observe that reported data and
take their own actions.

With k8s clusters, the control plane indirectly works with IP address management
tools, storage services, cloud provider APIs, and other services by extending
k8s to implement that.

## Desired versus current state

k8s takes a cloud-native view of systems, and is able to handle constant change.
Our cluster could be changing at any point as work happens and control loops
automatically fix failures. This means that, potentially, our cluster never
reaches a stable state.

As long as the controllers for our cluster are running and able to make useful
changes, it does not matter if the overall state is stable or not.

## Design 

as a tenet of its design, k8s uses lots of controllers that each manage
particular aspect of cluster state. Most commonly, a particular control loop
(controller) uses one kind of resource as it desired state, and has a different
kind of resource that it manages to make that desired state happen.

For example, a controller for jobs tracks job object (to discover new work) and
pod objects (to run the jobs, and then to see when the work is finished). In
this case something else creates the jobs, whereas the job controller creates
pod.

It is useful to have simple controllers rather than one, monolithic set of
control loops that are interlinked. Controller can fail, so k8s is designed to
allow for that.

There can be serveral controllers that create or update the same kind of object.
Behind the scenes, k8s controllers make sure that they only pay attention to the
resources linked to their controlling resource.

For example, we can have deployments and jobs; these both creates pods. The job
controller does not delete the pods that our deployment created, because there
is information (labels) the controllers can use to tell those pods apart.

## Way of running controllers

k8s comes with a set of built-in controllers that run inside
`kube-controller-manager`. These built-in controllers provide important core
behaviours. 

The deployment controller and job controller are examples of controllers that
come as part of k8s itself (built-in controllers). k8s lets us run a resilient
control plane, so that if any of the built-in controllers were to fail, another
part of the control-plane will take over the work.

We can find controllers that run outside the control plane, to extend k8s. Or,
if we want, we can write a new controller ourselves. We can run our own
controller as a set of pods, or externally to k8s. What first best will depend
on what that particular controller does.
