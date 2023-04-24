package responsemodel

import (
	"git.selly.red/Selly-Server/warehouse/external/constant"
	"time"

	"git.selly.red/Selly-Modules/natsio/model"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
)

// ResponseWarehouseAll ...
type ResponseWarehouseAll struct {
	List  []WarehouseBrief `json:"list"`
	Total int64            `json:"total"`
	Limit int64            `json:"limit"`
}

// WarehouseBrief ...
type WarehouseBrief struct {
	ID           string                    `json:"_id"`
	Name         string                    `json:"name"`
	BusinessType string                    `json:"businessType"`
	Status       string                    `json:"status"`
	Slug         string                    `json:"slug"`
	Supplier     WarehouseSupplier         `json:"supplier"`
	Location     ResponseWarehouseLocation `json:"location"`
	Contact      ResponseWarehouseContact  `json:"contact"`
	CreatedAt    *ptime.TimeResponse       `json:"createdAt"`
	UpdatedAt    *ptime.TimeResponse       `json:"updatedAt"`
}

type WarehouseSupplier struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

// ResponseWarehouseContact ...
type ResponseWarehouseContact struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

// ResponseWarehouseLocation ...
type ResponseWarehouseLocation struct {
	Province            ResponseWarehouseProvince `json:"province"`
	District            ResponseWarehouseDistrict `json:"district"`
	Ward                ResponseWarehouseWard     `json:"ward"`
	Address             string                    `json:"address"`
	LocationCoordinates ResponseLatLng            `json:"locationCoordinates"`
}

type ResponseWarehouseProvince struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code int    `json:"code"`
}
type ResponseWarehouseDistrict struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code int    `json:"code"`
}
type ResponseWarehouseWard struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code int    `json:"code"`
}

// ResponseLatLng ...
type ResponseLatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ResponseWarehouseDetail ...
type ResponseWarehouseDetail struct {
	ID                    string                    `json:"_id"`
	BusinessType          string                    `json:"businessType"`
	Name                  string                    `json:"name"`
	Supplier              WarehouseSupplier         `json:"supplier"`
	Location              ResponseWarehouseLocation `json:"location"`
	Contact               ResponseWarehouseContact  `json:"contact"`
	CreatedAt             *ptime.TimeResponse       `json:"createdAt"`
	UpdatedAt             *ptime.TimeResponse       `json:"updatedAt"`
	ReasonPendingInactive string                    `json:"reasonPendingInactive"`
	SupplierHolidayFrom   *ptime.TimeResponse       `json:"supplierHolidayFrom"`
	SupplierHolidayTo     *ptime.TimeResponse       `json:"supplierHolidayTo"`
	Status                string                    `json:"status"`
	StatusBeforeHoliday   string                    `json:"statusBeforeHoliday"`
}

// WarehouseNatsResponse ...
type WarehouseNatsResponse struct {
	ID             string                               `json:"_id"`
	Name           string                               `json:"name"`
	SearchString   string                               `json:"searchString"`
	Slug           string                               `json:"slug"`
	Status         string                               `json:"status"`
	Supplier       string                               `json:"supplier"`
	Contact        ResponseWarehouseContact             `json:"contact"`
	Location       ResponseWarehouseLocation            `json:"location"`
	Configurations ResponseWarehouseConfigurationDetail `json:"configurations"`
	CreatedAt      time.Time                            `json:"createdAt"`
	UpdatedAt      time.Time                            `json:"updatedAt"`
}

func GetWarehouseNatsResponse(w mgwarehouse.Warehouse, c mgwarehouse.Configuration, staff string) model.WarehouseNatsResponse {
	res := model.WarehouseNatsResponse{
		ID:           w.ID.Hex(),
		Staff:        staff,
		BusinessType: w.BusinessType,
		Name:         w.Name,
		SearchString: w.SearchString,
		Slug:         w.Slug,
		Status:       w.Status,
		Supplier:     w.Supplier.Hex(),
		Contact: model.ResponseWarehouseContact{
			Name:    w.Contact.Name,
			Phone:   w.Contact.Phone,
			Address: w.Contact.Address,
			Email:   w.Contact.Email,
		},
		Location: model.ResponseWarehouseLocation{
			Province: model.CommonLocation{
				Code: w.Location.Province,
			},
			District: model.CommonLocation{
				Code: w.Location.District,
			},
			Ward: model.CommonLocation{
				Code: w.Location.Ward,
			},
			Address:             w.Location.Address,
			LocationCoordinates: model.ResponseLatLng{},
		},
		Configurations: model.WarehouseConfiguration{
			Warehouse: c.Warehouse.Hex(),
			Supplier: model.WarehouseSupplier{
				CanAutoSendMail:       false,
				InvoiceDeliveryMethod: c.Supplier.InvoiceDeliveryMethod,
			},
			Order: model.WarehouseOrder{
				MinimumValue: c.Order.MinimumValue,
				PaymentMethod: model.WarehousePaymentMethod{
					Cod:          c.Order.PaymentMethod.Cod,
					BankTransfer: c.Order.PaymentMethod.BankTransfer,
				},
				IsLimitNumberOfPurchases: c.Order.IsLimitNumberOfPurchases,
				LimitNumberOfPurchases:   c.Order.LimitNumberOfPurchases,
			},
			Partner: model.WarehousePartner{
				IdentityCode:   c.Partner.IdentityCode,
				Code:           c.Partner.Code,
				Enabled:        c.Partner.Enabled,
				Authentication: c.Partner.Authentication,
			},
			Delivery: model.WarehouseDelivery{
				DeliveryMethods:      c.Delivery.DeliveryMethods,
				PriorityServiceCodes: c.Delivery.PriorityServiceCodes,
				EnabledSources:       c.Delivery.EnabledSources,
				Types:                c.Delivery.Types,
			},
			Other: model.WarehouseOther{
				DoesSupportSellyExpress: c.Other.DoesSupportSellyExpress,
			},
			Food: model.WarehouseFood{
				ForceClosed: c.Food.ForceClosed,
				IsClosed:    c.Food.IsClosed,
				TimeRange:   make([]model.TimeRange, 0),
			},
			AutoConfirmOrder: model.WarehouseOrderConfirm{
				IsEnable:              c.AutoConfirmOrder.IsEnable,
				ConfirmDelayInSeconds: c.AutoConfirmOrder.ConfirmDelayInSeconds,
			},
		},
		CreatedAt:             w.CreatedAt,
		UpdatedAt:             w.UpdatedAt,
		ReasonPendingInactive: w.ReasonPendingInactive,
		IsPendingInactive:     false,
	}

	if w.Status == constant.WarehouseStatusHoliday {
		res.IsPendingInactive = true
	}

	if len(c.Food.TimeRange) > 0 {
		for _, t := range c.Food.TimeRange {
			res.Configurations.Food.TimeRange = append(res.Configurations.Food.TimeRange, model.TimeRange{
				From: t.From,
				To:   t.To,
			})
		}
	}

	if len(w.Location.LocationCoordinates.Coordinates) == 2 {
		res.Location.LocationCoordinates = model.ResponseLatLng{
			Latitude:  w.Location.LocationCoordinates.Coordinates[1],
			Longitude: w.Location.LocationCoordinates.Coordinates[0],
		}
	}
	return res
}

// ResponseWarehouseShort ...
type ResponseWarehouseShort struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}
