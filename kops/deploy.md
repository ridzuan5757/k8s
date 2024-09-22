create:

staging

```bash
kops create cluster --name monitoring.staging.ronpos.com --state=s3://monitoring-state-store --node-count 3 --control-plane-count 3 --control-plane-size t3.medium --node-size c7i.xlarge	--control-plane-zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --ssh-public-key ./proactive-monitoring.pub --cloud-labels "silentmode:owner=engineering, silentmode:environment=ronpos-staging, silentmode:service=cluster-default" --network-id vpc-04d79ee4ad78189d5 --topology private --subnets subnet-098fa6f36930e1444,subnet-044b686495b3e8b54,subnet-0fef5199d3ae5f0ae --utility-subnets subnet-0fa644f0c5e8250c2,subnet-0ba0dc08b27ad8ac5,subnet-0f20762800bf2feb0 --bastion --dns-zone Z029400221WKTP44KDW6I
```

```bash
kops create cluster --name shell.ronpos.com --state s3://monitoring-state-store --node-count 1 --control-plane-count 1 --control-plane-size t3.medium --node-size t3.medium --zones ap-southeast-5a --control-plane-zones ap-southeast-5a --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --networking amazonvpc
```

kops create cluster --name shell.ronpos.com --state s3://monitoring-state-store --node-count 1 --control-plane-count 1 --control-plane-size t3.medium --node-size t3.medium --zones ap-southeast-5a --control-plane-zones ap-southeast-5a --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --networking amazonvpc

shell-private

```bash
kops create cluster --state=s3://monitoring-state-store --node-count 5 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-1a --zones ap-southeast-1a --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --topology private --bastion --network-id vpc-064d56cdcb000f690 --dns-zone Z03730671I4OKR1B7EROZ --dns private --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery
```

shell-private-v2

```bash
kops create cluster --state=s3://monitoring-state-store --node-count 5 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-5a,ap-southeast-5b,ap-southeast-5c --zones ap-southeast-5a,ap-southeast-5b,ap-southeast-5c --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --topology private --bastion --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery --network-id vpc-08a65d72ec4e0c6bf --subnets subnet-0557f73330b0bfab8,subnet-0e9f14b076fe6ca80,subnet-06230f3f8938144ad --utility-subnets subnet-055547bcc1d34e79b,subnet-02737750921e602c4,subnet-03391b225a0de9ae0 --dns-zone Z03687583RGJ2ACMJLDCW --dns private --networking amazonvpc
```

```bash
kops create cluster --state=s3://monitoring-state-store --node-count 3 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-5a,ap-southeast-5b,ap-southeast-5c --zones ap-southeast-5a,ap-southeast-5b,ap-southeast-5c --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --topology private --bastion --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery 
```


shell-public
```bash
kops create cluster --state=s3://monitoring-state-store --node-count 5 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-5a,ap-southeast-5b,ap-southeast-5c --zones ap-southeast-5a,ap-southeast-5b,ap-southeast-5c --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery --networking amazonvpc --network-id vpc-08a65d72ec4e0c6bf --subnets subnet-055547bcc1d34e79b,subnet-02737750921e602c4,subnet-03391b225a0de9ae0
```

```bash
kops create cluster --state=s3://monitoring-state-store --node-count 5 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-5a,ap-southeast-5b,ap-southeast-5c --zones ap-southeast-5a,ap-southeast-5b,ap-southeast-5c --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery --networking amazonvpc --network-id vpc-08a65d72ec4e0c6bf --topology private --subnets subnet-0557f73330b0bfab8,subnet-0e9f14b076fe6ca80,subnet-06230f3f8938144ad --utility-subnets subnet-055547bcc1d34e79b,subnet-02737750921e602c4,subnet-03391b225a0de9ae0 --dns private --dns-zone Z03687583RGJ2ACMJLDCW --bastion
```

```bash
kops create cluster --state=s3://monitoring-state-store --node-count 5 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery --networking amazonvpc --network-id vpc-04d79ee4ad78189d5 --topology private --subnets subnet-098fa6f36930e1444,subnet-044b686495b3e8b54,subnet-0fef5199d3ae5f0ae --utility-subnets subnet-0fa644f0c5e8250c2,subnet-0ba0dc08b27ad8ac5,subnet-0f20762800bf2feb0 --bastion --dns public --dns-zone Z02935632M8VYG5TTFKL3
```

```bash
kops create cluster --state=s3://monitoring-state-store --node-count 3 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery --network-id vpc-04d79ee4ad78189d5 --topology private --subnets subnet-098fa6f36930e1444,subnet-044b686495b3e8b54,subnet-0fef5199d3ae5f0ae --utility-subnets subnet-0fa644f0c5e8250c2,subnet-0ba0dc08b27ad8ac5,subnet-0f20762800bf2feb0 --bastion
```

edit:

```bash
kops edit cluster proactivemonitoring.silentmode.com --state=s3://proactive-monitoring-state
```

```yaml
cloudLabels:
    silentmode:environment: ronpos-staging
    silentmode:owner: engineering
    silentmode:service: proactive-monitoring
awsLoadBalancerController:
    enabled: true
certManager:
    enabled: true
```

add instance group:
```bash
kops create ig nodes-ap-southeast-1a-opensearch --name proactivemonitoring.silentmode.com --role node --state s3://proactive-monitoring-state
```

```yaml
apiVersion: kops.k8s.io/v1alpha2
kind: InstanceGroup
metadata:
  creationTimestamp: "2024-08-30T20:47:26Z"
  generation: 1
  labels:
    kops.k8s.io/cluster: shell.ronpos.com
  name: nodes-ap-southeast-1b
spec:
  image: 099720109477/ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20240607
  machineType: t3.large
  maxSize: 1
  minSize: 1
  nodeLabels:
    node: cluster-config
  role: Node
  subnets:
  - ap-southeast-1b
```

update:
```bash
kops edit cluster proactivemonitoring.silentmode.com --state=s3://proactive-monitoring-state
```

delete:
```bash
kops delete cluster  proactivemonitoring.silentmode.com --yes --state s3://proactive-monitoring-state
```