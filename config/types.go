package config

import (
	"io/ioutil"
	"log"

	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v2"
)

type Node struct {
	Name      string
	Cmd       string
	Exposed   bool
	SubDomain string `yaml:"sub-domain"`
	Balanced  bool
	Replicas  int
}

type Config struct {
	Name          string
	Description   string
	Github        string
	DefaultBranch string `yaml:"default-branch"`
	Host          string
	DevImage      string   `yaml:"dev-image"`
	SharedPaths   []string `yaml:"shared-paths"`

	Cmds     map[string]string
	Nodes    map[string]interface{}
	NodeSets map[string][]string `yaml:"node-sets"`
}

type App struct {
	Branch string
	Path   string
	Config *Config
}

// Extract yaml file data([]byte) into Config struct
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

// Allows for multiple shapes of yaml nodes
//
// Yaml may look like:
//
// nodes:
//   node1: bundle exec command
//
// Or it might look like
//
// nodes:
//   node1:
//		 cmd: bundle exec command
//     exposed: true
//
// Which is the default with only one command. Both commands are essentially
// equal.
func (config *Config) GetNodes() map[string]Node {
	var nodes = make(map[string]Node)
	for k, v := range config.Nodes {
		var res Node
		switch v.(type) {
		case string:
			cmd, _ := v.(string)
			res.Cmd = cmd
		default:
			mapstructure.Decode(v, &res)
		}

		if len(config.Nodes) == 1 {
			res.Exposed = true
		}

		res.Name = k
		nodes[k] = res
	}

	return nodes
}
