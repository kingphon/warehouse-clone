package requestmodel

import (
	"git.selly.red/Selly-Modules/mongodb"
	externalauth "git.selly.red/Selly-Server/warehouse/external/auth"
	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/convert"
	"git.selly.red/Selly-Server/warehouse/external/utils/parray"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SupplierHolidayCreate ...
type SupplierHolidayCreate struct {
	Supplier   string   `json:"supplier"`
	Title      string   `json:"title"`
	From       string   `json:"from"`
	To         string   `json:"to"`
	Reason     string   `json:"reason"`
	Warehouses []string `json:"warehouses"`
	IsApplyAll bool     `json:"isApplyAll"`
}

func (m SupplierHolidayCreate) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Supplier, is.MongoID.Error(errorcode.SupplierHolidayInvalidSupplierId)),
		validation.Field(&m.Title, validation.Required.Error(errorcode.SupplierHolidayIsRequireTitle)),
		validation.Field(&m.Reason, validation.Required.Error(errorcode.SupplierHolidayIsRequireReason)),
		validation.Field(&m.Warehouses, validation.Each(is.MongoID.Error(errorcode.WarehouseInvalidID))),
	)
}

// ConvertToBSON ...
func (m SupplierHolidayCreate) ConvertToBSON(staff externalauth.User) mgwarehouse.SupplierHoliday {
	result := mgwarehouse.SupplierHoliday{
		ID:           mongodb.NewObjectID(),
		Title:        m.Title,
		From:         ptime.TimeParseISODate(m.From),
		To:           ptime.TimeParseISODate(m.To),
		Reason:       m.Reason,
		SearchString: mongodb.NonAccentVietnamese(m.Title),
		Source:       constant.WarehouseSupplierHolidaySourceAdmin,
		Status:       constant.StatusInactive,
		Warehouses:   make([]primitive.ObjectID, 0),
		IsApplyAll:   m.IsApplyAll,
		Supplier:     mongodb.ConvertStringToObjectID(m.Supplier),
		CreatedBy:    mongodb.ConvertStringToObjectID(staff.ID),
		CreatedAt:    ptime.Now(),
		UpdatedAt:    ptime.Now(),
	}

	if len(m.Warehouses) > 0 {
		result.Warehouses, _ = convert.StringsToObjectIDs(parray.UniqueArrayStrings(m.Warehouses))
	}

	return result
}

// SupplierHolidayAll ...
type SupplierHolidayAll struct {
	Limit     int64  `query:"limit"`
	Page      int64  `query:"page"`
	Keyword   string `query:"keyword"`
	Status    string `query:"status"`
	FromAt    string `query:"fromAt"`
	ToAt      string `query:"toAt"`
	Supplier  string `query:"supplier"`
	Warehouse string `query:"warehouse"`
}

func (m SupplierHolidayAll) Validate() error {
	return validation.ValidateStruct(&m)
}

// SupplierHolidayUpdate ...
type SupplierHolidayUpdate struct {
	Title      string   `json:"title"`
	From       string   `json:"from"`
	To         string   `json:"to"`
	Supplier   string   `json:"supplier"`
	IsApplyAll bool     `json:"isApplyAll"`
	Warehouses []string `json:"warehouses"`
	Reason     string   `json:"reason"`
}

func (m SupplierHolidayUpdate) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Supplier, is.MongoID.Error(errorcode.SupplierHolidayInvalidSupplierId)),
		validation.Field(&m.Title, validation.Required.Error(errorcode.SupplierHolidayIsRequireTitle)),
		validation.Field(&m.Reason, validation.Required.Error(errorcode.SupplierHolidayIsRequireReason)),
		validation.Field(&m.Warehouses, validation.Each(is.MongoID.Error(errorcode.WarehouseInvalidID))),
	)
}

// ConvertToBSON ...
func (m SupplierHolidayUpdate) ConvertToBSON() mgwarehouse.SupplierHoliday {
	result := mgwarehouse.SupplierHoliday{
		Title:        m.Title,
		From:         ptime.TimeParseISODate(m.From),
		To:           ptime.TimeParseISODate(m.To),
		Reason:       m.Reason,
		SearchString: mongodb.NonAccentVietnamese(m.Title),
		UpdatedAt:    ptime.Now(),
		IsApplyAll:   m.IsApplyAll,
		Warehouses:   make([]primitive.ObjectID, 0),
		Supplier:     mongodb.ConvertStringToObjectID(m.Supplier),
	}

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
