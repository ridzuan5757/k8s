create:

staging

```bash
kops create cluster --state=s3://monitoring-state-store --node-count 2 --control-plane-count 1 --control-plane-size t2.medium --node-size t2.medium --control-plane-zones ap-southeast-1a --zones ap-southeast-1a --name shell.ronpos.com --ssh-public-key ./proactive-monitoring.pub --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --vpc vpc-04b47bd44a664c3a6
```

shell-private

```bash
kops create cluster --state=s3://monitoring-state-store --node-count 5 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-1a --zones ap-southeast-1a --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --topology private --bastion --network-id vpc-064d56cdcb000f690 --dns-zone Z03730671I4OKR1B7EROZ --dns private --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery
```

shell-private-v2

```bash
kops create cluster --state=s3://monitoring-state-store --node-count 5 --control-plane-count 3 --control-plane-size t3.medium --node-size t3.large --control-plane-zones ap-southeast-1a --zones ap-southeast-1a --name shell.ronpos.com --ssh-public-key ./monitoring-shell.pub --topology private --bastion --network-id vpc-064d56cdcb000f690 --dns-zone Z03730671I4OKR1B7EROZ --dns private --cloud-labels "silentmode:owner=engineering, silentmode:environment=shell-production, silentmode:service=cluster-default" --discovery-store s3://monitoring-oidc-store/shell.ronpos.com/discovery
```

shell-public


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
