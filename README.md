# Node Metrics

This project allows to retrieve the total memory available and used form a node
and export it as a prometheus metric.

# Usage

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    k8s-app: node-metrics
  name: metrics
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: node-metrics
  template:
    metadata:
      annotations:
        prometheus.io/port: "9093"
        prometheus.io/scrape: "true"
        prometheus.io/path: "/metrics"
      labels:
        k8s-app: node-metrics
    spec:
      containers:
      - image: "quay.io/aanm/node-metrics:v0.0.1"
        imagePullPolicy: IfNotPresent
        env:
        - name: PROMETHEUS_ADDR
          value: "0.0.0.0:9093"
        name: node-metrics
        ports:                
        - containerPort: 9093
          hostPort: 9093
          name: prometheus
          protocol: TCP
      hostNetwork: true
      restartPolicy: Always
      priorityClassName: system-node-critical
      tolerations:
      - operator: Exists
```