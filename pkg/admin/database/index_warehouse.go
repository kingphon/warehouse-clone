package database

import (
	"context"
	"fmt"

	"git.selly.red/Selly-Modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func indexWarehouse() {
	// outbound request
	{
		index := []mongo.IndexModel{
			mongodb.NewIndexKey("partner.code"),
			mongodb.NewIndexKey("orderCode"),
			mongodb.NewIndexKey("trackingCode"),
			mongodb.NewIndexKey("order", "status"),
		}
		indexCol(OutboundRequestCol(), index)
	}

	// outbound request history
	{
		index := []mongo.IndexModel{
			mongodb.NewIndexKey("request"),
		}
		indexCol(OutboundRequestHistoryCol(), index)
	}

	// WarehouseConfiguration
	{
		index := []mongo.IndexModel{
			mongodb.NewIndexKey("warehouse"),
		}
		indexCol(WarehouseConfigurationCol(), index)
	}
}

func indexCol(col *mongo.Collection, indexes []mongo.IndexModel) {
	_, err := col.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		fmt.Printf("Index collection %s err: %v", col.Name(), err)
	}
}
