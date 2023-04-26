package global

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// RegisterConsul 注册服务
func RegisterConsul() (serviceId string) {
	config := api.DefaultConfig()
	config.Address = ServerConfig.Consul.Host + ":" + ServerConfig.Consul.Port

	client, err := api.NewClient(config)
	if err != nil {
		Logger.Error("Consul连接失败", zap.String("error ", err.Error()))
		panic(err)
	}
	agent := client.Agent()

	ip := ServerConfig.Ip
	port := ServerConfig.Port

	serviceId = fmt.Sprintf("%s-%s", ip, uuid.NewV4())

	reg := new(api.AgentServiceRegistration)
	reg.ID = serviceId           // 服务节点的名称
	reg.Name = ServiceName       // 服务名称
	reg.Address = ip             // 服务 IP
	reg.Port = port              // 服务端口ßßßß
	reg.Tags = []string{"v1000"} // tag，可以为空

	checkUrl := fmt.Sprintf("%s:%d", ip, port)
	reg.Check = &api.AgentServiceCheck{ // 健康检查
		GRPC:                           checkUrl,
		GRPCUseTLS:                     false,
		Timeout:                        "3s",
		Interval:                       "5s",  // 健康检查间隔
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务，注销时间，相当于过期时间
	}

	if err = agent.ServiceRegister(reg); err != nil {
		Logger.Error("服务注册失败", zap.String("error ", err.Error()))
		panic(err)
	}

	Logger.Info("服务注册成功")
	return
}

// DeRegisterConsul 注销注册
func DeRegisterConsul(serviceId string) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", ServerConfig.Consul.Host, ServerConfig.Consul.Port)

	client, err := api.NewClient(config)
	if err != nil {
		return
	}
	if err := client.Agent().ServiceDeregister(serviceId); err != nil {
		Logger.Error("服务注销失败", zap.String("error ", err.Error()))
	}

	Logger.Info("服务注销成功")
}
