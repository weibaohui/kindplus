# kindplus
## 基础环境
### 镜像
* 运行：本方案使用Kind项目，因此采用Docker:dind为基础镜像
* 管理：使用nginx镜像，集成控制功能，建议使用TCP四层暴露，默认监听9443端口，根据多集群访问域名自动转发到对应的Kind中

0.运行在宿主集群，专门负责管理kind小集群
1.含有nginx，监听9443端口，将本集群内所有kind集群的证书收集，并加载到nginx.conf中
2.对外暴露9443端口，用于转发kind集群访问流量，需要使用4层方式对外开放。如NodePort、ingress-controller的四层转发
3.本程序不负责转发kind集群内的服务访问。
4.宿主集群上安装ingress-controller，专门负责监听32480的流量，使用ingress，转发到对应的kind集群的SVC上


## 执行
* /kind-install为集群安装目录，安装过程中产生的脚本或证书文件都保存在此目录下
* 检测/kind-init目录下是否有挂载的yaml文件，有即执行kubectl apply -f /kind-init 