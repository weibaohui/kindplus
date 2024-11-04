#!/bin/bash

# 检查输入参数
if [ $# -ne 3 ]; then
  echo "Usage: $0 <custom-domain> <cluster-name> <cluster-port>"
  exit 1
fi

DOMAIN=$1
KIND_CLUSTER_NAME=$2
KIND_CLUSTER_PORT=$3
INNER_DOMAIN="kind-svc.ns.svc.cluster.local"
mkdir -p /kind-install && cd /kind-install
#先进行清理，避免出现混乱
kind delete cluster --name $KIND_CLUSTER_NAME

#创建证书目录
rm -rf certs &&mkdir -p "certs"
# 检查 /path/to/file 是否存在
if [ -f "/kind-install/kind-config.yaml" ]; then
  echo "文件 /kind-install/kind-config.yaml 存在"
else
  echo "文件 /kind-install/kind-config.yaml 不存在,创建默认配置文件"
  # 创建 Kind 配置文件，指定挂载证书路径
  cat <<EOF > kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: $KIND_CLUSTER_NAME
networking:
  # disableDefaultCNI: true
  apiServerAddress: 0.0.0.0
  # ipFamily: dual
  apiServerPort: $KIND_CLUSTER_PORT
nodes:
  - role: control-plane
    kubeadmConfigPatches:
          - |
            ---
            kind: ClusterConfiguration
            apiServer:
                certSANs:
                  - ${DOMAIN} #外部访问域名
                  - ${INNER_DOMAIN} #内部访问域名,从nginx实例访问kind的域名
                  - 0.0.0.0
                  - 192.168.182.124 #从nginx实例访问kind的IP
                  - localhost
                  - 127.0.0.1
                extraArgs:
                  enable-admission-plugins: MutatingAdmissionWebhook,ValidatingAdmissionWebhook
    extraMounts:
      - hostPath: ./certs                  # 主机上的证书目录路径
        containerPath: /etc/kubernetes/pki # 挂载到容器中的路径
EOF

fi

# 创建 Kind 集群
kind create cluster --config kind-config.yaml



######签发kubeconfig
# 固定路径
CA_CERT="certs/ca.crt"
CA_KEY="certs/ca.key"
API_SERVER="https://${DOMAIN}:${KIND_CLUSTER_PORT}"

# 生成私钥
openssl genrsa -out "client.key" 2048

# 生成证书签名请求 (CSR) 时添加 SAN
openssl req -new -key "client.key" -out "client.csr" -subj "/CN=client-user/O=my-org" -reqexts SAN -config <(cat /etc/ssl/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:${DOMAIN},DNS:192.168.182.124"))

# 使用 CA 签发证书
openssl x509 -req -in "client.csr" -CA "${CA_CERT}" -CAkey "${CA_KEY}" -CAcreateserial -out "client.crt" -days 10000 -extfile <(printf "subjectAltName=DNS:%s,DNS:192.168.182.124" "${DOMAIN}")



# 将 CSR 文件内容编码为 Base64
CSR_BASE64=$(cat client.csr | base64 | tr -d '\n')
CSR_NAME="client-cert-request-$RANDOM"

# 创建 CSR 资源
cat <<EOF | kubectl apply -f -
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: $CSR_NAME
spec:
  request: $CSR_BASE64
  signerName: kubernetes.io/kube-apiserver-client
  usages:
    - client auth
#    - server auth
EOF

# 批准 CSR
kubectl certificate approve "$CSR_NAME"
# 赋权
# 创建 CSR 资源
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: client-admin
rules:
  - apiGroups: ["*"]            # 允许所有 API 组
    resources: ["*"]            # 允许所有资源类型
    verbs: ["*"]                # 允许所有操作
EOF
kubectl create clusterrolebinding client-admin-binding --clusterrole=client-admin --user=client-user


# 等待证书生成
echo "Waiting for CSR to be signed..."
for i in {1..20}; do
  CERT=$(kubectl get csr "$CSR_NAME" -o jsonpath='{.status.certificate}' 2>/dev/null)
  if [ -n "$CERT" ]; then
    echo "Certificate signed successfully!"
    break
  fi
  sleep 1
done

if [ -z "$CERT" ]; then
  echo "Error: Certificate signing request not approved or failed to retrieve certificate."
  exit 1
fi

# 生成 kubeconfig 文件
KUBECONFIG_FILE="kubeconfig.yaml"
cat <<EOF > "${KUBECONFIG_FILE}"
apiVersion: v1
kind: Config
clusters:
- name: kubernetes
  cluster:
    server: $API_SERVER
    certificate-authority-data: $(base64 < certs/ca.crt | tr -d '\n')
contexts:
- name: default
  context:
    cluster: kubernetes
    user: client-user
current-context: default
users:
- name: client-user
  user:
    client-certificate-data: $(base64 < client.crt | tr -d '\n')
    client-key-data: $(base64 < client.key | tr -d '\n')
EOF

### TODO 将kubeconfig文件post到控制器
echo "Kind cluster created with custom domain: ${DOMAIN}"
echo "Server certificate generated with SAN: ${DOMAIN}"
echo "Kubeconfig file generated: ${KUBECONFIG_FILE}"

# 验证安装情况
echo "查看集群信息"
kubectl get nodes
echo "查看集群状态"
kubectl get pods -A
echo "查看端口占用情况"
netstat -anpt | grep $KIND_CLUSTER_PORT
echo "查看证书信息"
openssl x509 -in $KIND_CLUSTER_NAME/client.crt -noout -text | grep -A1 "Subject Alternative Name"




# 清理临时文件
echo "清理临时文件"
#rm client.key client.csr client.crt
kubectl delete csr "$CSR_NAME"


# 查看/kind-init文件夹下的目录，输出文件列表

# 检查 /kind-init 目录是否存在
if [ -d "/kind-init" ]; then
  echo "目录 /kind-init 存在，执行初始化"
  ls -l /kind-init/
  kubectl apply -f /kind-init
else
  echo "目录 /kind-init 不存在，跳过初始化"
fi

echo "$KIND_CLUSTER_NAME https://$DOMAIN:$KIND_CLUSTER_PORT 集群创建完成！"
echo "======"
cat $KUBECONFIG_FILE
