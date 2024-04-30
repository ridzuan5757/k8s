from diagrams import Diagram, Cluster, Edge
from diagrams.custom import Custom
from diagrams.onprem.container import Docker, Containerd
from diagrams.onprem.monitoring import Grafana, Prometheus
from diagrams.onprem.tracing import Jaeger
from diagrams.onprem.logging import Rsyslog
from diagrams.onprem.network import Nginx
from diagrams.programming.language import Nodejs
from diagrams.aws.compute import EC2
from diagrams.k8s.infra import Master, Node, ETCD
from diagrams.k8s.controlplane import APIServer, ControllerManager, KubeProxy, Kubelet, Scheduler
from diagrams.aws.storage import SimpleStorageServiceS3
from diagrams.aws.integration import SimpleQueueServiceSqs

with Diagram("Station", show=True, direction="BT"):
    with Cluster("Amazon Web Services"):
        cc_s3 = SimpleStorageServiceS3("S3")
        cc_sqs = SimpleQueueServiceSqs("SQS")
        cc_s3 - cc_sqs

    with Cluster("Ronpos Hub"):
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

    Custom("POS", "./images/ronpos.png") >> cc_otel_agent
    Custom("CDS", "./images/ronpos.png") >> cc_otel_agent

    with Cluster("EC2 Instance - Master Node 1") as mn1:
        EC2("")
        Containerd("") - Master("Master Node 1")
        with Cluster(""):
            cc_api_m1 = APIServer("")
            ControllerManager("") >> cc_api_m1
            ETCD("") >> cc_api_m1
            Scheduler("") >> cc_api_m1

    with Cluster("EC2 Instance - Master Node 2") as mn2:
        EC2("")
        Containerd("") - Master("Master Node 2")
        with Cluster(""):
            cc_api_m2 = APIServer("")
            ControllerManager("") >> cc_api_m2
            ETCD("") >> cc_api_m2
            Scheduler("") >> cc_api_m2

    cc_api_m1 - cc_api_m2

    with Cluster("EC2 Instance - Load Balancer", direction="LR"):
        EC2("")
        Containerd("") - Node("Load Balancer Node")
        with Cluster(""):
            cc_kubelet_lb = Kubelet("")
            cc_kubelet_lb >> cc_api_m1
            cc_kubelet_lb >> cc_api_m2
            cc_kubeproxy_lb = KubeProxy("")
            cc_kubeproxy_lb >> cc_api_m1
            cc_kubeproxy_lb >> cc_api_m2
            cc_kubeproxy_lb - cc_kubelet_lb
            with Cluster(""):
                cc_lb = Custom("MetalLB", "./images/metallb.png")
                cc_ingress = Nginx("Ingress")
                cc_lb >> cc_ingress

    with Cluster("EC2 Instance - Worker Node 1", direction="LR"):
        EC2("")
        Containerd("") - Node("Worker Node 1")
        with Cluster(""):
            cc_kubelet_w1 = Kubelet()
            cc_kubelet_w1 >> cc_api_m1
            cc_kubelet_w1 >> cc_api_m2
            cc_kubeproxy_w1 = KubeProxy()
            cc_kubeproxy_w1 >> cc_api_m1
            cc_kubeproxy_w1 >> cc_api_m2
            cc_kubeproxy_w1 - cc_kubelet_w1
            with Cluster(""):
                cc_otel_w1 = Custom(
                    "OTel - Gateway", "./images/opentelemetry.png")
                cc_dataprepper_w1 = Custom(
                    "DataPrepper", "./images/dataprepper.png")
                cc_opensearch_w1 = Custom(
                    "OpenSearch", "./images/opensearch.png")
                cc_dashboard_w1 = Custom("OpenSearch\n Dashboard",
                                         "./images/opensearch_dashboard.png")
                cc_fluentbit_w1 = Custom("FluentBit", "./images/fluentbit.png")
                cc_otel_w1 >> Edge(label="trace & metric") >> cc_dataprepper_w1
                cc_dataprepper_w1 >> cc_opensearch_w1 >> cc_dashboard_w1
                cc_otel_w1 >> Edge(label="log") >> cc_fluentbit_w1
                cc_fluentbit_w1 >> cc_dataprepper_w1

    with Cluster("EC2 Instance - Worker Node 2", direction="LR"):
        EC2()
        Containerd() - Node("Worker Node 2")
        with Cluster():
            cc_kubelet_w2 = Kubelet()
            cc_kubelet_w2 >> cc_api_m1
            cc_kubelet_w2 >> cc_api_m2
            cc_kubeproxy_w2 = KubeProxy()
            cc_kubeproxy_w2 >> cc_api_m1
            cc_kubeproxy_w2 >> cc_api_m2
            cc_kubelet_w2 - cc_kubeproxy_w2
            with Cluster(""):
                cc_otel_w2 = Custom(
                    "OTel - Gateway", "./images/opentelemetry.png")
                cc_dataprepper_w2 = Custom(
                    "DataPrepper", "./images/dataprepper.png")
                cc_opensearch_w2 = Custom(
                    "OpenSearch", "./images/opensearch.png")
                cc_dashboard_w2 = Custom("OpenSearch\n Dashboard",
                                         "./images/opensearch_dashboard.png")
                cc_fluentbit_w2 = Custom("FluentBit", "./images/fluentbit.png")
                cc_otel_w2 >> Edge(label="trace & metric") >> cc_dataprepper_w2
                cc_dataprepper_w2 >> cc_opensearch_w2 >> cc_dashboard_w2
                cc_fluentbit_w2 >> cc_dataprepper_w2
                cc_otel_w2 >> Edge(label="log") >> cc_fluentbit_w2

    cc_otel_agent >> cc_lb
    cc_dataprepper_w1 >> cc_s3
    cc_dataprepper_w2 >> cc_s3
