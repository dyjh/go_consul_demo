package main

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
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

func registerServiceWithConsul(consulAddr string) error {
	config := consul.DefaultConfig()
	config.Address = consulAddr // 设置 Consul 服务的地址
	client, err := consul.NewClient(config)
	if err != nil {
		return err
	}

	registration := new(consul.AgentServiceRegistration)
	registration.ID = "hello-service"
	registration.Name = "hello-service"
	registration.Port = 1234

	// 获取服务器IP地址
	ip, err := getServerIP()
	if err != nil {
		return err
	}

	registration.Address = ip

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	ip, errIp := getServerIP()
	if errIp != nil {
		return
	}
	consulAddr := fmt.Sprintf("%s:8500", ip) // 替换为实际的 Consul 地址
	err := registerServiceWithConsul(consulAddr)
	if err != nil {
		fmt.Println("Error registering service with Consul:", err)
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
