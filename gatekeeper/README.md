* Install Gatekeeper
```
kubectl create namespace kubecon-2024

kubectl ns kubecon-2024

kubectl apply -f https://raw.githubusercontent.com/open-policy-agent/gatekeeper/v3.15.0/deploy/gatekeeper.yaml
```

* Install resources
```
kubectl apply -f templates/k8srequiredlabels_template.yaml

kubectl get crd

kubectl apply -f constraints/pods-must-have-app.yaml

kubectl get constraints -o yaml
```

* Try to run containers
```
kubectl run nginx --image=nginx

kubectl run nginx --image=nginx --labels=app=kubecon-2024

kubectl get pods
```

* Cleanup 

```
kubectl delete pods nginx
kubectl delete -f constraints/pods-must-have-app.yaml
kubectl delete -f templates/k8srequiredlabels_template.yaml
```