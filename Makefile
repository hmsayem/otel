JAEGER_OPERATOR_VERSION = v1.36.0

namespace:
	kubectl apply -f k8s/namespace.yaml

jaeger-operator:
	# Create the jaeger operator and necessary artifacts in ns observability
	kubectl create -n observability -f https://github.com/jaegertracing/jaeger-operator/releases/download/$(JAEGER_OPERATOR_VERSION)/jaeger-operator.yaml

jaeger:
	kubectl apply -f k8s/jaeger.yaml

prometheus:
	kubectl apply -f k8s/prometheus-service.yaml   # Prometheus instance
	kubectl apply -f k8s/prometheus-monitor.yaml   # Service monitor

otel-collector:
	kubectl apply -f k8s/otel-collector.yaml

clean:
	- kubectl delete -f k8s/otel-collector.yaml

	- kubectl delete -f k8s/prometheus-monitor.yaml
	- kubectl delete -f k8s/prometheus-service.yaml

	- kubectl delete -f k8s/jaeger.yaml

	- kubectl delete -n observability -f https://github.com/jaegertracing/jaeger-operator/releases/download/$(JAEGER_OPERATOR_VERSION)/jaeger-operator.yaml
