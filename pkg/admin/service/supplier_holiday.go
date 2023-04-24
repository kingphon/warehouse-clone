package service

import (
	"context"

	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/mgquery"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SupplierHolidayInterface ...
type SupplierHolidayInterface interface {
	// CreateWithClientData ...
	CreateWithClientData(ctx context.Context, payload requestmodel.SupplierHolidayCreate) (doc mgwarehouse.SupplierHoliday, err error)

	// All ...
	All(ctx context.Context, q mgquery.AppQuery) (docs responsemodel.ResponseSupplierHolidayAll)

	// Detail ...
	Detail(ctx context.Context, id primitive.ObjectID) (result *responsemodel.ResponseSupplierHolidayDetail, err error)

	// UpdateWithClientData ...
	UpdateWithClientData(ctx context.Context, id primitive.ObjectID, payload requestmodel.SupplierHolidayUpdate) (result responsemodel.ResponseUpdate, err error)

	// ChangeStatus ...
	ChangeStatus(ctx context.Context, id primitive.ObjectID, payload requestmodel.SupplierHolidayChangeStatus) (result responsemodel.ResponseChangeStatus, err error)

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}) []mgwarehouse.SupplierHoliday

	// UpdateHolidayStatusForSupplier ...
	UpdateHolidayStatusForSupplier()
}

// supplierHolidayImplement ...
type supplierHolidayImplement struct {
	CurrentStaff externalauth.User
}

func SupplierHoliday(staff externalauth.User) SupplierHolidayInterface {
	return &supplierHolidayImplement{
		CurrentStaff: staff,
	}
}
