create:
kops create cluster --state=s3://proactive-monitoring-state --node-count 3 --control-plane-count 3 --control-plane-size t2.medium --node-size t2.medium --control-plane-zones ap-southeast-1a --zones=ap-southeast-1a --name proactivemonitoring.silentmode.com

edit:
kops edit cluster proactivemonitoring.silentmode.com --state=s3://proactive-monitoring-state
silentmode:environment	    ronpos-staging
silentmode:owner	        engineering
silentmode:service	        proactive-monitoring

update:
kops edit cluster proactivemonitoring.silentmode.com --state=s3://proactive-monitoring-state

delete:
kops delete cluster  proactivemonitoring.silentmode.com --yes --state s3://proactive-monitoring-state
