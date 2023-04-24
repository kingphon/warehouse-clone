package requestmodel

import (
	"fmt"
	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/convert"
	"git.selly.red/Selly-Server/warehouse/external/utils/parray"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"git.selly.red/Selly-Server/warehouse/pkg/app/errorcode"
	"git.selly.red/Selly-Server/warehouse/pkg/app/locale"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/app/model/response"
	"github.com/friendsofgo/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SupplierHolidayCreate ...
type SupplierHolidayCreate struct {
	From       string   `json:"from"`
	To         string   `json:"to"`
	Warehouses []string `json:"warehouses"`
	IsApplyAll bool     `json:"isApplyAll"`
}

func (m SupplierHolidayCreate) Validate() error {
	if len(m.Warehouses) == 0 {
		return errors.New(errorcode.SupplierHolidayMusHaveAtLeastOneWarehouse)
	}

	return validation.ValidateStruct(&m,
		validation.Field(&m.Warehouses, validation.Each(is.MongoID.Error(errorcode.WarehouseInvalidID))),
	)
}

// ConvertToBSON ...
func (m SupplierHolidayCreate) ConvertToBSON(supplier responsemodel.ResponseSupplierInfo) mgwarehouse.SupplierHoliday {
	result := mgwarehouse.SupplierHoliday{
		ID:         mongodb.NewObjectID(),
		Supplier:   mongodb.ConvertStringToObjectID(supplier.ID),
		Title:      locale.GetSupplierHolidayTitle(supplier.Name),
		From:       ptime.TimeParseISODate(m.From),
		To:         ptime.TimeParseISODate(m.To),
		Source:     constant.WarehouseSupplierHolidaySourceSupplier,
		Status:     constant.StatusActive,
		Warehouses: make([]primitive.ObjectID, 0),
		CreatedAt:  ptime.Now(),
		UpdatedAt:  ptime.Now(),
		CreatedBy:  mongodb.ConvertStringToObjectID(supplier.ID),
		IsApplyAll: m.IsApplyAll,
	}

	result.SearchString = mongodb.NonAccentVietnamese(result.Title)
	result.Reason = locale.GetSupplierHolidayReason(result.To.Format("02/01/2006"))

	if len(m.Warehouses) > 0 && !m.IsApplyAll {
		result.Warehouses, _ = convert.StringsToObjectIDs(parray.UniqueArrayStrings(m.Warehouses))
	}

	return result
}

// SupplierHolidayUpdate ...
type SupplierHolidayUpdate struct {
	From       string   `json:"from"`
	To         string   `json:"to"`
	Warehouses []string `json:"warehouses"`
	IsApplyAll bool     `json:"isApplyAll"`
}

func (m SupplierHolidayUpdate) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Warehouses, validation.Each(is.MongoID.Error(errorcode.WarehouseInvalidID))),
	)
}

// ConvertToBSON ...
func (m SupplierHolidayUpdate) ConvertToBSON() mgwarehouse.SupplierHoliday {
	result := mgwarehouse.SupplierHoliday{
		From:       ptime.TimeParseISODate(m.From),
		To:         ptime.TimeParseISODate(m.To),
		Status:     constant.StatusActive,
		UpdatedAt:  ptime.Now(),
		IsApplyAll: m.IsApplyAll,
		Warehouses: make([]primitive.ObjectID, 0),
	}

	result.Reason = fmt.Sprintf("Tạm ngưng bán đến ngày %s", result.To.Format("02/01/2006"))

	if len(m.Warehouses) > 0 {
		result.Warehouses, _ = convert.StringsToObjectIDs(parray.UniqueArrayStrings(m.Warehouses))
	}

	return result
}

// SupplierHolidayChangeStatus ...
type SupplierHolidayChangeStatus struct {
	Status string `json:"status"`
}

func (m SupplierHolidayChangeStatus) Validate() error {
	statuses := []interface{}{
		constant.StatusInactive,
		constant.StatusActive,
	}

	return validation.ValidateStruct(&m,
		validation.Field(&m.Status, validation.In(statuses...).Error(errorcode.SupplierHolidayInvalidStatus)),
	)
}

// SupplierHolidayPayload ...
type SupplierHolidayPayload struct {
	Create *SupplierHolidayCreate
	Update *SupplierHolidayUpdate
}
