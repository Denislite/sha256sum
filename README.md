# sha256sum for k8s container

Implementation on Golang.<br> 
The sha256sum command computes and checks a SHA256 encrypted message digest.

---
## :hammer: Installation:
#### 1. Database:
Apply all annotations in directory "manifests/database/..":
```
kubectl apply -f manifests/db/postgres-secret.yaml
kubectl apply -f manifests/db/postgres-db-deployment.yaml
kubectl apply -f manifests/db/postgres-db-service.yaml
```
#### 2. K8S checksum:
Generate ca in /tmp :
```
cfssl gencert -initca ./webhook/tls/ca-csr.json | cfssljson -bare /tmp/ca
```

Generate private key and certificate for SSL connection:
```
cfssl gencert \
-ca=/tmp/ca.pem \
-ca-key=/tmp/ca-key.pem \
-config=./webhook/tls/ca-config.json \
-hostname="tcpdump-webhook,tcpdump-webhook.default.svc.cluster.local,tcpdump-webhook.default.svc,localhost,127.0.0.1" \
-profile=default \
./webhook/tls/ca-csr.json | cfssljson -bare /tmp/tcpdump-webhook
```

Move your SSL key and certificate to the ssl directory:
```
mv /tmp/tcpdump-webhook.pem ./pkg/webhook/ssl/tcpdump.pem
mv /tmp/tcpdump-webhook-key.pem ./pkg/webhook/ssl/tcpdump.key
```

Update ConfigMap data in the manifests/webhook/webhook-deployment.yaml file with your key and certificate:
```
cat ./pkg/webhook/ssl/tcpdump.key | base64 | tr -d '\n'
cat ./pkg/webhook/ssl/tcpdump.pem | base64 | tr -d '\n'
```

Update caBundle value in the manifests/webhook/webhook-configuration.yaml file with your base64 encoded CA certificate:
```
cat /tmp/ca.pem | base64 | tr -d '\n'
```
Build docker images webhook and hasher:
```
eval $(minikube docker-env)
docker build -t webhook -f pkg/webhook/Dockerfile .
docker build -t hasher .
```
Apply webhook annotation:
```
kubectl apply -f manifests/webhook/webhook-deployment.yaml
kubectl apply -f manifests/webhook/webhook-configuration.yaml
```
You can test working by applying test annotation:
```
kubectl apply -f manifests/hasher/test-deployment.yaml
```
