package main

import (
	"encoding/json"
	"github.com/Bpazy/ssManager/util"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	SsConfPath string `json:"ss_conf_path"`
	Addr       string `json:"addr"`
}

var config Config

func init() {
	configBytes, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatal(err)
	}
}

type SsConfig struct {
	Server       string            `json:"server"`
	LocalAddress string            `json:"local_address"`
	LocalPort    int               `json:"local_port"`
	Timeout      int               `json:"timeout"`
	Method       string            `json:"method"`
	FastOpen     bool              `json:"fast_open"`
	PortPassword map[string]string `json:"port_password"`
}

func findSsConfig() (*SsConfig, error) {
	fileBytes, err := ioutil.ReadFile(config.SsConfPath)
	if err != nil {
		return nil, err
	}
	ssConfig := SsConfig{}
	err = json.Unmarshal(fileBytes, &ssConfig)
	if err != nil {
		return nil, err
	}
	return &ssConfig, nil
}

func saveSsConfig(ssConfig *SsConfig) error {
	ssConfigBytes, err := json.Marshal(ssConfig)
	if err != nil {
		return err
	}

	ssConfigBytes, err = util.PrettyJson(ssConfigBytes)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(config.SsConfPath, ssConfigBytes, os.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}
