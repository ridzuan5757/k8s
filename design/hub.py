from diagrams import Diagram, Cluster
from diagrams.custom import Custom
from diagrams.onprem.container import Docker
from diagrams.onprem.monitoring import Grafana, Prometheus
from diagrams.onprem.tracing import Jaeger
from diagrams.onprem.logging import Rsyslog
from diagrams.programming.language import Nodejs

with Diagram("Ronpos Hub", show=True, direction="LR"):

    with Cluster("Hub Container Runtime", direction="LR"):
        cc_ronpos = Custom("hub", "./images/ronpos.png")
        cc_hub_tracing = Nodejs("hub-tracing")
        cc_docker_hub = Docker("")

    with Cluster("Observer Container Runtime", direction="LR"):
        cc_otel_agent = Custom("OTel Collector Agent",
                               "./images/opentelemetry.png")
        cc_docker_observer = Docker("")
        cc_prom = Prometheus("Prometheus")
        cc_grafana = Grafana("Grafana")
        cc_jaeger = Jaeger("Jaeger")

        cc_otel_agent >> cc_prom >> cc_grafana
        cc_otel_agent >> cc_jaeger

    cc_hub_tracing >> cc_otel_agent

    Rsyslog("RSyslog") >> cc_otel_agent
