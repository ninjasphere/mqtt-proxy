package store

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ninjablocks/mqtt-proxy/conf"
)

type MysqlStore struct {
	db   *sql.DB
	conf *conf.MysqlConfiguration
}

func NewMysqlStore(conf *conf.MysqlConfiguration) *MysqlStore {
	db, err := sql.Open("mysql", conf.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	//	defer db.Close()
	return &MysqlStore{
		db:   db,
		conf: conf,
	}
}

// Sends a PING request to Redis.
func (s *MysqlStore) Health() bool {
	//	defer s.db.Close()
	err := s.db.Ping()
	if err != nil {
		return false
	}
	return true
}

// Validates the credentials against MySQL.
func (s *MysqlStore) FindUser(token string) (*User, error) {

	var uid uint
	var mqttId string
	err := s.db.QueryRow(s.conf.Select, token).Scan(&uid, &mqttId)
	if err != nil {
		return nil, err
	}
	return &User{
		UserId: uid,
		MqttId: mqttId,
	}, nil
}

// Sends a PING request to Redis.
func (s *MysqlStore) Close() {
	s.db.Close()
}
