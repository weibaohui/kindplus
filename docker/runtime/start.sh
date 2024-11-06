# !/bin/bash
#将entrypoint.sh在后台运行
nohup dockerd-entrypoint.sh > /dev/null 2>&1 &
# 输出Echo 命令

# 检查 Docker 是否已启动的函数
check_docker() {
  until docker info >/dev/null 2>&1; do
    echo "等待 Docker 启动..."
    sleep 2
  done
  echo "Docker 已启动。"
}
# 调用函数检查 Docker 是否启动
check_docker

# 调用安装kind脚本
# 取环境变量，启动容器时注入
KIND_CLUSTER_NAME=${KIND_CLUSTER_NAME:-'kind-noname'}
KIND_CLUSTER_IP=${KIND_CLUSTER_IP:-'0.0.0.0'}
KIND_CLUSTER_PORT=${KIND_CLUSTER_PORT:-'6443'}
DOMAIN=${DOMAIN:-'kind-noname.k8m.site'}


bash kind-maker.sh  -d ${DOMAIN} -n ${KIND_CLUSTER_NAME} -p ${KIND_CLUSTER_PORT} -i ${KIND_CLUSTER_IP}

tail -f /dev/null