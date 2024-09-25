#!/bin/sh

set -e

# Function to perform base64 encoding without line breaks using cat and pipe
encode_base64() {
  cat "$1" | base64 | tr -d '\n'
}


# Root CA
openssl genrsa -out root-ca-key.pem 2048
openssl req -new -x509 -sha256 -key root-ca-key.pem -subj "/C=MY/ST=SELANGOR/L=BANGI/O=SILENTMODE/OU=ENGINEERING/CN=root.dns.a-record" -out root-ca.pem -days 36525
# Admin cert
openssl genrsa -out admin-key-temp.pem 2048
openssl pkcs8 -inform PEM -outform PEM -in admin-key-temp.pem -topk8 -nocrypt -v1 PBE-SHA1-3DES -out admin-key.pem
openssl req -new -key admin-key.pem -subj "/C=MY/ST=SELANGOR/L=BANGI/O=SILENTMODE/OU=ENGINEERING/CN=A" -out admin.csr
openssl x509 -req -in admin.csr -CA root-ca.pem -CAkey root-ca-key.pem -CAcreateserial -sha256 -out admin.pem -days 36525
# Node cert 1
openssl genrsa -out node1-key-temp.pem 2048
openssl pkcs8 -inform PEM -outform PEM -in node1-key-temp.pem -topk8 -nocrypt -v1 PBE-SHA1-3DES -out node1-key.pem
openssl req -new -key node1-key.pem -subj "/C=MY/ST=SELANGOR/L=BANGI/O=SELANGOR/OU=ENGINEERING/CN=opensearch-cluster-master." -out node1.csr
echo 'subjectAltName=DNS:opensearch-cluster-master.' > node1.ext
openssl x509 -req -in node1.csr -CA root-ca.pem -CAkey root-ca-key.pem -CAcreateserial -sha256 -out node1.pem -days 36525 -extfile node1.ext
# Node cert 2
openssl genrsa -out node2-key-temp.pem 2048
openssl pkcs8 -inform PEM -outform PEM -in node2-key-temp.pem -topk8 -nocrypt -v1 PBE-SHA1-3DES -out node2-key.pem
openssl req -new -key node2-key.pem -subj "/C=MY/ST=SELANGOR/L=BANGI/O=SILENTMODE/OU=ENGINEERING/CN=node2.dns.a-record" -out node2.csr
echo 'subjectAltName=DNS:node2.dns.a-record' > node2.ext
openssl x509 -req -in node2.csr -CA root-ca.pem -CAkey root-ca-key.pem -CAcreateserial -sha256 -out node2.pem -days 36525 -extfile node2.ext
# Client cert
openssl genrsa -out client-key-temp.pem 2048
openssl pkcs8 -inform PEM -outform PEM -in client-key-temp.pem -topk8 -nocrypt -v1 PBE-SHA1-3DES -out client-key.pem
openssl req -new -key client-key.pem -subj "/C=MY/ST=SELANGOR/L=BANGI/O=SILENTMODE/OU=ENGINEERING/CN=opensearch-dashboards.default.svc.cluster.local" -out client.csr
echo 'subjectAltName=DNS:opensearch-dashboards.default.svc.cluster.local' > client.ext
openssl x509 -req -in client.csr -CA root-ca.pem -CAkey root-ca-key.pem -CAcreateserial -sha256 -out client.pem -days 36525 -extfile client.ext
# Cleanup
# rm admin-key-temp.pem
# rm admin.csr
# rm node1-key-temp.pem
# rm node1.csr
# rm node1.ext
# rm node2-key-temp.pem
# rm node2.csr
# rm node2.ext
# rm client-key-temp.pem
# rm client.csr
# rm client.ext

# Generate Kubernetes Secret manifest
SECRET_NAME=opensearch-certs
NAMESPACE=default  # Change to desired namespace

cat <<EOF > secret_certificate.yaml
apiVersion: v1
kind: Secret
metadata:
  name: ${SECRET_NAME}
  namespace: ${NAMESPACE}
type: Opaque
data:
  root-ca-key.pem: $(encode_base64 root-ca-key.pem)
  root-ca.pem: $(encode_base64 root-ca.pem)
  admin-key.pem: $(encode_base64 admin-key.pem)
  admin.pem: $(encode_base64 admin.pem)
  node1-key.pem: $(encode_base64 node1-key.pem)
  node1.pem: $(encode_base64 node1.pem)
  node2-key.pem: $(encode_base64 node2-key.pem)
  node2.pem: $(encode_base64 node2.pem)
  client-key.pem: $(encode_base64 client-key.pem)
  client.pem: $(encode_base64 client.pem)
EOF

echo "Kubernetes Secret manifest 'secret_manifest.yaml' has been created successfully."