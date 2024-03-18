* Prepare

```
# guac
helm install guac ./charts/guac
kubectl port-forward svc/graphql-server 8080:8080

# gatekeeper
helm install gatekeeper/gatekeeper  \
    --name-template=gatekeeper \
    --namespace gatekeeper-system --create-namespace \
    --set enableExternalData=true \
    --set controllerManager.dnsPolicy=ClusterFirst,audit.dnsPolicy=ClusterFirst

# provider
minikube image load --overwrite ghcr.io/dejanb/guac-provider:latest

kubectl apply -f manifest/deployment.yaml -n gatekeeper-system
kubectl apply -f manifest/provider.yaml -n gatekeeper-system
kubectl apply -f manifest/service.yaml -n gatekeeper-system

```

* Demo

```
# guac
go run ./cmd/guacone collect files ../../cdx_vuln.json
go run ./cmd/guacone collect files ~/go/src/github.com/guacsec/guac-data/docs/cyclonedx/syft-cyclonedx-docker.io-library-bash.latest.json
go run ./cmd/guacone collect files ~/go/src/github.com/guacsec/guac-data/docs/cyclonedx/syft-cyclonedx-docker.io-library-alpine.latest.json
go run ./cmd/guacone collect files ../../guac-slsa-v0.5.json
go run ./cmd/guacone certifier osv
go run ./cmd/guacone certify package "critical vulnerability reported by maintainer" "pkg:alpine/alpine-baselayout@3.2.0-r18?arch=x86_64&upstream=alpine-baselayout&distro=alpine-3.15.6"

# provider
kubectl apply -f policy/template.yaml
kubectl apply -f policy/constraint.yaml

kubectl get constraint -o yaml

# deployments
kubectl apply -f policy/examples/vulnerable.yaml
kubectl apply -f policy/examples/bad.yaml
kubectl apply -f policy/examples/sbom.yaml
kubectl apply -f policy/examples/slsa.yaml
```