* Create Kind Cluster

If you have go 1.16+ and docker, podman or nerdctl installed `go install sigs.k8s.io/kind@v0.22.0` && `kind create cluster` is all you need!


* Setup GUAC

cdx_vuln.json is attached to this repo

```

git clone git@github.com:pxp928/kusari-helm-charts.git

cd kusari-helm-charts

git checkout update-guacrest

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

go run ./cmd/guacone certify package "critical vulnerability reported by maintainer" "pkg:alpine/alpine-baselayout@3.2.0-r18?arch=x86_64&upstream=alpine-baselayout&distro=alpine-3.15.6"

```



* Install Gatekeeper
```
helm repo add gatekeeper https://open-policy-agent.github.io/gatekeeper/charts
helm install gatekeeper/gatekeeper  \
    --name-template=gatekeeper \
    --namespace gatekeeper-system --create-namespace \
    --set enableExternalData=true \
    --set controllerManager.dnsPolicy=ClusterFirst,audit.dnsPolicy=ClusterFirst
```


* Build and install Guac provider
```
docker build . -t ghcr.io/dejanb/guac-provider:latest

# kind
kind load docker-image ghcr.io/dejanb/guac-provider:latest --name kind
# minikube
minikube image load --overwrite ghcr.io/dejanb/guac-provider:latest

kubectl apply -f manifest/deployment.yaml -n gatekeeper-system
kubectl apply -f manifest/provider.yaml -n gatekeeper-system
kubectl apply -f manifest/service.yaml -n gatekeeper-system
```

* Apply resources
```
kubectl apply -f policy/template.yaml
kubectl apply -f policy/constraint.yaml
```

* Try to run deployments
```
kubectl apply -f policy/examples/vulnerable.yaml
kubectl apply -f policy/examples/bad.yaml
kubectl apply -f policy/examples/sbom.yaml
kubectl apply -f policy/examples/slsa.yaml
```

* Receive an validation failure
```

Error from server (Forbidden): error when creating "policy/examples/valid.yaml": admission webhook "validation.gatekeeper.sh" denied the request: [guac] Image ghcr.io/guacsec/vul-image:latest@sha256:b6f1a6e034d40c240f1d8b0a3f5481aa0a315009f5ac72f736502939419c1855 contains 9 vulnerabilities
```

* Delete
```
kubectl delete -f manifest/deployment.yaml -n gatekeeper-system
kubectl delete -f manifest/provider.yaml -n gatekeeper-system
kubectl delete -f manifest/service.yaml -n gatekeeper-system
```

* Delete resources
```
kubectl delete -f policy/template.yaml
kubectl delete -f policy/constraint.yaml
```


* Delete deployments
```
kubectl delete -f policy/examples/vulnerable.yaml
kubectl delete -f policy/examples/sbom.yaml
kubectl delete -f policy/examples/bad.yaml
kubectl apply -f policy/examples/slsa.yaml
```