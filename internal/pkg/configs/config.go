package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CassandraFirstNode  string
	CassandraSecondNode string
	CassandraThirdNode  string
	CassandraKeySpace   string
	SecretKey           string
	AppPort             string
}

func LoadConnectionConfig() (*Config, error) {
	var config Config
	err := godotenv.Load("./internal/pkg/configs/connection.env")
	if err != nil {
		return nil, err
	}

	config.SecretKey = os.Getenv("SECRET_KEY")
	config.CassandraFirstNode = os.Getenv("CASSANDRA_NODE1")
	config.CassandraSecondNode = os.Getenv("CASSANDRA_NODE2")
	config.CassandraThirdNode = os.Getenv("CASSANDRA_NODE3")
	config.CassandraKeySpace = os.Getenv("CASSANDRA_KEYSPACE")
	config.AppPort = os.Getenv("APP_PORT")
	return &config, err
}
