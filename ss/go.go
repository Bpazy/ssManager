package ss

import (
	"encoding/json"
	"github.com/Bpazy/ssManager/util"
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
	util.ShouldPanic(err)

	portPasswords := config.PortPassword
	for k := range portPasswords {
		portPasswords[k] = "******"
	}
	return portPasswords
}

func (c GoClient) AddPortPassword(port, password string) {
	config, err := loadConfig(c.Filename)
	util.ShouldPanic(err)
	config.PortPassword[port] = password
	writeConfig(c.Filename, config)
}

func (c GoClient) DeletePort(port string) error {
	config, err := loadConfig(c.Filename)
	util.ShouldPanic(err)

	delete(config.PortPassword, port)
	writeConfig(c.Filename, config)

	return nil
}

func (GoClient) Restart() error {
	panic("implement me")
}

func loadConfig(filename string) (c *GoConfig, err error) {
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
	util.ShouldPanic(err)

	err = ioutil.WriteFile(filename, configBytes, 0644)
	util.ShouldPanic(err)
}
