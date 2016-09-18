package docker

import "github.com/docker/docker/api/types"

// CreateVolume
// ensure named volume is created
func (op DockerOperation) CreateVolume(name string) {
	client.VolumeCreate(
		ctx,
		types.VolumeCreateRequest{
			Name: name,
		},
	)
}
