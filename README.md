* Take a look at [gatekeeper demo](gatekeeper/) for info on how to install OPA Gatekeeper and run basic demo

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

If you have go 1.16+ and docker, podman or nerdctl installed go install sigs.k8s.io/kind@v0.22.0 && kind create cluster is all you need!


kind create cluster


bring up guac via helm

kubectl port-forward svc/graphql-server 8080:8080
go run ./cmd/guacone collect files ../guac-data/docs/cdx_vuln.json

go run ./cmd/guacone certifier osv




docker build . -t ghcr.io/dejanb/guac-provider:latest

kind load docker-image ghcr.io/dejanb/guac-provider:latest --name kind

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
kubectl apply -f policy/examples/valid.yaml
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
kubectl delete -f policy/examples/valid.yaml
```