package config

import (
	"reflect"
	"testing"
)

func TestGetNodes(t *testing.T) {
	cases := []struct {
		yaml     string
		expected Node
	}{
		{
			yaml: "nodes:\n  test: here is a test command",
			expected: Node{
				Name:    "test",
				Exposed: true,
				Cmd:     "here is a test command",
			},
		},
		{
			yaml: `nodes:
  test:
    cmd: here is a test command`,
			expected: Node{
				Name:    "test",
				Exposed: true,
				Cmd:     "here is a test command",
			},
		},
	}

	for _, c := range cases {
		node := ParseConfig([]byte(c.yaml)).GetNodes()["test"]
		if !reflect.DeepEqual(node, c.expected) {
			t.Errorf("Expected \n%+v\nto be equal to\n%+v", node, c.expected)
		}
	}
}
