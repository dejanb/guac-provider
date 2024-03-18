# guac Data Provider

A repository for using guac as a data provider for Gatekeeper.

## Prerequisites

- [`docker`](https://docs.docker.com/get-docker/)
- [`kind`](https://kind.sigs.k8s.io/)
- [`helm`](https://helm.sh/)
- [`kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl)

## Quick Start

1. Create a [kind cluster](https://kind.sigs.k8s.io/docs/user/quick-start/).

1.  Setup GUAC

```bash
git clone git@github.com:pxp928/kusari-helm-charts.git

cd kusari-helm-charts

git checkout update-guacrest

# TODO
#kubectl create namespace guac
#kubectl ns guac

helm dependency update ./charts/guac

helm install guac ./charts/guac

kubectl port-forward svc/graphql-server 8080:8080

git clone git@github.com:pxp928/artifact-ff.git

cd artifact-ff

git checkout issue-1734

go run ./cmd/guacone collect files ../../cdx_vuln.json
go run ./cmd/guacone collect files ~/go/src/github.com/guacsec/guac-data/docs/cyclonedx/syft-cyclonedx-docker.io-library-bash.latest.json
go run ./cmd/guacone collect files ~/go/src/github.com/guacsec/guac-data/docs/cyclonedx/syft-cyclonedx-docker.io-library-alpine.latest.json
go run ./cmd/guacone collect files ../../guac-slsa-v0.5.json

# optional
go run ./cmd/guacone collect files ~/go/src/github.com/guacsec/guac-data

go run ./cmd/guacone certifier osv

go run ./cmd/guacone collect files ../guac-data/cdx_guac.json

go run ./cmd/guacone certify package "critical vulnerability reported by maintainer" "pkg:alpine/alpine-baselayout@3.2.0-r18?arch=x86_64&upstream=alpine-baselayout&distro=alpine-3.15.6"
```

1. Install the latest version of Gatekeeper and enable the external data feature.

```bash
# Install the latest version of Gatekeeper with the external data feature enabled.
helm repo add gatekeeper https://open-policy-agent.github.io/gatekeeper/charts
helm install gatekeeper/gatekeeper  \
    --name-template=gatekeeper \
    --namespace gatekeeper-system --create-namespace \
    --set enableExternalData=true \
    --set controllerManager.dnsPolicy=ClusterFirst,audit.dnsPolicy=ClusterFirst
```

1. Build and deploy the guac data provider.

```bash
git clone https://github.com:dejanb/guac-provider.git
cd guac-provider

# generate a self-signed certificate for the guac data provider
./scripts/generate-tls-certificate.sh

# build the image via docker 
docker build . -t ghcr.io/dejanb/guac-provider:latest

# load the image into kind
kind load docker-image ghcr.io/dejanb/guac-provider:latest --name kind

# Install guac data provider into gatekeeper-system to use mTLS
helm install guac-provider charts/guac-provider \
    --set provider.tls.caBundle="$(cat certs/ca.crt | base64 | tr -d '\n\r')" \
    --namespace gatekeeper-system
```

1. Install constraint template and constraint.

```bash
kubectl apply -f policy/template.yaml
kubectl apply -f policy/constraint.yaml
```

1. Check the logs for the guac-provider
```bash
kubectl logs -n gatekeeper-system deployments/guac-provider -f
```

1. Examples
```
kubectl create ns test
kubectl apply -f policy/examples/vulnerable.yaml -n test
kubectl apply -f policy/examples/bad.yaml -n test
kubectl apply -f policy/examples/sbom.yaml -n test
kubectl apply -f policy/examples/slsa.yaml -n test
kubectl apply -f policy/examples/good.yaml -n test
```

1. Delete

```bash
kubectl delete -f policy/

helm uninstall guac
helm uninstall guac-provider --namespace gatekeeper-system
helm uninstall gatekeeper --namespace gatekeeper-system
```