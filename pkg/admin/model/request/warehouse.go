package requestmodel

import (
	"git.selly.red/Selly-Modules/mongodb"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/response"
	"git.selly.red/Selly-Server/warehouse/external/utils/pstring"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	updatemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/update"
)

const (
	BusinessTypeOther = "other"
	BusinessTypeFood  = "food"
)

// WarehouseCreate ...
type WarehouseCreate struct {
	Name         string             `json:"name"`
	BusinessType string             `json:"businessType"`
	Supplier     string             `json:"supplier"`
	Contact      WarehouseContact   `json:"contact"`
	Location     WarehouseLocation  `json:"location"`
	Config       WarehouseCfgCreate `json:"config"`
	SupplierID   primitive.ObjectID `json:"-"`
}

// WarehouseContact ...
type WarehouseContact struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

func (m WarehouseContact) Validate() error {
	return validation.ValidateStruct(
		&m,
		validation.Field(&m.Name, validation.Required.Error(errorcode.WarehouseIsRequiredContactName)),
		validation.Field(&m.Phone, validation.Required.Error(errorcode.WarehouseIsRequiredContactPhone)),
		validation.Field(&m.Address, validation.Required.Error(errorcode.WarehouseIsRequiredContactAddress)),
		validation.Field(&m.Email, validation.Required.Error(errorcode.WarehouseIsRequiredContactEmail)),
	)
}

// WarehouseLocation ...
type WarehouseLocation struct {
	Province            int    `json:"province"`
	District            int    `json:"district"`
	Ward                int    `json:"ward"`
	Address             string `json:"address"`
	FullAddress         string `json:"fullAddress"`
	LocationCoordinates LatLng `json:"locationCoordinates"`
}

func (m WarehouseLocation) Validate() error {
	return validation.ValidateStruct(
		&m,
		validation.Field(&m.Province, validation.Required.Error(errorcode.WarehouseIsRequiredLocationProvince)),
		validation.Field(&m.District, validation.Required.Error(errorcode.WarehouseIsRequiredLocationProvince)),
		validation.Field(&m.Ward, validation.Required.Error(errorcode.WarehouseIsRequiredLocationProvince)),
		validation.Field(&m.Address, validation.Required.Error(errorcode.WarehouseIsRequiredLocationProvince)),
	)
}

// LatLng ...
type LatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Validate ...
func (m WarehouseCreate) Validate() error {
	listBusinessType := []interface{}{
		BusinessTypeFood,
		BusinessTypeOther,
	}
	return validation.ValidateStruct(
		&m,
		validation.Field(&m.Name, validation.Required.Error(errorcode.WarehouseIsRequiredName)),
		validation.Field(&m.Supplier, validation.Required.Error(errorcode.WarehouseIsRequiredSupplier)),
		validation.Field(&m.BusinessType, validation.Required.Error(errorcode.WarehouseCfgIsRequiredBusinessType), validation.In(listBusinessType...).Error(errorcode.WarehouseCfgInvalidBusinessType)),
		validation.Field(&m.Contact),
		validation.Field(&m.Location),
		validation.Field(&m.Config),
	)
}

// ConvertToBSON ...
func (m WarehouseCreate) ConvertToBSON() mgwarehouse.Warehouse {
	return mgwarehouse.Warehouse{
		ID:           mongodb.NewObjectID(),
		BusinessType: m.BusinessType,
		Name:         m.Name,
		SearchString: mongodb.NonAccentVietnamese(m.Name),
		Slug:         pstring.ToSlug(m.Name),
		Status:       constant.StatusInactive,
		Supplier:     m.SupplierID,
		Contact: mgwarehouse.Contact{
			Name:    pstring.TrimSpace(m.Contact.Name),
			Phone:   pstring.TrimSpace(m.Contact.Phone),
			Address: m.Contact.Address,
			Email:   m.Contact.Email,
		},
		Location: mgwarehouse.Location{
			Province: m.Location.Province,
			District: m.Location.District,
			Ward:     m.Location.Ward,
			Address:  m.Location.Address,
			LocationCoordinates: mgwarehouse.Coordinates{
				Type: constant.WarehouseCoordinatesType,
				Coordinates: []float64{
					m.Location.LocationCoordinates.Longitude,
					m.Location.LocationCoordinates.Latitude},
			},
		},
		CreatedAt: ptime.Now(),
		UpdatedAt: ptime.Now(),
	}
}

// WarehouseAll ...
type WarehouseAll struct {
	Page         int64  `query:"page"`
	Limit        int64  `query:"limit"`
	Keyword      string `query:"keyword"`
	BusinessType string `query:"businessType" enums:"food,all,other"`
	Status       string `query:"status"`
	Supplier     string `query:"supplier"`
	Partner      string `query:"partner"`
}

// Validate ...
func (m WarehouseAll) Validate() error {
	return validation.ValidateStruct(
		&m,
		validation.Field(&m.Page, validation.Min(0).Error(response.CommonInvalidPagination)),
		validation.Field(&m.Limit, validation.Min(0).Error(response.CommonInvalidPagination)),
		validation.Field(&m.Supplier, is.MongoID.Error(errorcode.WarehouseInvalidSupplier)),
	)
}

// WarehouseUpdate ...
type WarehouseUpdate struct {
	Name string `json:"name"`
	// BusinessType string            `json:"businessType"`
	Contact  WarehouseContact  `json:"contact"`
	Location WarehouseLocation `json:"location"`
}

// Validate ...
func (m WarehouseUpdate) Validate() error {
	return validation.ValidateStruct(
		&m,
		validation.Field(&m.Name, validation.Required.Error(errorcode.WarehouseIsRequiredName)),
		validation.Field(&m.Contact),
		validation.Field(&m.Location),
	)
}

// ConvertToBSON ...
func (m WarehouseUpdate) ConvertToBSON() updatemodel.WarehouseUpdate {
	return updatemodel.WarehouseUpdate{
		Name:         m.Name,
		SearchString: mongodb.NonAccentVietnamese(m.Name),
		Slug:         pstring.ToSlug(m.Name),
		Contact: updatemodel.WarehouseContact{
			Name:    pstring.TrimSpace(m.Contact.Name),
			Phone:   pstring.TrimSpace(m.Contact.Phone),
			Address: m.Contact.Address,
			Email:   m.Contact.Email,
		},
		Location: updatemodel.LocationWarehouse{
			Province: m.Location.Province,
			District: m.Location.District,
			Ward:     m.Location.Ward,
			Address:  m.Location.Address,
			LocationCoordinates: updatemodel.WarehouseCoordinates{
				Type: constant.WarehouseCoordinatesType,
				Coordinates: []float64{
					m.Location.LocationCoordinates.Longitude,
					m.Location.LocationCoordinates.Latitude,
				},
			},
		},
		UpdatedAt: ptime.Now(),
	}
}

// WarehouseUpdateStatus ...
type WarehouseUpdateStatus struct {
	Status string `json:"status"`
}

// Validate ...
func (m WarehouseUpdateStatus) Validate() error {
	warehouseMethods := []interface{}{
		constant.StatusActive,
		constant.StatusInactive,
	}
	return validation.ValidateStruct(
		&m,
		validation.Field(
			&m.Status,
			validation.Required.Error(errorcode.WarehouseInvalidStatus),
			validation.In(warehouseMethods...).Error(errorcode.WarehouseInvalidStatus),
		),
	)
}
