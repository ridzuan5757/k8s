# Workload Management

k8s provides several built-in APIs for declarative management of the workloads
and the components of those workloads.

Ultimately, the applications running as containers inside pods. However,
managing individual pods would requires lot of effort. For example, if a pod
fails, we probably want to run a new pod to replace it via k8s.

We use the k8s API to create the workload object that represents a higher
abstraction level than a pod, and then the k8s control plane automatically
manages pod objects on our behalf, based on the specification for the workload
object that has been define using the manifest YAML file.


