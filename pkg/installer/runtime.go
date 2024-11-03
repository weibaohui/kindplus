package installer

type Runtime struct {
	BaseImage string
}

func NewRuntime() *Runtime {
	return &Runtime{
		BaseImage: "weibh/kind-in-docker:docker-27-kind-0.24",
	}
}
