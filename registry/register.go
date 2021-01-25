package registry

import "google.golang.org/grpc/metadata"

//规定了注册服务的格式
type ServiceInfo struct {
	InstanceId string      //实例ID
	Name       string      //名称
	Version    string      //版本
	Address    string      //地址
	Metadata   metadata.MD //元数据
}

//注册器接口
type RegisterIns interface {
	Register(service *ServiceInfo) error
	UnRegister(service *ServiceInfo) error
	Close()
}
