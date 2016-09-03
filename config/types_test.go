package config

import "testing"

func TestBasicNode(t *testing.T) {
	yaml := "nodes:\n  test: here is a test command"
	config := ParseConfig([]byte(yaml))
	if config.GetNodes()[0].Cmd != "here is a test command" {
		t.Fail()
	}
}

func TestNestedNode(t *testing.T) {
	yaml := `nodes:
  test:
    cmd: here is a test command`
	config := ParseConfig([]byte(yaml))
	if config.GetNodes()[0].Cmd != "here is a test command" {
		t.Fail()
	}
}
