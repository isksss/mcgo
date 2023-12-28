package config

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/isksss/mcgo/docs"
)

const (
	configFile      = "mcgo.toml"
	ProjectPaper    = "paper"
	ProjectVelocity = "velocity"
)

type Config struct {
	Server  Server   `toml:"server"`
	Plugins []Plugin `toml:"plugins"`
}

type Server struct {
	Project string `toml:"project"`
	Version string `toml:"version"`
	Memory  string `toml:"memory"`
}

type Plugin struct {
	Name string `toml:"name"`
	Url  string `toml:"url"`
}

// ReadConfig is read config file
func GetConfig() (Config, error) {
	var config Config
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func DefaultConfig(project string) {
	var conf []byte
	if project == "velocity" {
		conf = docs.VelocityToml
	} else {
		conf = docs.PaperToml
	}

	// write config file
	f, err := os.Create(configFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(conf)
	if err != nil {
		panic(err)
	}

}

// Config is config struct
func (c Config) Check() error {
	if c.Server.Project == "" {
		return errors.New("project is empty")
	}

	if c.Server.Memory == "" {
		return errors.New("memory is empty")
	}

	// check project
	if err := checkProject(c.Server.Project); err != nil {
		return err
	}
	return nil
}

// 指定されたプロジェクトが存在するか確認する
func checkProject(p string) error {
	j, err := GetJson(apiUrl)
	if err != nil {
		return err
	}

	for _, v := range j.Projects {
		if v == p {
			return nil
		}
	}
	return errors.New("project not found")
}

// 指定されたバージョンが存在するか確認する
func checkVersion(p string, v string) error {
	j, err := GetJson(apiUrl + "/" + p)
	if err != nil {
		return err
	}

	for _, version := range j.Versions {
		// fmt.Println("version: " + version + " v: " + v + "")
		if version == v {
			return nil
		}
	}
	return errors.New("version not found")
}
