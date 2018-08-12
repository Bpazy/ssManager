package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type SsConfig struct {
	Server       string            `json:"server"`
	LocalAddress string            `json:"local_address"`
	LocalPort    int               `json:"local_port"`
	Timeout      int               `json:"timeout"`
	Method       string            `json:"method"`
	FastOpen     string            `json:"fast_open"`
	PortPassword map[string]string `json:"port_password"`
}

func main() {
	fileBytes, e := ioutil.ReadFile("./test.json")
	if e != nil {
		log.Fatal(e)
	}
	ssConfig := SsConfig{}
	json.Unmarshal(fileBytes, &ssConfig)
	fmt.Println(ssConfig)

	ssConfig.PortPassword["9999"] = "hahaha"

	ssConfigBytes, err := json.Marshal(ssConfig)
	if err != nil {
		log.Fatal(err)
	}

	ssConfigBytes, err = prettyPrint(ssConfigBytes)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./test2.json", ssConfigBytes, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
}

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
