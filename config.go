package main

import (
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type DynamicHost struct {
	Hostname    string `yaml:"host"`
	Login       string `yaml:"login"`
	Password    string `yaml:"passwd"`
	UrlTemplate string `yaml:"url"`
}

type YamlConfig struct {
	Login       string `yaml:"login"`
	Password    string `yaml:"passwd"`
	UrlTemplate string `yaml:"url"`

	Hosts []DynamicHost `yaml:"hosts"`
}

const (
	TEMPLATE_HOST string = "{{HOSTNAME}}"
	TEMPLATE_IP   string = "{{IP}}"
)

func GetDynHosts(filename string) ([]DynamicHost, []bool) {
	source, err := ioutil.ReadFile(filename)
	var config YamlConfig
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}

	updateHost := make([]bool, len(config.Hosts))

	for i := range config.Hosts {
		// Use global config if it is not set for each host specifically
		if config.Hosts[i].Login == "" {
			config.Hosts[i].Login = config.Login
		}
		if config.Hosts[i].Password == "" {
			config.Hosts[i].Password = config.Password
		}
		if config.Hosts[i].UrlTemplate == "" {
			config.Hosts[i].UrlTemplate = config.UrlTemplate
		}
		if !(strings.Contains(config.Hosts[i].UrlTemplate, TEMPLATE_HOST) && strings.Contains(config.Hosts[i].UrlTemplate, TEMPLATE_IP)) {
			log.Warn("Url defined for '" + config.Hosts[i].Hostname + "' doesn't use the templates " + TEMPLATE_HOST + " or " + TEMPLATE_IP + " correctly")
			updateHost[i] = false
		} else {
			updateHost[i] = true
		}
	}

	return config.Hosts, updateHost
}
