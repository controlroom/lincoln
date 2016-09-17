package docker

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
)

func getImage(imageTag string) {
	matches, _ := client.ImageList(ctx, types.ImageListOptions{
		MatchName: imageTag,
	})

	if len(matches) == 0 || strings.HasSuffix(imageTag, "latest") {
		fmt.Printf("Check image pull %v\n", imageTag)
		resp, err := client.ImagePull(ctx, imageTag, types.ImagePullOptions{})
		defer resp.Close()

		if err != nil {
			panic(err)
		}

		jsonmessage.DisplayJSONMessagesStream(resp, os.Stderr, os.Stderr.Fd(), true, nil)
	}
}
