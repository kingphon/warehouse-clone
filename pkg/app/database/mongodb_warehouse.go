package database

import (
	"git.selly.red/Selly-Modules/mongodb"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	"git.selly.red/Selly-Server/warehouse/pkg/app/config"
)

var warehouseDB *mongo.Database

// ConnectMongoDBWarehouse ...
func ConnectMongoDBWarehouse() {
	var (
		cfg = config.GetENV().MongoDB
		err error
		tls *mongodb.ConnectTLSOpts
	)
	if cfg.ReplicaSet != "" {
		tls = &mongodb.ConnectTLSOpts{
			ReplSet:             cfg.ReplicaSet,
			CaFile:              cfg.CAPem,
			CertKeyFile:         cfg.CertPem,
			CertKeyFilePassword: cfg.CertKeyFilePassword,
			ReadPreferenceMode:  cfg.ReadPrefMode,
		}
	}

	// Connect
	warehouseDB, err = mongodb.Connect(mongodb.Config{
		Host:    cfg.URI,
		DBName:  cfg.DBName,
		Monitor: apmmongo.CommandMonitor(),
		TLS:     tls,
		Standalone: &mongodb.ConnectStandaloneOpts{
			AuthMechanism: "",
			AuthSource:    "",
			Username:      "",
			Password:      "",
		},
	})
	if err != nil {
		panic(err)
	}
}

// GetMongoDBWarehouse ...
func GetMongoDBWarehouse() *mongo.Database {
	return warehouseDB
}

// WarehouseCol ...
func WarehouseCol() *mongo.Collection {
	return warehouseDB.Collection(constant.ColWarehouse)
}

// WarehouseConfigurationCol ...
func WarehouseConfigurationCol() *mongo.Collection {
	return warehouseDB.Collection(constant.ColWarehouseConfiguration)
}

// SupplierHolidayCol ...
func SupplierHolidayCol() *mongo.Collection {
	return warehouseDB.Collection(colSupplierHoliday)
}
