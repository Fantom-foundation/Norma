package network

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// ServiceDescription is
type ServiceDescription struct {
	Name     string
	Port     Port
	Protocol string
}

type ServiceGroup map[Port]*ServiceDescription

// RegisterService installs a supported service in the registry.
func (s *ServiceGroup) RegisterService(service *ServiceDescription) error {
	if _, exists := (*s)[service.Port]; exists {
		return fmt.Errorf("port %d already assigned - it is not supported to bind the same port many times", service.Port)
	}

	(*s)[service.Port] = service
	return nil
}

// Services returns all registered services from the registry.
func (s *ServiceGroup) Services() []*ServiceDescription {
	res := make([]*ServiceDescription, 0, len(*s))
	for _, v := range *s {
		res = append(res, v)
	}

	return res
}

// Port provides an alias type for a TCP port.
type Port uint16

// GetFreePort obtains a free TCP port on the local system. Note, that after
// this call the port is not reserved. Thus, consecutive calls may produce the
// same free port until it is actually bound to some application.
func GetFreePort() (Port, error) {
	ports, err := GetFreePorts(1)
	if err != nil {
		return 0, err
	}
	return ports[0], nil
}

// GetFreePorts obtains a list of free TCP ports on the local system.  Note
// that after this call the ports are not reserved. Thus, consecutive calls may
// produce the same free ports until it is actually bound to some application.
func GetFreePorts(num int) ([]Port, error) {
	ports := make([]Port, 0, num)
	for len(ports) < num {
		found := false
		for i := 0; !found && i < 10; i++ {
			listener, err := net.Listen("tcp", "")
			if err != nil {
				log.Printf("failed to create a new listening port")
				continue
			}
			defer listener.Close()

			port := listener.Addr().String()
			columnPos := strings.LastIndex(port, ":")
			if columnPos < 0 {
				log.Printf("invalid port format: %s", port)
				continue
			}
			port = port[columnPos+1:]

			res, err := strconv.ParseUint(port, 10, 16)
			if err != nil {
				log.Printf("invalid port format: %s, err: %v", port, err)
				continue
			}
			ports = append(ports, Port(res))
			found = true
		}
		if !found {
			return nil, fmt.Errorf("failed to allocate a free port on the system")
		}
	}
	return ports, nil
}
