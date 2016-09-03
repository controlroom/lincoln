package config

import (
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
