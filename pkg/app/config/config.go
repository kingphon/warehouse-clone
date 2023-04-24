package config

import (
	"fmt"
	"os"

	"git.selly.red/Selly-Server/warehouse/external/constant"

	"git.selly.red/Selly-Modules/mongodb"
)

// ENV ...
type ENV struct {
	Env                            string
	ZookeeperPrefixExternal        string
	ZookeeperPrefixWarehouseCommon string
	ZookeeperPrefixWarehouseApp    string
	MongoDB                        struct {
		URI    string
		DBName string

		ReplicaSet          string
		CAPem               string
		CertPem             string
		CertKeyFilePassword string
		ReadPrefMode        string
	}
	Nats       NatsConfig
	SecretKey  string
	MongoAudit MongoConfig `env:",prefix=MONGO_AUDIT_"`
}

// NatsConfig
type NatsConfig struct {
	URL      string
	Username string
	Password string
	APIKey   string
}

// GetConnectOptions ...
func (dbCfg MongoConfig) GetConnectOptions() mongodb.Config {
	return mongodb.Config{
		Host:       dbCfg.Host,
		DBName:     dbCfg.DBName,
		Standalone: &mongodb.ConnectStandaloneOpts{},
		TLS:        &mongodb.ConnectTLSOpts{},
	}
}

// Audit
type MongoConfig struct {
	Host   string
	DBName string
}

var env ENV

// GetENV ...
func GetENV() *ENV {
	return &env
}

// IsEnvDevelop ...
func IsEnvDevelop() bool {
	return env.Env == constant.EnvDevelop
}

// IsEnvStaging ...
func IsEnvStaging() bool {
	return env.Env == constant.EnvStaging
}

// IsEnvProduction ...
func IsEnvProduction() bool {
	return env.Env == constant.EnvProduction
}

// Init ...
func Init() {
	env = ENV{
		Env:                            os.Getenv("ENV"),
		ZookeeperPrefixExternal:        os.Getenv("ZOOKEEPER_PREFIX_EXTERNAL"),
		ZookeeperPrefixWarehouseCommon: os.Getenv("ZOOKEEPER_PREFIX_WAREHOUSE_COMMON"),
		ZookeeperPrefixWarehouseApp:    os.Getenv("ZOOKEEPER_PREFIX_WAREHOUSE_APP"),
	}

	fmt.Println("env", env)
}
