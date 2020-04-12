package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type InfluxDBConfig struct {
	Addr string
	Name string
	User string
	Pwd  string
}

type Config struct {
	Port int             `yaml:"port"`
	DB   *InfluxDBConfig `yaml:"influxdb"`
}

var cfg = &Config{}

func initConfig() {
	cfgFile := os.Args[1]
	err := load(cfgFile, cfg)
	if err != nil {
		log.Printf("Load configure file \"%s\" failed :%s\n", cfgFile, err)
		os.Exit(-1)
	}
}

func load(file string, out interface{}) error {
	conf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return loadData(conf, out)
}

func loadData(data []byte, out interface{}) error {
	if err := yaml.Unmarshal(data, out); err != nil {
		return err
	}

	return nil
}
