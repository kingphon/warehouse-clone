package service

import (
	"context"

	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/response"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/request"
)

// SupplierHolidayInterface ...
type SupplierHolidayInterface interface {
	// CreateWithClientData ...
	CreateWithClientData(ctx context.Context, payload requestmodel.SupplierHolidayCreate) (doc mgwarehouse.SupplierHoliday, err error)

	// Detail ...
	Detail(ctx context.Context) (result *responsemodel.ResponseSupplierHolidayDetail)

	// FindByCondition ...
	FindByCondition(ctx context.Context, cond interface{}) []mgwarehouse.SupplierHoliday

	// UpdateWithClientData ...
	UpdateWithClientData(ctx context.Context, payload requestmodel.SupplierHolidayUpdate) (result responsemodel.ResponseUpdate, err error)

	// ChangeStatus ...
	ChangeStatus(ctx context.Context, payload requestmodel.SupplierHolidayChangeStatus) (result responsemodel.ResponseChangeStatus, err error)
}

// supplierHolidayImplement ...
type supplierHolidayImplement struct {
	CurrentUser *responsemodel.ResponseUserInfo
}

func SupplierHoliday(user *responsemodel.ResponseUserInfo) SupplierHolidayInterface {
	return &supplierHolidayImplement{
		CurrentUser: user,
	}
}
