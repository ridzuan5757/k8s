# Cloud Controller Manager

The cloud controller manager is a k8s control plane component that embeds cloud
specific control logic. It lets us link our clouster into our cloud provider's
API, and separates out the compoentns that interact with that cloud platform
from compoennts that only interact with our clouster.

By decoupling the inteoperability logic between k8s and the underlying cloud
infrastructure, the cloud controller manager compoentn enables cloud providers
to release features at a differnt pace compared to the main k8s project.

it is structured using a plugin meachnism that allow different cloud providers
to integrate the platforms with k8s.

The cloud controller manager runs in the control palne as a replicated set of
processes. Usually, these are containers in pods. Each coud controller manager
implements multiple controllers in a single process. 

We can also run the cloud controller manager as a k8s addon rather than as part
of the control plane.

## Cloud controller manager functions:

The controllers inside the cloud controller manager include:

### Node controller

The node controller is responsible for updating node objects when new servers
are created in our cloud infrastructure. The node controller obtains information
about the hosts running inside our tenancy with the cloud provider. The node
controller performs the following functions:
- Update a node project with the corrsponding server's unique identifier
  obtained form the cloud provider API.
- Annotating and labelling the node object with cloud specific informatioon,
  such as the region the node is deployed into and the resources such as CPU and
  memory that it has available.
- Obtain the node's hostname and network addresses.
- Verifying the node's health. In case a node becomes unresponsive, this
  controller checks with the cloud provider's API to see if the server has been
  deactivated / deleted / terminated. If the node has been deleted from the
  cloud, the controller deletes the node object from the cluster.

Some cloud provider implementations split this into a ndoe controller and a
separate node lifecycle controller.
