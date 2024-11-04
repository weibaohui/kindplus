# mkind
## 基础环境
### 镜像
* 运行：本方案使用Kind项目，因此采用Docker:dind为基础镜像
* 管理：使用nginx镜像，集成控制功能，建议使用TCP四层暴露，默认监听9443端口，根据多集群访问域名自动转发到对应的Kind中


## 执行
* /kind-install为集群安装目录，安装过程中产生的脚本或证书文件都保存在此目录下
* 检测/kind-init目录下是否有挂载的yaml文件，有即执行kubectl apply -f /kind-init 