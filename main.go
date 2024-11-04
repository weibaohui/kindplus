package main

import (
	"github.com/weibaohui/kindplus/pkg/installer"
	"github.com/weibaohui/kom/kom_starter"
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
}
