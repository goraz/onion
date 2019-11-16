#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

ROOT=$(dirname "${BASH_SOURCE}")/../..

ETCD_VERSION=${ETCD_VERSION:-v3.3.17}

mkdir -p "${ROOT}/third_party"
cd "${ROOT}/third_party"
curl -sL https://github.com/etcd-io/etcd/releases/download/${ETCD_VERSION}/etcd-${ETCD_VERSION}-linux-amd64.tar.gz \
  | tar xzf -

exec etcd-${ETCD_VERSION}-linux-amd64/etcd
