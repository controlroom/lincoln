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

func FindLocalApp(path string, name string) (*App, error) {
	appPath := fmt.Sprintf("%v/%v", path, name)
	configPath := fmt.Sprintf("%v/lincoln.yml", appPath)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("%s does not exist or does not contain a lincoln.yml", name))
	}

	return &App{
		Path:   appPath,
		Config: ParseConfigFromPath(configPath),
	}, nil
}

func FindAllLocalApps(path string) []App {
	var apps []App
	matches, _ := filepath.Glob(fmt.Sprintf("%v/*/lincoln.yml", path))
	if len(matches) > 0 {
		for _, fpath := range matches {
			os.Chdir(filepath.Dir(fpath))
			branch, _ := exec.Command("git", "rev-parse", "--abbrev-ref", "head").Output()
			apps = append(apps, App{
				Branch: string(branch),
				Path:   fpath,
				Config: ParseConfigFromPath(fpath),
			})
		}
	}
	return apps
}
