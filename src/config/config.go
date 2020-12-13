package config

import "errors"

type Configs []Config

type Config struct {
	//ServiceName is unique
	ServiceName string                 `json:"serviceName"`
	Reference   map[string]interface{} `json:"reference"`
	Param       map[string]interface{} `json:"param"`
	Result      map[string]interface{} `json:"result"`
}

//hold configs in memory
var configs Configs

/**
find match config by serviceName
*/
func LoadAndGetConfigByServiceName(serviceName string) (*Config, error) {
	//load configs from file if not in memory
	if len(configs) == 0 {
		loader := NewLoader(FILE)
		data, err := loader.Load()
		if err != nil {
			return nil, err
		}
		configs = data
	}
	config, err := findConfigByServiceName(serviceName)
	if err != nil {
		return nil, err
	}
	return config, nil
}

/**
loop configs to find config by serviceName
*/
func findConfigByServiceName(serviceName string) (*Config, error) {
	//to check if a slice is empty, always use len(s) == 0. Do not check for nil
	if len(configs) == 0 {
		return nil, errors.New("Configs is nil")
	}
	for _, config := range configs {
		if config.ServiceName == serviceName {
			return &config, nil
		}
	}
	return nil, errors.New("can't find config by serviceName")
}
