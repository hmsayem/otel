### Install Zookeeper

```bash
helm install zookeeper \
    --set zookeeper.enabled=false \
    --set persistence.size=100Mi \
    --set replicaCount=1 \
    oci://registry-1.docker.io/bitnamicharts/zookeeper
```


### Install Clickhouse Operator

```bash
helm repo add clickhouse-operator https://docs.altinity.com/clickhouse-operator/
helm repo update
helm install clickhouse-operator clickhouse-operator/altinity-clickhouse-operator
```

### Install OpenTelemetry Collector

```bash
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update
helm install otel-collector open-telemetry/opentelemetry-collector  --values ./values.yaml
```