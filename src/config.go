package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	ServiceBasicInfo serviceBasicInfo
	LogConfig        logConfig
	Topics           map[string]topicSetting
}

type serviceBasicInfo struct {
	Cpu           int
	ZookeeperAddr string
	EsAddr        string
	DataSource    string
}
type logConfig struct {
	LogPath string
}
type topicSetting struct {
	UdpAlarm    int
	PushBuffNum int
}

const configPath = "../conf/config.toml"

var configInfo config

//  set cpu
func configCPUS() int {
	if _, err := toml.DecodeFile(configPath, &configInfo); err != nil {
		fmt.Println(err)
	}
	return configInfo.ServiceBasicInfo.Cpu
}

// kafka zookeeper address
func zookeeperAddr() *string {
	if s, err := toml.DecodeFile(configPath, &configInfo); err != nil {
		fmt.Println(err)
		fmt.Println(s)
	}
	if configInfo.ServiceBasicInfo.ZookeeperAddr == "" {
		fmt.Println("zookeeper is empty ! please check it !")
		os.Exit(2)
	}
	return &configInfo.ServiceBasicInfo.ZookeeperAddr
}
func configDataSource() string {
	if _, err := toml.DecodeFile(configPath, &configInfo); err != nil {
		fmt.Println(err)
	}
	return configInfo.ServiceBasicInfo.DataSource
}
func logPath() string {
	if _, err := toml.DecodeFile(configPath, &configInfo); err != nil {
		fmt.Println(err)
	}
	return configInfo.LogConfig.LogPath
}
func topicInfo() map[string]topicSetting {
	if _, err := toml.DecodeFile(configPath, &configInfo); err != nil {
		fmt.Println(err)
	}
	return configInfo.Topics
}
func esAddr() string {
	if _, err := toml.DecodeFile(configPath, &configInfo); err != nil {
		fmt.Println(err)

	}
	return configInfo.ServiceBasicInfo.EsAddr
}

// func main() {
// 	fmt.Println(esAddr())
// 	// fmt.Println(configDataSource())
// }
