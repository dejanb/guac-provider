#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
cd "${REPO_ROOT}" || exit 1
GATEKEEPER_NAMESPACE=${NAMESPACE:-gatekeeper-system}

generate() {
    # generate CA key and certificate
    echo "Generating CA key and certificate for guac provider..."
    openssl genrsa -out ca.key 2048
    openssl req -new -x509 -days 365 -key ca.key -subj "/O=Gatekeeper/CN=Gatekeeper Root CA" -out ca.crt

    # generate server key and certificate
    echo "Generating server key and certificate for guac provider..."
    openssl genrsa -out server.key 2048
    openssl req -newkey rsa:2048 -nodes -keyout server.key -subj "/CN=guac-provider.${GATEKEEPER_NAMESPACE}" -out server.csr
    openssl x509 -req -extfile <(printf "subjectAltName=DNS:guac-provider.${GATEKEEPER_NAMESPACE}") -days 365 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
}

mkdir -p "${REPO_ROOT}/certs"
pushd "${REPO_ROOT}/certs"
generate
popd