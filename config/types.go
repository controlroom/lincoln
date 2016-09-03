package config

import (
	"log"

	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v2"
)

type Node struct {
	Name     string
	Cmd      string
	Replicas int
}

type Config struct {
	Nodes map[string]interface{}
}

func GetConfig(fileData []byte) *Config {
	t := Config{}

	err := yaml.Unmarshal(fileData, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &t
}

func (config *Config) GetNodes() []Node {
	var nodes []Node
	for k, v := range config.Nodes {
		var res Node
		switch v.(type) {
		case string:
			cmd, _ := v.(string)
			res.Cmd = cmd
		default:
			mapstructure.Decode(v, &res)
		}

		res.Name = k
		nodes = append(nodes, res)
	}

	return nodes
}
