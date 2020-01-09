package config

//go:generate pflags Config --default-var=defaultConfig

import (
	"context"
	"net/url"

	"github.com/lyft/flytestdlib/config"
	"github.com/lyft/flytestdlib/logger"

	pluginsConfig "github.com/lyft/flyteplugins/go/tasks/config"
)

const quboleConfigSectionKey = "qubole"

func MustParse(s string) config.URL {
	r, err := url.Parse(s)
	if err != nil {
		logger.Panicf(context.TODO(), "Bad Qubole URL Specified as default, error: %s", err)
	}
	if r == nil {
		logger.Panicf(context.TODO(), "Nil Qubole URL specified.", err)
	}
	return config.URL{URL: *r}
}

type ClusterConfig struct {
	PrimaryLabel string   `json:"primaryLabel" pflag:",The primary label of a given service cluster"`
	Labels       []string `json:"labels" pflag:",Labels of a given service cluster"`
	Limit        int      `json:"limit" pflag:",Resource quota (in the number of outstanding requests) of the service cluster"`
}

type ProjectDestinationClusterConfig struct {
	Project      string `json:"project" pflag:",Project of the task which the query belongs to"`
	Domain       string `json:"domain" pflag:",Domain of the task which the query belongs to"`
	ClusterLabel string `json:"clusterLabel" pflag:",The label of the destination cluster this query to be submitted to"`
}

var (
	defaultConfig = Config{
		Endpoint:                         MustParse("https://wellness.qubole.com"),
		CommandAPIPath:                   MustParse("/api/v1.2/commands/"),
		AnalyzeLinkPath:                  MustParse("/v2/analyze"),
		TokenKey:                         "FLYTE_QUBOLE_CLIENT_TOKEN",
		LruCacheSize:                     2000,
		Workers:                          15,
		ClusterConfigs:                   []ClusterConfig{{PrimaryLabel: "default", Labels: []string{"default"}, Limit: 250}},
		ProjectDestinationClusterConfigs: []ProjectDestinationClusterConfig{},
	}

	quboleConfigSection = pluginsConfig.MustRegisterSubSection(quboleConfigSectionKey, &defaultConfig)
)

// Qubole plugin configs
type Config struct {
	Endpoint                         config.URL                        `json:"endpoint" pflag:",Endpoint for qubole to use"`
	CommandAPIPath                   config.URL                        `json:"commandApiPath" pflag:",API Path where commands can be launched on Qubole. Should be a valid url."`
	AnalyzeLinkPath                  config.URL                        `json:"analyzeLinkPath" pflag:",URL path where queries can be visualized on qubole website. Should be a valid url."`
	TokenKey                         string                            `json:"quboleTokenKey" pflag:",Name of the key where to find Qubole token in the secret manager."`
	LruCacheSize                     int                               `json:"lruCacheSize" pflag:",Size of the AutoRefreshCache"`
	Workers                          int                               `json:"workers" pflag:",Number of parallel workers to refresh the cache"`
	ClusterConfigs                   []ClusterConfig                   `json:"clusterConfigs" pflag:"-,A list of cluster configs. Each of the configs corresponds to a service cluster"`
	ProjectDestinationClusterConfigs []ProjectDestinationClusterConfig `json:"projectDestinationClusterConfigs" pflag:"-,A list configs specifying the destination service cluster for (project, domain)"`
}

// Retrieves the current config value or default.
func GetQuboleConfig() *Config {
	return quboleConfigSection.GetConfig().(*Config)
}

func SetQuboleConfig(cfg *Config) error {
	return quboleConfigSection.SetConfig(cfg)
}
