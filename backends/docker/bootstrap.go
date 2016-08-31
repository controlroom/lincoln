package docker

import (
	"fmt"

	"github.com/controlroom/lincoln/interfaces"
)

var rootNetworkName string = "root"
var proxyName string = "root-proxy"
var dnsName string = "root-dns"

func (op DockerOperation) EnsureBootstrapped() {
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
			Image: "jwilder/nginx-proxy",
			Stack: interfaces.Stack{
				Name: rootNetworkName,
			},
			Volumes: []string{"/var/run/docker.sock:/tmp/docker.sock:ro"},
		})
	}

	proxyData, _ := client.ContainerInspect(ctx, proxy.ID)
	proxyIP := proxyData.NetworkSettings.Networks["bridge"].IPAddress
	fmt.Println(proxyIP)

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
			Cmd:          []string{"-S", fmt.Sprintf("/funky.net/%s", proxyIP)},
			PortBindings: []string{"53:53/tcp", "53:53/udp"},
		})
	}
}
