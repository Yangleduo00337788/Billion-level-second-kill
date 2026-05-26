#!/bin/bash
set -euo pipefail

#=======================================
# Seckill System K8s Deploy Script
#=======================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="${SCRIPT_DIR}/.."
K8S_DIR="${ROOT_DIR}/k8s"
DEPLOY_ENV="${DEPLOY_ENV:-prod}"
NAMESPACE="seckill"
MONITORING_NS="monitoring"
LOGGING_NS="logging"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC}  $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }
log_step()  { echo -e "\n${BLUE}==>${NC} $*"; }

check_prerequisites() {
    log_step "Checking prerequisites..."

    local required_cmds=("kubectl" "docker" "helm")
    local missing=()

    for cmd in "${required_cmds[@]}"; do
        if ! command -v "$cmd" &>/dev/null; then
            missing+=("$cmd")
        fi
    done

    if [ ${#missing[@]} -ne 0 ]; then
        log_error "Missing required tools: ${missing[*]}"
        log_error "Please install them first."
        exit 1
    fi

    if ! kubectl cluster-info &>/dev/null; then
        log_error "Cannot connect to Kubernetes cluster. Check kubeconfig."
        exit 1
    fi

    log_info "All prerequisites satisfied"
}

deploy_namespaces() {
    log_step "Creating namespaces..."
    kubectl apply -f "${K8S_DIR}/namespace.yaml"
    log_info "Namespaces created"
}

deploy_secrets() {
    log_step "Deploying secrets..."
    kubectl apply -f "${K8S_DIR}/secret.yaml"
    log_info "Secrets deployed"
}

deploy_configmaps() {
    log_step "Deploying ConfigMaps..."
    kubectl apply -f "${K8S_DIR}/configmap.yaml"
    log_info "ConfigMaps deployed"
}

deploy_infrastructure() {
    log_step "Deploying infrastructure (MySQL, Redis, Kafka)..."

    helm repo add bitnami https://charts.bitnami.com/bitnami 2>/dev/null || true
    helm repo update

    # Redis Cluster
    if ! kubectl get statefulset redis -n "${NAMESPACE}" &>/dev/null; then
        helm upgrade --install redis bitnami/redis \
            --namespace "${NAMESPACE}" \
            --set architecture=standalone \
            --set auth.enabled=true \
            --set auth.existingSecret=seckill-redis-secret \
            --set auth.existingSecretPasswordKey=redis-password \
            --set master.persistence.size=100Gi \
            --set master.persistence.storageClass=ssd \
            --set master.resources.requests.cpu=1 \
            --set master.resources.requests.memory=2Gi \
            --set master.resources.limits.cpu=2 \
            --set master.resources.limits.memory=6Gi \
            --set master.disableCommands="" \
            --set master.extraFlags[0]="maxmemory-policy allkeys-lru" \
            --wait --timeout 10m
    else
        log_warn "Redis already deployed, skipping"
    fi

    # MySQL
    if ! kubectl get statefulset mysql -n "${NAMESPACE}" &>/dev/null; then
        helm upgrade --install mysql bitnami/mysql \
            --namespace "${NAMESPACE}" \
            --set auth.rootPasswordSecretKey=mysql-root-password \
            --set auth.existingSecret=seckill-db-secret \
            --set auth.database=user_system \
            --set auth.username=seckill \
            --set auth.passwordSecretKey=mysql-password \
            --set primary.persistence.size=500Gi \
            --set primary.persistence.storageClass=ssd \
            --set primary.resources.requests.cpu=2 \
            --set primary.resources.requests.memory=4Gi \
            --set primary.resources.limits.cpu=4 \
            --set primary.resources.limits.memory=8Gi \
            --set primary.configuration="[mysqld]\nmax_connections=2000\ninnodb_buffer_pool_size=4G" \
            --wait --timeout 10m
    else
        log_warn "MySQL already deployed, skipping"
    fi

    # Kafka
    if ! kubectl get statefulset kafka -n "${NAMESPACE}" &>/dev/null; then
        helm upgrade --install kafka bitnami/kafka \
            --namespace "${NAMESPACE}" \
            --set replicaCount=3 \
            --set persistence.size=200Gi \
            --set persistence.storageClass=ssd \
            --set resources.requests.cpu=1 \
            --set resources.requests.memory=2Gi \
            --set resources.limits.cpu=2 \
            --set resources.limits.memory=4Gi \
            --set defaultReplicationFactor=3 \
            --set offsetsTopicReplicationFactor=3 \
            --set numPartitions=16 \
            --set logRetentionHours=168 \
            --set maxMessageBytes=10485760 \
            --wait --timeout 15m
    else
        log_warn "Kafka already deployed, skipping"
    fi

    log_info "Infrastructure deployed"
}

deploy_app() {
    log_step "Deploying seckill application..."

    kubectl apply -f "${K8S_DIR}/deployment.yaml"
    kubectl apply -f "${K8S_DIR}/service.yaml"
    kubectl apply -f "${K8S_DIR}/hpa.yaml"
    kubectl apply -f "${K8S_DIR}/ingress.yaml"

    kubectl rollout status deployment/seckill-app -n "${NAMESPACE}" --timeout=10m

    log_info "Application deployed successfully"
}

deploy_monitoring() {
    log_step "Deploying monitoring stack..."

    kubectl create namespace "${MONITORING_NS}" --dry-run=client -o yaml | kubectl apply -f -

    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts 2>/dev/null || true
    helm repo add grafana https://grafana.github.io/helm-charts 2>/dev/null || true
    helm repo update

    # Prometheus Stack (includes AlertManager, Grafana, NodeExporter)
    if ! kubectl get deployment prometheus-server -n "${MONITORING_NS}" &>/dev/null; then
        helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
            --namespace "${MONITORING_NS}" \
            --set prometheus.prometheusSpec.retention=30d \
            --set prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.storageClassName=ssd \
            --set prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.accessModes[0]=ReadWriteOnce \
            --set prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.resources.requests.storage=200Gi \
            --set prometheus.prometheusSpec.resources.requests.cpu=1 \
            --set prometheus.prometheusSpec.resources.requests.memory=4Gi \
            --set prometheus.prometheusSpec.resources.limits.cpu=2 \
            --set prometheus.prometheusSpec.resources.limits.memory=8Gi \
            --set prometheus.prometheusSpec.additionalScrapeConfigs="$(cat "${ROOT_DIR}/monitoring/prometheus/prometheus.yml" | sed 's/^/    /')" \
            --set grafana.adminPassword="${GRAFANA_PASSWORD:-admin123}" \
            --set grafana.persistence.enabled=true \
            --set grafana.persistence.size=50Gi \
            --set grafana.persistence.storageClassName=ssd \
            --set grafana."grafana\.ini".server.root_url="%(protocol)s://%(domain)s:%(http_port)s/grafana" \
            --set grafana."grafana\.ini".server.serve_from_sub_path=true \
            --set alertmanager.alertmanagerSpec.storage.volumeClaimTemplate.spec.storageClassName=ssd \
            --set alertmanager.alertmanagerSpec.storage.volumeClaimTemplate.spec.resources.requests.storage=10Gi \
            --wait --timeout 10m
    else
        log_warn "Prometheus stack already deployed, skipping"
    fi

    log_info "Monitoring stack deployed"
}

deploy_logging() {
    log_step "Deploying logging stack (ELK)..."

    kubectl create namespace "${LOGGING_NS}" --dry-run=client -o yaml | kubectl apply -f -

    helm repo add elastic https://helm.elastic.co 2>/dev/null || true
    helm repo update

    # Elasticsearch
    if ! kubectl get statefulset elasticsearch-master -n "${LOGGING_NS}" &>/dev/null; then
        helm upgrade --install elasticsearch elastic/elasticsearch \
            --namespace "${LOGGING_NS}" \
            --set replicas=3 \
            --set minimumMasterNodes=2 \
            --set resources.requests.cpu=2 \
            --set resources.requests.memory=4Gi \
            --set resources.limits.cpu=4 \
            --set resources.limits.memory=8Gi \
            --set volumeClaimTemplate.storageClassName=ssd \
            --set volumeClaimTemplate.resources.requests.storage=500Gi \
            --set esJavaOpts="-Xms2g -Xmx2g" \
            --wait --timeout 15m
    else
        log_warn "Elasticsearch already deployed, skipping"
    fi

    # Filebeat
    if ! kubectl get daemonset filebeat -n "${LOGGING_NS}" &>/dev/null; then
        kubectl create configmap filebeat-config \
            --from-file=filebeat.yml="${ROOT_DIR}/logging/filebeat/filebeat.yml" \
            -n "${LOGGING_NS}" --dry-run=client -o yaml | kubectl apply -f -

        kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: filebeat
  namespace: ${LOGGING_NS}
spec:
  selector:
    matchLabels:
      app: filebeat
  template:
    metadata:
      labels:
        app: filebeat
    spec:
      serviceAccountName: filebeat-sa
      terminationGracePeriodSeconds: 30
      containers:
        - name: filebeat
          image: docker.elastic.co/beats/filebeat:8.11.3
          args: ["-c", "/etc/filebeat.yml", "-e", "-strict.perms=false"]
          resources:
            requests:
              cpu: 100m
              memory: 256Mi
            limits:
              cpu: 500m
              memory: 512Mi
          volumeMounts:
            - name: config
              mountPath: /etc/filebeat.yml
              subPath: filebeat.yml
            - name: varlog
              mountPath: /var/log/containers
              readOnly: true
            - name: varlibdockercontainers
              mountPath: /var/lib/docker/containers
              readOnly: true
            - name: pod-logs
              mountPath: /app/logs
              readOnly: true
      volumes:
        - name: config
          configMap:
            name: filebeat-config
        - name: varlog
          hostPath:
            path: /var/log/containers
        - name: varlibdockercontainers
          hostPath:
            path: /var/lib/docker/containers
        - name: pod-logs
          hostPath:
            path: /data/logs/seckill
EOF
    else
        log_warn "Filebeat already deployed, skipping"
    fi

    # Logstash
    if ! kubectl get deployment logstash -n "${LOGGING_NS}" &>/dev/null; then
        kubectl create configmap logstash-config \
            --from-file=logstash.conf="${ROOT_DIR}/logging/logstash/logstash.conf" \
            --from-file=logstash.yml="${ROOT_DIR}/logging/logstash/logstash.yml" \
            -n "${LOGGING_NS}" --dry-run=client -o yaml | kubectl apply -f -

        kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: logstash
  namespace: ${LOGGING_NS}
spec:
  replicas: 2
  selector:
    matchLabels:
      app: logstash
  template:
    metadata:
      labels:
        app: logstash
    spec:
      containers:
        - name: logstash
          image: docker.elastic.co/logstash/logstash:8.11.3
          ports:
            - containerPort: 5044
              name: beats
            - containerPort: 9600
              name: api
          env:
            - name: LS_JAVA_OPTS
              value: "-Xms1g -Xmx1g"
          resources:
            requests:
              cpu: 1
              memory: 2Gi
            limits:
              cpu: 2
              memory: 3Gi
          volumeMounts:
            - name: config
              mountPath: /usr/share/logstash/pipeline/logstash.conf
              subPath: logstash.conf
            - name: config
              mountPath: /usr/share/logstash/config/logstash.yml
              subPath: logstash.yml
      volumes:
        - name: config
          configMap:
            name: logstash-config
---
apiVersion: v1
kind: Service
metadata:
  name: logstash-service
  namespace: ${LOGGING_NS}
spec:
  selector:
    app: logstash
  ports:
    - name: beats
      port: 5044
    - name: api
      port: 9600
EOF
    else
        log_warn "Logstash already deployed, skipping"
    fi

    log_info "Logging stack deployed"
}

verify_deployment() {
    log_step "Verifying deployment..."

    echo "=== Namespace Resources ==="
    kubectl get all -n "${NAMESPACE}"

    echo ""
    echo "=== Pod Status ==="
    kubectl get pods -n "${NAMESPACE}" -o wide

    echo ""
    echo "=== HPA Status ==="
    kubectl get hpa -n "${NAMESPACE}"

    echo ""
    echo "=== Ingress Status ==="
    kubectl get ingress -n "${NAMESPACE}"

    echo ""
    echo "=== Top Nodes ==="
    kubectl top nodes 2>/dev/null || echo "(metrics-server not available)"

    echo ""
    echo "=== Top Pods ==="
    kubectl top pods -n "${NAMESPACE}" 2>/dev/null || echo "(metrics-server not available)"

    log_info "Deployment verification complete"
}

usage() {
    cat <<EOF
Usage: $0 [OPTIONS] [COMPONENT]

OPTIONS:
    -e, --env ENV         Deployment environment (dev|staging|prod) [default: prod]
    -h, --help            Show this help message

COMPONENTS:
    all                   Deploy everything (default)
    infra                 Deploy only infrastructure (MySQL, Redis, Kafka)
    app                   Deploy only application
    monitoring            Deploy only monitoring stack
    logging               Deploy only logging stack

Examples:
    $0                          # Deploy everything to prod
    $0 -e dev app               # Deploy only the app to dev
    $0 infra                    # Deploy only infrastructure
EOF
}

main() {
    local component="${1:-all}"
    shift || true

    while [[ $# -gt 0 ]]; do
        case "$1" in
            -e|--env)
                DEPLOY_ENV="$2"
                export DEPLOY_ENV
                shift 2
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                component="$1"
                shift
                ;;
        esac
    done

    log_info "=============================================="
    log_info "  Seckill System K8s Deploy"
    log_info "  Environment: ${DEPLOY_ENV}"
    log_info "  Component:   ${component}"
    log_info "=============================================="

    case "${component}" in
        all)
            check_prerequisites
            deploy_namespaces
            deploy_secrets
            deploy_configmaps
            deploy_infrastructure
            deploy_app
            deploy_monitoring
            deploy_logging
            verify_deployment
            ;;
        infra)
            check_prerequisites
            deploy_namespaces
            deploy_secrets
            deploy_infrastructure
            ;;
        app)
            check_prerequisites
            deploy_namespaces
            deploy_secrets
            deploy_configmaps
            deploy_app
            ;;
        monitoring)
            check_prerequisites
            deploy_monitoring
            ;;
        logging)
            check_prerequisites
            deploy_logging
            ;;
        *)
            log_error "Unknown component: ${component}"
            usage
            exit 1
            ;;
    esac

    log_info "=============================================="
    log_info "  Deploy Complete!"
    log_info "  API:      https://api.seckill.com"
    log_info "  Grafana:  https://monitor.seckill.com/grafana"
    log_info "  Swagger:  https://admin.seckill.com/swagger-ui.html"
    log_info "=============================================="
}

main "$@"