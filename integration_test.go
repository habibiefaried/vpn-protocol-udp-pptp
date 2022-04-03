package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
	"time"
)

func InfraDown(t *testing.T) {
	cmd := exec.Command("docker-compose", "down")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(stdoutStderr))
		t.Fatal("Running command error: ", err)
	}
}

func InfraBuild(t *testing.T) {
	cmd := exec.Command("docker-compose", "build")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(stdoutStderr))
		t.Fatal("Running command error: ", err)
	}
}

func InfraUp(t *testing.T) {
	cmd := exec.Command("docker-compose", "up", "-d")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(stdoutStderr))
		t.Fatal("Running command error: ", err)
	}
}

func IsPingOK(t *testing.T, containername string, ip string, pnum int) {
	t.Parallel()
	cmd := exec.Command("docker", "exec", containername, "ping", "-c", fmt.Sprintf("%v", pnum), ip)
	cmd.Stdin = os.Stdin
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(stdoutStderr))
		t.Fatal("Running command error: ", err)
	}
	assert.Contains(t, string(stdoutStderr), fmt.Sprintf("%v received", pnum))
	t.Log(string(stdoutStderr))
}

func IsRouteOK(t *testing.T, containername, ip, iphop string) {
	cmd := exec.Command("docker", "exec", containername, "traceroute", ip)
	cmd.Stdin = os.Stdin
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(stdoutStderr))
		t.Fatal("Running command error: ", err)
	}
	assert.Contains(t, string(stdoutStderr), iphop)
	t.Log(string(stdoutStderr))
}

func TestInitial(t *testing.T) {
	InfraBuild(t)
	InfraUp(t)
	time.Sleep(5)

	pingtest := []struct {
		name          string
		containername string
		ip            string
		totalpacket   int
	}{
		{"vpnclient--vpnserver", "vpnclient", "10.5.0.101", 5},
		{"vpnserver--vpnclient", "vpnserver", "10.5.0.102", 6},
		{"externalclient--vpnserver", "externalclient", "10.72.0.101", 5},
		{"vpnserver--externalclient", "vpnserver", "10.72.0.102", 5},
		{"VPN::vpnclient--externalclient", "vpnclient", "10.72.0.102", 7},
		{"VPN::vpnserver--vpnclient", "vpnserver", "192.168.0.102", 4},
		{"VPN::vpnclient--vpnserver", "vpnclient", "10.72.0.101", 5},
		{"VPN::externalclient--vpnclient", "externalclient", "192.168.0.102", 7},
		{"VPN::externalclient--vpnserver", "externalclient", "192.168.0.101", 3},
	}

	for _, tc := range pingtest {
		t.Run(tc.name, func(t *testing.T) {
			IsPingOK(t, tc.containername, tc.ip, tc.totalpacket)
		})
	}

	t.Run("traceroute--external--vpnclient", func(t *testing.T) {
		IsRouteOK(t, "externalclient", "192.168.0.102", "10.72.0.101")
	})
}
