package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

/**
default config loader
*/

type FileLoader struct {
	name string
	path string
}

/**
read config from file and parse to map
*/
func (fl *FileLoader) Load() (Configs, error) {
	//path and name divied I think better
	fileName := filepath.Join(fl.path, fl.name)
	_, err := os.Stat(fileName)
	if err != nil {
		return nil, errors.New("file not exist")
	}
	return readConfig(fileName)
}

func readConfig(fileName string) (Configs, error) {
	//read from file
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("error happen when read file")
	}

	config := Configs{}
	//parse to json
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.New("error happen when parse file cotent")
	}
	return config, nil
}
