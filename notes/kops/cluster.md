# Cluster State Store

In order to store the state of the cluster, and the representation of the
cluster, a dedicated `S3` bucket for `kops` will be created. This bucket will be
the source of truth for the cluster configuration.

It is always a good practice to enable versioning because we are storing the
state of the cluster because if something goes wrong, we can go to the s3
bucket and revert the state. 

```bash
bucket-name: proactive-monitoring-state
region: ap-southeast-1
```

# Cluster OIDC store

In order for `ServiceAccounts` to use external permissions (IAM Roles for
`ServiceAccounts`), we also need a bucket for hosting the OIDC documents. While
we can reuse the bucket for the cluster state if we grant it public ACL, a
separate bucket is recommended.

The ACL must be public so that the AWS STS service can access them.

```bash
aws s3api create-bucket \
    --bucket prefix-example-com-oidc-store \
    --region us-east-1 \
    --object-ownership BucketOwnerPreferred
aws s3api put-public-access-block \
    --bucket prefix-example-com-oidc-store \
    --public-access-block-configuration BlockPublicAcls=false,IgnorePublicAcls=false,BlockPublicPolicy=false,RestrictPublicBuckets=false
aws s3api put-bucket-acl \
    --bucket prefix-example-com-oidc-store \
    --acl public-read
```

# Creating Cluster

## Prepare local environment

We are ready to start creating the first cluster. The following environment
variables are recommended to make the process easier.

```bash
export NAME=proactive.monitoring.example.com
export KOPS_STATE_STORE=s3://proactive-monitoring-state
```

For gossip-based cluster, make sure the name ends with `k8s.local`. For example:

> [!NOTE]
> Gossip-based cluster use a peer-to-peer network instead of externally hosted
> DNS for propagating the `k8s` API address. This means that an externally
> hosted DNS service is not needed.
>
> Gossip does not suffer potential disruptions due to out of date records in DNS
> caches as the propagation is almost instant.
>
> Gossip is also the only option if we want to deploy a cluster in any of the
> AWS regions without Route 53.


```bash
export NAME=myfirstcluster.k8s.local
export KOPS_STATE_STORE=s3://prefix-example-com-state-store
```

Environment variables is not mandatory. We can always define the values using
the `--name` and `--state` flags later.

## Create cluster configuration

We will need to note which availability zones are available to us. In example we
will be deploying the cluster to the `ap-southeast-1` region.

```bash
aws ec2 describe-availability-zones --region ap-southeast-1
```

Below is a create cluster command. We will use the most basic example possible,
with more verbose examples in high availability. The below command will generate
a cluster configuration, but will not start building it. Make sure an SSH key
pair have been generated before the cluster is created.

```bash
kops create cluster \
    --name=${NAME}  \
    --cloud=aws \
    --zones=ap-southeast-1 \
    --discovery-store=s3://prefix-example-com-oidc-store/${NAME}/discovery
```

All instances created by `kops` will be built within auto scaling groups, which
means each instance will be automatically monitored and rebuilt by AWS if it
suffers any failure.

## Customize Cluster Configuration

Now we have a cluster configuration, we can look at every aspect that defines
the cluster by editing the description.

```bash
kops edit cluster --name ${NAME}
```

This opes the editor `$EDITOR` and allows us to edit the configuration. The
configuration is loaded from the S3 bucket we created earlier and automatically
updated when we save and exit the editor.

We will leave everyhting to set to the defaults for now, the the rest of `kops`
documentation covers additional settings and configuration we can enable.

## Building Cluster

Now we take the final step of actually building the cluster. This will take a
while. Once it finishes we will have to wait longer while the booted instances
finish downloading `k8s` components and reach `ready` state.

```bash
kops update cluster --name ${NAME} --yes --admin
```

## Using Cluster

The configuration for the cluster will be automatically generated and written to
`~/.kube/config`. A simple `k8s` API call can be used to check if the API is
online and listening.

```bash
kubectl get nodes
```

We will see a list of nodes that should match the `--zones` flag defined
earlier. This is a sign that the `k8s` is online and working. `kops` also ships
with a handy validation tool that can be ran to ensure the cluster is working as
expected.

```bash
kops validate cluster --wait 10m
```

We can look at all the system components with the following command:

```bash
kubectl -n kube-system get all
```

## Deleting Cluster

Running a `k8s` cluster within AWS can be costly. So we may want to delete the
cluster once the experiment is completed. We cam preview all of the  AWS
resources that will be destroyed when the cluster is deleted by issueing the
following command.

```bash
kops delete cluster --name ${NAME}
```

When we are sure that we want to delete the cluster, issue the delete command
with the `--yes` flag. ***Note that this command is very destructive,a dnw ill
delete the cluster and everything contained within it.***
