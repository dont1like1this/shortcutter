package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Name  string
	Links map[string]Link
}

type Link struct {
	Self      string
	Links     map[string]Link
}

type ParameterValue struct {
	Name  string
	Value string
}

func readConfig() Config {
	filename, _ := filepath.Abs("./file.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
}