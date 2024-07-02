install:
	helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts

deploy:
	helm install opentelemetry-collector-daemonset open-telemetry/opentelemetry-collector \
   --set image.repository="otel/opentelemetry-collector-k8s" \
   --set mode=daemonset \

delete:
	helm uninstall opentelemetry-collector-daemonset
