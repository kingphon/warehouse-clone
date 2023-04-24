package service

import (
	"context"

	"git.selly.red/Selly-Modules/natsio/client"
	"git.selly.red/Selly-Modules/natsio/model"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
)

type WarehouseLocation struct {
}

// GetLocationByCode ...
func (wl WarehouseLocation) GetLocationByCode(ctx context.Context, w mgwarehouse.Warehouse) (*model.ResponseLocationAddress, error) {
	body := model.LocationRequestPayload{
		Province: w.Location.Province,
		District: w.Location.District,
		Ward:     w.Location.Ward,
	}
	return client.GetLocation().GetLocationByCode(body)
}
