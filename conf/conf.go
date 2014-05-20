package conf

import (
	"errors"
	"io/ioutil"
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

type MysqlConfiguration struct {
	ConnectionString string `toml:"connection-string"`
	Select           string `toml:"select"`
}

type MqttConfiguration struct {
	ListenAddress string `toml:"listen-address"`
	Cert          string `toml:"cert"`
	Key           string `toml:"key"`
}

type InfluxConfiguration struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Pass     string `toml:"pass"`
	Database string `toml:"database"`
}

type LibratoConfiguration struct {
	Email string `toml:"email"`
	Token string `toml:"token"`
}

type Configuration struct {
	BackendServers []string `toml:"backend-servers"`
	User           string   `toml:"user"`
	Pass           string   `toml:"pass"`

	// typically us-west | us-east
	// prepended to metrics
	Region string `toml:"region"`

	// typically develop | beta | prod
	// prepended to metrics
	Environment string `toml:"env"`

	ReadTimeout int `toml:"read-timeout"`

	MqttStoreMysql MysqlConfiguration   `toml:"mqtt-store"`
	Mqtt           MqttConfiguration    `toml:"mqtt"`
	Influx         InfluxConfiguration  `toml:"influx"`
	Librato        LibratoConfiguration `toml:"librato"`
}

func (c *Configuration) GetReadTimeout() time.Duration {
	return time.Second * time.Duration(c.ReadTimeout)
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

	if tomlConfiguration.Region == "" {
		tomlConfiguration.Region = "us-east"
	}

	if tomlConfiguration.Environment == "" {
		tomlConfiguration.Environment = "develop"
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

	// need a way to merge defaults..
	if tomlConfiguration.MqttStoreMysql.ConnectionString == "" {
		tomlConfiguration.MqttStoreMysql.ConnectionString = "root:@tcp(127.0.0.1:3306)/mqtt"
	}

	if tomlConfiguration.MqttStoreMysql.Select == "" {
		tomlConfiguration.MqttStoreMysql.Select = "select uid, mqtt_id from users where mqtt_id = ?"
	}

	return tomlConfiguration, nil
}
