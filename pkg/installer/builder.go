package installer

type Builder struct {
	Domain string // Kind 外部访问Kind集群使用的域名
	Port   int    // Kind APIServer 端口
	Name   string // Kind 集群名称
}

func NewBuilder() *Builder {
	return &Builder{}
}
func (b *Builder) SetDomain(domain string) *Builder {
	b.Domain = domain
	return b
}
func (b *Builder) SetPort(port int) *Builder {
	b.Port = port
	return b
}
func (b *Builder) SetName(name string) *Builder {
	b.Name = name
	return b
}
func (b *Builder) Build() error {

	return nil
}
