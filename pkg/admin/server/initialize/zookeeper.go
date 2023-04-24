package initialize

import (
	"fmt"
	"os"

	zk "git.selly.red/Selly-Modules/zookeeper"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
)

// zookeeper ...
func zookeeper() {
	var (
		uri = os.Getenv("ZOOKEEPER_URI")
	)

	// Connect
	if err := zk.Connect(uri); err != nil {
		panic(err)
	}

	envVars := config.GetENV()
	externalValues(envVars)
	commonValues(envVars)
}

// get value in zookeeper
func externalValues(envVars *config.ENV) {
	// MongoDB
	mongodbPrefix := getExternalPrefix("/mongodb/warehouse")
	envVars.MongoDB.URI = zk.GetStringValue(fmt.Sprintf("%s/uri", mongodbPrefix))
	envVars.MongoDB.DBName = zk.GetStringValue(fmt.Sprintf("%s/db_name", mongodbPrefix))

	envVars.MongoDB.ReplicaSet = zk.GetStringValue(fmt.Sprintf("%s/replica_set", mongodbPrefix))
	envVars.MongoDB.CAPem = zk.GetStringValue(fmt.Sprintf("%s/ca_pem", mongodbPrefix))
	envVars.MongoDB.CertPem = zk.GetStringValue(fmt.Sprintf("%s/cert_pem", mongodbPrefix))
	envVars.MongoDB.CertKeyFilePassword = zk.GetStringValue(fmt.Sprintf("%s/cert_key_file_password", mongodbPrefix))
	envVars.MongoDB.ReadPrefMode = zk.GetStringValue(fmt.Sprintf("%s/read_pref_mode", mongodbPrefix))

	// NATS
	natsPrefix := getExternalPrefix("/nats/warehouse")
	envVars.Nats.URL = zk.GetStringValue(fmt.Sprintf("%s/uri", natsPrefix))
	envVars.Nats.Username = zk.GetStringValue(fmt.Sprintf("%s/user", natsPrefix))
	envVars.Nats.Password = zk.GetStringValue(fmt.Sprintf("%s/password", natsPrefix))
	envVars.Nats.APIKey = zk.GetStringValue(fmt.Sprintf("%s/api_key", natsPrefix))

	// MongoDB_AUDIT
	mongoAuditPrefix := getExternalPrefix("/mongodb/warehouse_audit")
	envVars.MongoAudit.Host = zk.GetStringValue(fmt.Sprintf("%s/host", mongoAuditPrefix))
	envVars.MongoAudit.DBName = zk.GetStringValue(fmt.Sprintf("%s/db_name", mongoAuditPrefix))

	// Redis value
	redisPrefix := getExternalPrefix("/redis/warehouse")
	envVars.Redis.URI = zk.GetStringValue(redisPrefix + "/uri")
	envVars.Redis.PWD = zk.GetStringValue(redisPrefix + "/password")

	// Telegram
	//telegramPrefix := getExternalPrefix("/telegram/warehouse")
	//envVars.Telegram.Token = zk.GetStringValue(fmt.Sprintf("%s/token", telegramPrefix))
	//envVars.Telegram.WarehouseID = zk.GetIntValue(fmt.Sprintf("%s/warehouseID", telegramPrefix))

}

func getExternalPrefix(group string) string {
	return fmt.Sprintf("%s%s", config.GetENV().ZookeeperPrefixExternal, group)
}

func commonValues(envVars *config.ENV) {
	envVars.MaxTimeWaitForGCUpdateStatusInMinute = zk.GetIntValue(getCommonPrefix("/max_time_wait_for_global_care_update_status_in_minute"))
}

func getCommonPrefix(group string) string {
	return fmt.Sprintf("%s%s", config.GetENV().ZookeeperPrefixWarehouseCommon, group)
}

// func adminValues(envVars *config.ENV) {
//	// For admin values
// }

// func getAdminPrefix(group string) string {
//	return fmt.Sprintf("%s%s", config.GetENV().ZookeeperLocationAdminPrefix, group)
// }
