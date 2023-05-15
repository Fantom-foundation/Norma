package network

import (
	"fmt"
	"net"
	"testing"
)

func TestGetFreePort(t *testing.T) {
	port, err := GetFreePort()
	if err != nil {
		t.Errorf("failed to obtain a free port: %v", err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Errorf("provided port %d is not free", port)
	}
	listener.Close()
}

func TestGetFreePorts(t *testing.T) {
	ports, err := GetFreePorts(5)
	if err != nil {
		t.Errorf("failed to obtain a free ports: %v", err)
	}
	if got, want := len(ports), 5; got != want {
		t.Errorf("invalid number of ports obtained, got %d, wanted %d", got, want)
	}
	for _, port := range ports {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			t.Errorf("provided port %d is not free", port)
		}
		defer listener.Close()
	}
}
