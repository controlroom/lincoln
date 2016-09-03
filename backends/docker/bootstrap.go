package docker

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/controlroom/lincoln/interfaces"
)

const (
	rootNetworkName string = "root"
	proxyName       string = "root-proxy"
	dnsName         string = "root-dns"
	rootPort        int    = 4040
	domain          string = "renew"
)

func (op DockerOperation) EnsureBootstrapped() {
	if _, err := os.Stat("/etc/resolver/renew"); os.IsNotExist(err) {
		fmt.Println("Setting up DNS. It might ask for your password")
		cmd := "sudo mkdir -p /etc/resolver && echo 'nameserver 127.0.0.1' | sudo tee /etc/resolver/renew"
		exec.Command("/bin/sh", "-c", cmd).Output()
	}

	if !hasNetwork(rootNetworkName) {
		fmt.Println("Creating root network")

		createNetwork(rootNetworkName, false)
	}

	newProxy := false
	proxy := op.FindContainerByName(proxyName)

	if proxy == nil {
		newProxy = true
		fmt.Println("Creating root proxy")

		proxy = op.StartContainer(interfaces.ContainerStartOptions{
			Name:  proxyName,
			Image: "jwilder/nginx-proxy:0.4.0",
			Stack: interfaces.Stack{
				Name: rootNetworkName,
			},
			Volumes:      []string{"/var/run/docker.sock:/tmp/docker.sock:ro"},
			PortBindings: []string{fmt.Sprintf("%v:80", rootPort)},
		})
	}

	dns := op.FindContainerByName(dnsName)

	if dns == nil || newProxy {
		op.RemoveContainer(dns)
		op.StartContainer(interfaces.ContainerStartOptions{
			Name:   dnsName,
			Image:  "andyshinn/dnsmasq:2.76",
			CapAdd: []string{"NET_ADMIN"},
			Stack: interfaces.Stack{
				Name: "bridge",
			},
			Cmd:          []string{"-S", fmt.Sprintf("/%s/127.0.0.1", domain)},
			PortBindings: []string{"53:53/tcp", "53:53/udp"},
		})
	}
}
