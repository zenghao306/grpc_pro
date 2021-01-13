package etc

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type (
	Server struct {
		HttpPort  int    `yaml:"http_port"`
		RunMode   string `yaml:"run_mode"`
		EtcdAddr  string `yaml:"etcd_addr"`
	}

	LogInfo struct {
		LogPath    string `yaml:"log_path"`
		LogAdapter string `yaml:"log_adapter"`
		LogLevel   string `yaml:"log_level"`
	}

	config struct {
		Server   *Server     `yaml:"server"`
		LogInfo  *LogInfo    `yaml:"log_info"`
	}
)

var Conf config

func InitConfig(file string) error {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bs, &Conf)
}