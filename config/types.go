package config

import (
	"io/ioutil"
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
	Name          string
	Description   string
	Github        string
	DefaultBranch string "yaml:default-branch"
	Host          string

	Nodes map[string]interface{}
}

func ParseConfig(fileData []byte) *Config {
	t := Config{}

	err := yaml.Unmarshal(fileData, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &t
}

func ParseConfigFromPath(path string) *Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return ParseConfig(file)
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
