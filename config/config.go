package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type App struct {
	Branch string
	Path   string
	Config *Config
}

func allLocalMatches(path string) []App {
	matches, _ := filepath.Glob(fmt.Sprintf("%v/*/lincoln.yml", path))
	apps := make([]App, len(matches))
	for i, match := range matches {
		os.Chdir(filepath.Dir(match))
		branch, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "head").Output()
		apps[i] = App{
			Branch: string(branch),
			Path:   filepath.Dir(match),
			Config: ParseConfigFromPath(match),
		}
	}

	return apps
}

func FindLocalApp(path string, name string) (*App, error) {
	matches := allLocalMatches(path)

	for _, match := range matches {
		if match.Config.Name == name {
			return &match, nil
		}
	}

	return nil, errors.New("Could not find app")
}

func FindAllLocalApps(path string) []App {
	return allLocalMatches(path)
}
