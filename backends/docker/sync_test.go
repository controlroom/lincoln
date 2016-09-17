package docker

import (
	"fmt"
	"testing"
)

func TestFreePort(t *testing.T) {
	port := freePort()
	fmt.Println(fmt.Sprintf("%v", port))
}
