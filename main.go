package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct{}

func (h *HelloService) SayHello(request string, reply *string) error {
	// 获取服务器IP地址
	ip, err := getServerIP()
	if err != nil {
		return err
	}

	*reply = fmt.Sprintf("Hello world from %s!", ip)
	return nil
}

func getServerIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("Unable to determine server IP address")
}

func registerWithConsul() error {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	agent := client.Agent()

	ip, err := getServerIP()
	if err != nil {
		return err
	}

	port := 1234 // 修改为你的实际端口

	service := &api.AgentServiceRegistration{
		ID:      "jsonrpc-service",
		Name:    "jsonrpc",
		Address: ip,
		Port:    port,
	}

	err = agent.ServiceRegister(service)
	if err != nil {
		return err
	}

	fmt.Printf("Service registered with Consul: %s:%d\n", ip, port)
	return nil
}

func main() {
	// 注册服务到Consul
	err := registerWithConsul()
	if err != nil {
		fmt.Println("Error registering with Consul:", err)
		return
	}

	helloService := new(HelloService)
	rpc.Register(helloService)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server registered with Consul and listening on :1234...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
