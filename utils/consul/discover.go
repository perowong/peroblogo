package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

func ConsulHealthService(addr, serviceName string) (services []*consulapi.ServiceEntry, err error) {
	config := consulapi.DefaultConfig()
	config.Address = addr
	client, err := consulapi.NewClient(config)
	if err != nil {
		return
	}

	services, _, err = client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return
	}
	return
}

func ConsulDiscover(addr, serviceName string, ipOnly bool) (string, error) {
	services, err := ConsulHealthService(addr, serviceName)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no passing service found")
	}

	format := "%s:%d"
	if !ipOnly {
		format = "http://%s:%d"
	}

	// Future: load balance
	return fmt.Sprintf(format, services[0].Service.Address, services[0].Service.Port), nil
}
