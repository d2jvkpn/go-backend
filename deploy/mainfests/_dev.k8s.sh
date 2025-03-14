#!/bin/bash
set -eu -o pipefail; _wd=$(pwd); _path=$(dirname $0)


exit

#### copy go-backend to node k8s-cp01 and create /data/logs on all worker nodes
ansible k8s-cp01 -m copy -a 'src=./mainfests dest=./go-backend'

# ansible k8s_workers -m shell --become -a 'mkdir -p /data/local && chmod -R 777 /data/local'

ssh k8s-cp01

#### Configmap
# kubectl -n dev create configmap go-backend --from-file=dev.yaml
# kubectl create configmap go-backend --from-file=mainfests/dev.yaml

kubectl -n dev create configmap go-backend \
  --from-file=dev.yaml -o yaml --dry-run=client |
  kubectl apply -f -

kubectl get configmap go-backend -o yaml

##### Deployment, ClusterIP, Ingress(http) and HPA
kubectl apply -f dev.k8s-app.yaml
# kubectl get deploy/go-backend
# kubectl describe deploy/go-backend

kubectl apply -f dev.k8s-ctrl.yaml

kubectl -n dev get svc
# kubectl -n dev patch svc go-backend -p '{"spec":{"type":"LoadBalancer"}}'

####
kubectl get pods -o wide

kubectl get pods -l app=go-backend -o=custom-columns=NAME:.metadata.name |
  sed '1d' |
  xargs -i kubectl describe pod/{}

kubectl get pods -l app=go-backend |
  awk 'NR>1{print $1}' |
  xargs -i kubectl logs pod/{}

kubectl scale --replicas=1 deploy/go-backend

pod=$(kubectl get pods -l app=go-backend | awk 'NR==2{print $1; exit}')
kubectl get pod/$pod -o wide
kubectl exec -it $pod -- ls

####
curl -H 'Host: go-backend.k8s-dev.local' k8s.local/meta | jq
curl -H 'Host: go-backend.k8s-dev.local' k8s.local/api/v1/open/hello | jq

exit
# method 1
kubectl rollout restart deploy/go-backend

# method 2
kubectl scale --replicas=0 deploy/go-backend
kubectl scale --replicas=3 deploy/go-backend

# method 3
# imagePullPolicy: "Always"
# imagePullPolicy: "IfNotPresent"
# imagePullPolicy: "Nerver"
kubectl set image deploy/go-backend \
  go-backend=registry.cn-shanghai.aliyuncs.com/d2jvkpn/go-backend:dev@sha256:xxxxxx

# method 4
kubectl patch deploy/go-backend -p \
  '{"spec":{"template":{"spec":{"terminationGracePeriodSeconds":30}}}}'

date +"==> %FT%T%:z" && \
  kubectl get pods -o wide && \
  curl -H 'Host: go-backend.k8s-dev.local' k8s.local/meta | jq
