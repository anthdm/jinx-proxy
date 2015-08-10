package proxy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Listen   int
	TLS      bool
	CertFile string `json:"cert"`
	KeyFile  string `json:"key"`
	Servers  map[string]*Settings
}

func LoadConfig(file string) *Config {
	p, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	config := Config{}
	if err := json.NewDecoder(bytes.NewBuffer(p)).Decode(&config); err != nil {
		panic(err)
	}
	return &config
}

type Settings struct {
	Proxy  string
	Scheme string
	Path   string
}
