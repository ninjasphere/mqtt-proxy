package conf

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

type MysqlConfiguration struct {
	ConnectionString string `toml:"connection-string"`
	MqttSelect       string `toml:"mqtt-select"`
}

type HttpConfiguration struct {
	ListenAddress string `toml:"listen-address"`
}

type MqttConfiguration struct {
	ListenAddress string `toml:"listen-address"`
	Cert          string `toml:"cert"`
	Key           string `toml:"key"`
}

type Configuration struct {
	BackendServers []string `toml:"backend-servers"`
	User           string   `toml:"user"`
	Pass           string   `toml:"pass"`

	Http  HttpConfiguration  `toml:"http"`
	Mysql MysqlConfiguration `toml:"mysql"`
	Mqtt  MqttConfiguration  `toml:"mqtt"`
}

func LoadConfiguration(fileName string) *Configuration {
	config, err := parseTomlConfiguration(fileName)
	if err != nil {
		log.Println("Couldn't parse configuration file: " + fileName)
		panic(err)
	}
	return config
}

func parseTomlConfiguration(filename string) (*Configuration, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	tomlConfiguration := &Configuration{}
	_, err = toml.Decode(string(body), tomlConfiguration)
	if err != nil {
		return nil, err
	}
	if len(tomlConfiguration.BackendServers) == 0 {
		return nil, errors.New("At least one backend servers required.")
	}
	if tomlConfiguration.Http.ListenAddress == "" {
		tomlConfiguration.Http.ListenAddress = ":9000"
	}
	if tomlConfiguration.Mqtt.ListenAddress == "" {
		tomlConfiguration.Mqtt.ListenAddress = ":1883"
	}
	if tomlConfiguration.User == "" {
		tomlConfiguration.User = "guest"
	}
	if tomlConfiguration.Pass == "" {
		tomlConfiguration.Pass = "guest"
	}

	if tomlConfiguration.Mysql.ConnectionString == "" {
		tomlConfiguration.Mysql.ConnectionString = "root:@tcp(127.0.0.1:3306)/mqtt"
	}
	if tomlConfiguration.Mysql.MqttSelect == "" {
		tomlConfiguration.Mysql.MqttSelect = "select uid, mqtt_id from users where mqtt_id = ?"
	}

	return tomlConfiguration, nil
}
