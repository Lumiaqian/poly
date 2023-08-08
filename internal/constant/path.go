package constant

import (
	"fmt"
	"os"
	P "path"
)

const (
	Name  = "poly"
	Focus = "focus.yaml"
)

type path struct {
	homeDir    string
	configFile string
}

// Path is used to get the configuration path
var Path = func() *path {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir, _ = os.Getwd()
	}

	homeDir = P.Join(homeDir, ".config", Name)
	return &path{homeDir: homeDir, configFile: "config.yaml"}
}()

// SetHomeDir is used to set the configuration path
func SetHomeDir(root string) {
	Path.homeDir = root
}

// SetConfig is used to set the configuration file
func SetConfig(file string) {
	Path.configFile = file
}

func (p *path) HomeDir() string {
	return p.homeDir
}

func (p *path) Config() string {
	return p.configFile
}

func (p *path) Focus() string {
	return fmt.Sprintf("%s/%s", p.homeDir, Focus)
}
