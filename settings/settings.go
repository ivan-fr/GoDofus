package settings

import (
	"gopkg.in/yaml.v2"
	"os"
)

type settings struct {
	Ndc                 string `yaml:"nomdecompte"`
	Pass                string `yaml:"motdepasse"`
	LocalAddress        string `yaml:"localAddress"`
	LocalLoginPort      int32  `yaml:"localLoginPort"`
	LocalGamePort       int32  `yaml:"localGamePort"`
	ServerAnkamaAddress string `yaml:"serverAnkamaAddress"`
	ServerAnkamaPort    int32  `yaml:"serverAnkamaPort"`
}

func getConf() *settings {
	var settings = &settings{}

	yamlFile, err := os.ReadFile("./settings/settings.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, settings)
	if err != nil {
		panic(err)
	}

	return settings
}

var Settings = getConf()
