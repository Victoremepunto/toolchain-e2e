#!/usr/bin/env bash
set -euo pipefail

# Perf test stages: baseline then 2000-user load test.
# SM3 and RHOAI must be installed before running this script.

export KUBECONFIG=setup/infra/clusters/vmugicag-perf/auth/kubeconfig

WORKLOADS="openshift-operators:servicemesh-operator3"
WORKLOADS+=",cert-manager-operator:cert-manager-operator-controller-manager"
WORKLOADS+=",cert-manager:cert-manager"
WORKLOADS+=",cert-manager:cert-manager-cainjector"
WORKLOADS+=",cert-manager:cert-manager-webhook"
WORKLOADS+=",openshift-jobset-system:jobset-operator"
WORKLOADS+=",openshift-nfd:nfd-controller-manager"
WORKLOADS+=",nvidia-gpu-operator:gpu-operator"
WORKLOADS+=",redhat-ods-operator:rhods-operator"
WORKLOADS+=",redhat-ods-applications:dashboard-redirect"
WORKLOADS+=",redhat-ods-applications:kserve-controller-manager"
WORKLOADS+=",redhat-ods-applications:llama-stack-k8s-operator-controller-manager"
WORKLOADS+=",redhat-ods-applications:llmisvc-controller-manager"
WORKLOADS+=",redhat-ods-applications:mlflow-operator-controller-manager"
WORKLOADS+=",redhat-ods-applications:model-serving-api"
WORKLOADS+=",redhat-ods-applications:notebook-controller-deployment"
WORKLOADS+=",redhat-ods-applications:odh-model-controller"
WORKLOADS+=",redhat-ods-applications:odh-notebook-controller-manager"
WORKLOADS+=",redhat-ods-applications:rhods-dashboard"

run_test() {
  local desc="$1"; shift
  echo "--- ${desc} ---"
  timeout --signal=KILL 3600 go run setup/main.go "$@" || {
    echo "WARNING: test '${desc}' exited non-zero or was killed by timeout"
  }
}

echo "===== Stage 1: Baseline (1 user, with SM3 + RHOAI installed) ====="

run_test "Stage 1: 1-user baseline" --users 1 --default 1 --custom 0 \
  --username baseline --testname=baseline \
  --workloads "${WORKLOADS}" --interactive=false

echo "===== Stage 2: 2000-user load test ====="

run_test "Stage 2: 2000-user" --users 2000 --default 1750 --custom 250 \
  --template setup/resources/rhoai-user-workloads.yaml \
  --username rhoai34 --testname=rhoai34 \
  --workloads "${WORKLOADS}" --interactive=false

echo ""
echo "===== ALL STAGES COMPLETE ====="
echo "Results in tmp/results/"
ls -lt tmp/results/*.csv | head -6
