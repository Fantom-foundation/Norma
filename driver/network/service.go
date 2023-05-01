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
	Name string
	Port Port
}

// Port provides an alias type for a TCP port.
type Port uint16

// GetFreePort obtains a free TCP port on the local system. Note, that after
// this call the port is not reserved. Thus, consecutive calls may produce the
// same free port until it is actually bound to some application.
func GetFreePort() (Port, error) {
	for i := 0; i < 10; i++ {
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
		return Port(res), nil
	}
	return 0, fmt.Errorf("failed to allocate a free port on the system")
}
