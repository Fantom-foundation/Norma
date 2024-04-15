// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

package network

import (
	"fmt"
	"net"
	"testing"
	"time"
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
		t.Cleanup(func() {
			_ = listener.Close()
		})
	}
}

func TestRegisterDuplicatedPortsForServices(t *testing.T) {
	service := ServiceDescription{
		Name:     "OperaPprof",
		Port:     6060,
		Protocol: "http",
	}

	serviceGroup := ServiceGroup{}

	if err := serviceGroup.RegisterService(&service); err != nil {
		t.Errorf("first registration must succeed")
	}

	if err := serviceGroup.RegisterService(&service); err == nil {
		t.Errorf("first registration must fail")
	}

}

func TestRetry(t *testing.T) {
	t.Parallel()

	var count int
	err := Retry(5, 1*time.Millisecond, func() error {
		count++
		if count >= 5 {
			return nil
		} else {
			return fmt.Errorf("no time to end yet")
		}
	})

	if err != nil {
		t.Errorf("Retry should success eventually")
	}

	if got, want := count, 5; got < want {
		t.Errorf("Retry finished early: %d < %d", got, want)
	}
}
