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

tail -f /dev/null