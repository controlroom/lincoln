package docker

import (
	"context"

	dockerClient "github.com/docker/engine-api/client"
)

var ctx = context.Background()
var client *dockerClient.Client = getClient()

func getClient() *dockerClient.Client {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := dockerClient.NewClient("unix:///var/run/docker.sock", "v1.23", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}

	return cli
}
