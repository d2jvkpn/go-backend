#!/usr/bin/env bash
set -eu -o pipefail # -x
_wd=$(pwd); _path=$(dirname $0 | xargs -i readlink -f {})

exit
kubectl get ns

kubectl config get-contexts
kubectl config current-context

kubectl config use-context k8s-prod --namespace=prod
kubectl config set-context --current --namespace=prod

#### ingress tls(https)
kubectl create secret tls k8s.domain --key k8s.domain.key --cert k8s.domain.pem
# domains: *.k8s.domain

kubectl -n prod get secret/k8s.domain

#### create secret and ingress-https
kubectl create secret tls k8s.domain \
  --key k8s.domain.key --cert k8s.domain.pem \
  --dry-run=client -o yaml |
  kubectl apply -f -

kubectl create secret docker-registry k8s.domain \
  --docker-server=registry.k8s.domain \
  --docker-email=EMAIL \
  --docker-username=USERNAME \
  --docker-password=PASSWORD

kubectl apply -f go-backend/prod.k8s-apps.yaml
kubectl apply -f go-backend/prod.k8s-ctrl.yaml

#### get image sha256 of containers
kubectl get pods -l app=go-backend -o json |
  jq -r '.items[].status.containerStatuses[0].imageID'

kubectl get pods -l app=go-backend -o json |
  jq -r '.items[].status | .hostIP + ", " + .containerStatuses[0].imageID + ", " + .phase'
