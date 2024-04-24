# ReplicaSet

A ReplicaSet's purpose is to maintain a stable set of replica Pods running at
any given time. As such, it is often used to guarantee the availability of a
specified number of identical Pods.

## Mechanism

A ReplicaSet is defined with fields, including a selector that specifies how to
identify Pods it can acquire, a number of replicas indicating how many Pods it
should be maintaining, and a pod template specifying the data of new Pods it
should create to meet the number of replicas criteria.

A ReplicaSet then fulfills its purpose by creating and deleting Pods as needed
to reach the desired number. When a ReplicaSet needs to create new Pods, it uses
its Pod template.

A ReplicaSet is linked to its Pods via the Pods' `metadata.ownerReferences`
field, which specifies what resource the current object is owned by. All Pods
acquired by a ReplicaSet have their owning ReplicaSet's identifying information
withinb their ownerReferences field. It is through this link that the ReplicaSet
knows the state of the Pods it is maintaining and plans accordingly.

A ReplicaSet identifies new Pods to acquire by using its selector. If there is a
Pod that has no OwnerReference or the OwnerReference is not a Controller and it
matches a ReplicaSet's selector, it will be immediately acquired by said
ReplicaSet.
