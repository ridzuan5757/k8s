create:

```bash
kops create cluster --state=s3://proactive-monitoring-state --node-count 2 --control-plane-count 1 --control-plane-size t2.medium --node-size t2.medium --control-plane-zones ap-southeast-1a --zones ap-southeast-1a --name proactivemonitoring.silentmode.com --ssh-public-key ./proactive-monitoring.pub
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
