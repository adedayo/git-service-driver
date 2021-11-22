package model

import (
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

const (
	GIT_SERVICE_CONFIG_FILE = "GitServiceConfig.yaml"
)

type GitServiceType int

const (
	GitHub GitServiceType = iota
	GitLab
)

type GitService struct {
	GraphQLEndPoint, API_Key string
	ID                       string //some unique ID for this service instance
	Type                     GitServiceType
}

type GitServiceConfig struct {
	GitServices map[GitServiceType]map[string]*GitService
	manager     *configManager
}

func (gsc GitServiceConfig) IsServiceConfigured(service GitServiceType) bool {

	log.Printf("%#v", gsc.GitServices)

	if services, present := gsc.GitServices[service]; present && len(services) > 0 {
		return true
	}
	return false
}

func (gsc *GitServiceConfig) AddService(service *GitService) error {
	serviceType := service.Type
	if services, exist := gsc.GitServices[serviceType]; exist {
		services[service.ID] = service
		gsc.GitServices[serviceType] = services
	} else {
		gsc.GitServices[serviceType] = map[string]*GitService{service.ID: service}
	}
	return gsc.manager.SaveConfig(gsc)
}

type configManager struct {
	configLocation string
}

func (cm configManager) GetConfig() *GitServiceConfig {
	conf := &GitServiceConfig{
		GitServices: make(map[GitServiceType]map[string]*GitService),
		manager:     &cm,
	}

	file, err := os.Open(path.Join(cm.configLocation, GIT_SERVICE_CONFIG_FILE))
	if err != nil {
		log.Printf("Error opening Git Service Configuration: %v", err)
		return conf
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(conf); err != nil {
		log.Printf("Error opening Git Service Configuration: %v", err)
	}

	return conf
}

func (cm configManager) SaveConfig(conf *GitServiceConfig) error {
	file, err := os.Create(path.Join(cm.configLocation, GIT_SERVICE_CONFIG_FILE))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	return encoder.Encode(conf)
}

type ConfigManager interface {
	GetConfig() *GitServiceConfig
	SaveConfig(*GitServiceConfig) error
}

func MakeConfigManager() ConfigManager {
	location := "."
	if loc, err := homedir.Expand("~/.checkmate/config"); err == nil {
		location = loc
	}

	cm := configManager{
		configLocation: location,
	}

	//attempt to create the project location if it doesn't exist
	os.MkdirAll(location, 0755)

	return &cm
}
