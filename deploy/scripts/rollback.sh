#!/bin/bash
set -euo pipefail

#=======================================
# Seckill System Rollback Script
#=======================================

NAMESPACE="${NAMESPACE:-seckill}"
DEPLOYMENT="${DEPLOYMENT:-seckill-app}"
REVISION="${1:-0}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC}  $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

rollback_deployment() {
    log_info "=== Rollback History ==="
    kubectl rollout history deployment/"${DEPLOYMENT}" -n "${NAMESPACE}"

    if [ "${REVISION}" -eq 0 ]; then
        log_info "Rolling back to previous revision..."
        kubectl rollout undo deployment/"${DEPLOYMENT}" -n "${NAMESPACE}"
    else
        log_info "Rolling back to revision ${REVISION}..."
        kubectl rollout undo deployment/"${DEPLOYMENT}" -n "${NAMESPACE}" --to-revision="${REVISION}"
    fi

    log_info "Waiting for rollback to complete..."
    kubectl rollout status deployment/"${DEPLOYMENT}" -n "${NAMESPACE}" --timeout=5m

    log_info "=== Pod Status After Rollback ==="
    kubectl get pods -n "${NAMESPACE}" -l app="${DEPLOYMENT}" -o wide

    log_info "=== Current Revision ==="
    kubectl rollout history deployment/"${DEPLOYMENT}" -n "${NAMESPACE}" --revision=0 | tail -1

    log_info "Rollback completed successfully"
}

scale_down() {
    local replicas="${1:-0}"
    log_warn "Scaling ${DEPLOYMENT} down to ${replicas} replicas..."
    kubectl scale deployment/"${DEPLOYMENT}" -n "${NAMESPACE}" --replicas="${replicas}"
}

scale_up() {
    local replicas="${1:-3}"
    log_info "Scaling ${DEPLOYMENT} up to ${replicas} replicas..."
    kubectl scale deployment/"${DEPLOYMENT}" -n "${NAMESPACE}" --replicas="${replicas}"

    log_info "Waiting for pods to be ready..."
    kubectl rollout status deployment/"${DEPLOYMENT}" -n "${NAMESPACE}" --timeout=5m
}

restart_pods() {
    log_warn "Restarting all pods for ${DEPLOYMENT}..."
    kubectl rollout restart deployment/"${DEPLOYMENT}" -n "${NAMESPACE}"

    log_info "Waiting for rollout to complete..."
    kubectl rollout status deployment/"${DEPLOYMENT}" -n "${NAMESPACE}" --timeout=5m

    log_info "Pods restarted successfully"
}

usage() {
    cat <<EOF
Usage: $0 <COMMAND> [OPTIONS]

COMMANDS:
    rollback [REVISION]     Rollback to previous revision (or specified revision)
    history                 Show rollout history
    scale-down [REPLICAS]   Scale down deployment (default: 0)
    scale-up [REPLICAS]     Scale up deployment (default: 3)
    restart                 Rolling restart all pods
    status                  Show current deployment status

ENVIRONMENT VARIABLES:
    NAMESPACE               Kubernetes namespace (default: seckill)
    DEPLOYMENT              Deployment name (default: seckill-app)

Examples:
    $0 rollback                # Rollback to previous revision
    $0 rollback 5              # Rollback to revision 5
    $0 history                 # View rollout history
    $0 scale-down              # Scale to 0 replicas
    $0 scale-up 5              # Scale to 5 replicas
    $0 status                  # Show current status
EOF
}

main() {
    local command="${1:-}"

    case "${command}" in
        rollback)
            shift
            REVISION="${1:-0}"
            rollback_deployment
            ;;
        history)
            kubectl rollout history deployment/"${DEPLOYMENT}" -n "${NAMESPACE}"
            ;;
        scale-down)
            shift
            scale_down "${1:-0}"
            ;;
        scale-up)
            shift
            scale_up "${1:-3}"
            ;;
        restart)
            restart_pods
            ;;
        status)
            kubectl describe deployment/"${DEPLOYMENT}" -n "${NAMESPACE}"
            echo ""
            kubectl get pods -n "${NAMESPACE}" -l app="${DEPLOYMENT}" -o wide
            echo ""
            kubectl get hpa -n "${NAMESPACE}"
            ;;
        -h|--help|"")
            usage
            exit 0
            ;;
        *)
            log_error "Unknown command: ${command}"
            usage
            exit 1
            ;;
    esac
}

main "$@"