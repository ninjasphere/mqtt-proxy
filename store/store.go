package store

import "errors"

var ErrUserNotFound = errors.New("User not found in store")

type Store interface {
	Health() bool
	FindUser(token string) (*User, error)
	Close()
}

type User struct {
	UserId uint   `json:"uid"`
	MqttId string `json:"mqttId"`
}
