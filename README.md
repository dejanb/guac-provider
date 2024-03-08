* Take a look at [gatekeeper demo](gatekeeper/) for info on how to install OPA Gatekeeper and run basic demo

* Build and install Guac provider
```
docker build . -t ghcr.io/dejanb/guac-provider:latest

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
kubectl apply -f policy/examples/error.yaml
kubectl apply -f policy/examples/valid.yaml
```