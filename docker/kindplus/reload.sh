#!/bin/bash

APP_PATH="/app/kindplus"           # 上传的新文件名始终为 kindplus
CURRENT_PATH="/app/kindplus.current" # 运行中的文件名为 kindplus.current
WATCH_PATH="/app"                  # 监控的目录

# 启动 kindplus 程序
start_app() {
    echo "Starting kindplus..."
    mv $APP_PATH $CURRENT_PATH          # 将 kindplus 重命名为 kindplus.current
    chmod +x $CURRENT_PATH              # 确保文件可执行
    nohup $CURRENT_PATH > /dev/stdout 2> /dev/stderr &
    APP_PID=$!
    echo "kindplus started with PID: $APP_PID"
}

# 停止 kindplus 程序
stop_app() {
    if [ -n "$APP_PID" ]; then
        echo "Stopping kindplus with PID: $APP_PID"
        kill "$APP_PID"
        wait "$APP_PID" 2>/dev/null
        echo "kindplus stopped"
    fi
}

# 初次启动程序
start_app

# 使用 inotifywait 监控 APP_PATH 的 close_write 事件，确保只在新的文件写入完成后触发
while true; do
    # 监听 close_write 事件，且排除文件重命名可能带来的误触发
    inotifywait -e close_write --exclude "kindplus.current" $WATCH_PATH

    # 检查是否是新的 kindplus 文件写入完成
    if [ -f "$APP_PATH" ]; then
        echo "Detected new kindplus version, restarting..."

        # 停止当前程序
        stop_app

        # 删除旧的 kindplus.current 文件并启动新的程序
        rm -f $CURRENT_PATH
        start_app
    fi
done
