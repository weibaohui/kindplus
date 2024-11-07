package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/weibaohui/kindplus/pkg/installer"
	"github.com/weibaohui/kom/kom_starter"
	"k8s.io/klog/v2"
)

func main() {
	kom_starter.Init()

	i := installer.Installer{
		Config: &installer.Config{
			Name:       "kind-1",
			BaseDomain: "dev.power.sd.istio.space",
			Port:       6552,
			Namespace:  "default",
		},
		Runtime: installer.NewRuntime(),
	}
	i.Deploy()

	// todo 自动发现集群中的kind小集群，并且自动生成kubeconfig
	// 是否可以采用一个公共pvc，各个小集群将kubeconfig写入到文件夹中，然后读取，并且注册到kom方便控制？
	// 如果kom可以控制，那么可以从程序中安装集群的初始化应用

	r := gin.Default()

	r.Use(cors.Default())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	klog.Infof("listen and serve on 0.0.0.0:%d", "80")
	err := r.Run(fmt.Sprintf(":%d", 80))
	if err != nil {
		klog.Fatalf("Error %v", err)
	}
}
