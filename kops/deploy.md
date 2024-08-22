create:

staging

```bash
kops create cluster --name shell.ronpos.com --state=s3://monitoring-state-store --node-count 5 --control-plane-count 3 --control-plane-size t2.medium --node-size t2.medium --control-plane-zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --ssh-public-key ./proactive-monitoring.pub --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --network-id vpc-06aeb8d9751af25ce --topology private --subnets subnet-05210b7fd982a905b,subnet-0f1e84404c2c72470,subnet-06e8534ee6afc7c17 --utility-subnets subnet-07d9682ff52eb7371,subnet-053773833538d9aca,subnet-080800dd77492ca92 --bastion --dns private --dns-zone Z08353583QWN5RFKKJD4
```

shell-private

```bash
kops create cluster --state=s3://monitoring-state-store --node-count 5 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-1a --zones ap-southeast-1a --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --topology private --bastion --network-id vpc-064d56cdcb000f690 --dns-zone Z03730671I4OKR1B7EROZ --dns private --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery
```

shell-private-v2

```bash
kops create cluster --state=s3://monitoring-state-store --node-count 5 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --topology private --bastion --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery --network-id vpc-0d1e5a5a7d9c0fa91 --subnets subnet-031c364c869689ebf,subnet-05544cb3194734fa4,subnet-09a556aca9f26282c --utility-subnets subnet-031c364c869689ebf,subnet-05544cb3194734fa4,subnet-09a556aca9f26282c --dns-zone Z07725001JC0ZTQQYE2IE --dns private
```

shell-public
```bash
kops create cluster --state=s3://monitoring-state-store --node-count 5 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --zones ap-southeast-1a,ap-southeast-1b,ap-southeast-1c --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery 
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
kops create ig nodes-ap-southeast-1a-opensearch --name proactivemonitoring.silentmode.com --role node --state s3://proactive-monitoring-state

update:
kops edit cluster proactivemonitoring.silentmode.com --state=s3://proactive-monitoring-state

delete:
kops delete cluster  proactivemonitoring.silentmode.com --yes --state s3://proactive-monitoring-state
