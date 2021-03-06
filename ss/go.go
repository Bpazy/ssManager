package ss

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type GoClient struct {
	Filename string
}

type GoConfig struct {
	Server       string            `json:"server"`
	LocalAddress string            `json:"local_address"`
	LocalPort    int               `json:"local_port"`
	Timeout      int               `json:"timeout"`
	Method       string            `json:"method"`
	FastOpen     bool              `json:"fast_open"`
	PortPassword map[string]string `json:"port_password"`
}

func (c GoClient) QueryPortPasswords() map[string]string {
	config, err := loadConfig(c.Filename)
	if err != nil {
		panic(err)
	}

	portPasswords := config.PortPassword
	for k := range portPasswords {
		portPasswords[k] = "******"
	}
	return portPasswords
}

func (c GoClient) AddPortPassword(port, password string) {
	config, err := loadConfig(c.Filename)
	if err != nil {
		panic(err)
	}

	config.PortPassword[port] = password
	writeConfig(c.Filename, config)
}

func (c GoClient) DeletePort(port string) error {
	config, err := loadConfig(c.Filename)
	if err != nil {
		panic(err)
	}

	delete(config.PortPassword, port)
	writeConfig(c.Filename, config)

	return nil
}

func (GoClient) Restart() error {
	panic("implement me")
}

func loadConfig(filename string) (c *GoConfig, err error) {
	logrus.Debug(filename)
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileBytes, &c)
	if err != nil {
		return nil, err
	}
	return
}

func writeConfig(filename string, c *GoConfig) {
	configBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filename, configBytes, 0644)
	if err != nil {
		panic(err)
	}
}
