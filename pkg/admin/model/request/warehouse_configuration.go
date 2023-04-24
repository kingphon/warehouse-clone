package requestmodel

import (
	"errors"
	"sort"
	"time"

	"git.selly.red/Selly-Modules/mongodb"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	updatemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/update"
)

// WarehouseCfgCreate ...
type WarehouseCfgCreate struct {
	DoesSupportSellyExpress bool                           `json:"doesSupportSellyExpress"`
	Supplier                WarehouseCfgSupplier           `json:"supplier"`
	Order                   WarehouseCfgOrder              `json:"order"`
	Partner                 WarehouseCfgPartner            `json:"partner"`
	Delivery                WarehouseCfgDelivery           `json:"delivery"`
	Food                    WarehouseCfgFood               `json:"food"`
	OrderConfirm            mgwarehouse.ConfigOrderConfirm `json:"orderConfirm"`
}

// WarehouseCfgOrder ...
type WarehouseCfgOrder struct {
	MinimumValue             float64                            `json:"minimumValue"`
	PaymentMethod            WarehouseCfgPaymentMethod          `json:"paymentMethod"`
	IsLimitNumberOfPurchases bool                               `json:"isLimitNumberOfPurchases" json:"isLimitNumberOfPurchases"`
	LimitNumberOfPurchases   int64                              `json:"limitNumberOfPurchases" json:"limitNumberOfPurchases"`
	NotifyOnNewOrder         mgwarehouse.ConfigNotifyOnNewOrder `json:"notifyOnNewOrder"`
}

func (m WarehouseCfgFood) IsClosed() bool {
	if len(m.TimeRange) == 0 {
		return false
	}
	now := ptime.TimeOfDayInHCM(time.Now())
	timeStart := ptime.TimeStartOfDayInHCM(time.Now())
	minute := now.Sub(timeStart).Minutes()

	var (
		isClose = true
	)
	for _, rangeTime := range m.TimeRange {
		if rangeTime.From <= int64(minute) && rangeTime.To >= int64(minute) {
			isClose = false
			break
		}
	}
	return isClose
}

func (m WarehouseCfgOrder) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.LimitNumberOfPurchases, validation.When(m.IsLimitNumberOfPurchases, validation.Required.Error(errorcode.WarehouseIsRequiredLimitNumberOfPurchases))),
	)
}

// WarehouseCfgFood ...
type WarehouseCfgFood struct {
	ForceClosed bool        `json:"forceClosed"`
	TimeRange   []TimeRange `json:"timeRange"`
}

type TimeRange struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

// WarehouseCfgSupplier ...
type WarehouseCfgSupplier struct {
	InvoiceDeliveryMethod string `bson:"invoiceDeliveryMethod"`
}

func (m WarehouseCfgSupplier) Validate() error {
	listInvoiceMethod := []interface{}{
		constant.WarehouseInvoiceMethodNone,
		constant.WarehouseInvoiceMethodInstant,
		constant.WarehouseInvoiceMethodAfterDelivered,
	}

	return validation.ValidateStruct(
		&m,
		validation.Field(&m.InvoiceDeliveryMethod,
			validation.Required.Error(errorcode.WarehouseCfgInvalidInvoiceMethod),
			validation.In(listInvoiceMethod...).Error(errorcode.WarehouseCfgInvalidInvoiceMethod)),
	)
}

// WarehouseCfgPaymentMethod ...
type WarehouseCfgPaymentMethod struct {
	Cod          bool `json:"cod"`
	BankTransfer bool `json:"bankTransfer"`
}

// WarehouseCfgDelivery ...
type WarehouseCfgDelivery struct {
	DeliveryMethods      []string `json:"deliveryMethods"`
	PriorityServiceCodes []string `json:"priorityServiceCodes"`
	EnabledSources       []int    `json:"enabledSources"`
	Types                []string `json:"types"`
}

func (m WarehouseCfgDelivery) Validate() error {
	deliveryMethods := []interface{}{
		constant.WarehouseDeliveryInnerProvince,
		constant.WarehouseDeliveryOuterProvince,
	}
	deliveryTypes := []interface{}{
		constant.WarehouseCfgDeliveryTypeSelf,
		constant.WarehouseCfgDeliveryTypeSelly,
	}

	return validation.ValidateStruct(
		&m,
		validation.Field(&m.DeliveryMethods,
			validation.Required.Error(errorcode.WarehouseCfgIsRequiredDeliveryMethods),
			validation.Each(validation.In(deliveryMethods...).Error(errorcode.WarehouseCfgInvalidDeliveryMethod), validation.Required.Error(errorcode.WarehouseCfgIsRequiredDeliveryMethods))),
		validation.Field(&m.Types,
			validation.Each(validation.In(deliveryTypes...).Error(errorcode.WarehouseCfgInvalidDeliveryType))),
		validation.Field(&m.EnabledSources),
		validation.Field(&m.Types, validation.Required.Error(errorcode.WarehouseCfgIsRequiredTypes)),
	)
}

// WarehouseCfgPartner ...
type WarehouseCfgPartner struct {
	IdentityCode   string `json:"identityCode"`
	Code           string `json:"code"`
	Enabled        bool   `json:"enabled"`
	Authentication string `json:"authentication"`
}

// Validate ...
func (m WarehouseCfgCreate) Validate() error {
	return validation.ValidateStruct(
		&m,
		validation.Field(&m.Supplier),
		validation.Field(&m.Partner),
		validation.Field(&m.Delivery),
		validation.Field(&m.Order),
	)
}

// WarehouseCfgOther ...
type WarehouseCfgOther struct {
	DoesSupportSellyExpress bool `json:"doesSupportSellyExpress"`
}

func (m WarehouseCfgOther) Validate() error {
	return validation.ValidateStruct(&m)
}

// ConvertToBSON ...
func (m WarehouseCfgCreate) ConvertToBSON(id primitive.ObjectID) mgwarehouse.Configuration {
	res := mgwarehouse.Configuration{
		ID:        mongodb.NewObjectID(),
		Warehouse: id,
		Food: mgwarehouse.ConfigFood{
			ForceClosed: m.Food.ForceClosed,
			IsClosed:    m.Food.IsClosed() || m.Food.ForceClosed,
			TimeRange:   make([]mgwarehouse.TimeRange, 0),
		},
		Other: mgwarehouse.ConfigOther{
			DoesSupportSellyExpress: m.DoesSupportSellyExpress,
		},
		Supplier: mgwarehouse.ConfigSupplier{
			InvoiceDeliveryMethod: m.Supplier.InvoiceDeliveryMethod,
		},
		Order: mgwarehouse.ConfigOrder{
			MinimumValue: m.Order.MinimumValue,
			PaymentMethod: mgwarehouse.ConfigPaymentMethod{
				Cod:          m.Order.PaymentMethod.Cod,
				BankTransfer: m.Order.PaymentMethod.BankTransfer,
			},
			IsLimitNumberOfPurchases: m.Order.IsLimitNumberOfPurchases,
			LimitNumberOfPurchases:   m.Order.LimitNumberOfPurchases,
			NotifyOnNewOrder:         m.Order.NotifyOnNewOrder,
		},
		Partner: mgwarehouse.ConfigPartner{
			IdentityCode:   m.Partner.IdentityCode,
			Code:           m.Partner.Code,
			Enabled:        m.Partner.Enabled,
			Authentication: m.Partner.Authentication,
		},
		Delivery: mgwarehouse.ConfigDelivery{
			DeliveryMethods:      m.Delivery.DeliveryMethods,
			PriorityServiceCodes: m.Delivery.PriorityServiceCodes,
			EnabledSources:       m.Delivery.EnabledSources,
			Types:                m.Delivery.Types,
		},
		AutoConfirmOrder: m.OrderConfirm,
	}

	if len(m.Food.TimeRange) > 0 {
		for _, t := range m.Food.TimeRange {
			res.Food.TimeRange = append(res.Food.TimeRange, mgwarehouse.TimeRange{
				From: t.From,
				To:   t.To,
			})
		}
	}
	return res
}

// WarehouseCfgFoodUpdate ...
type WarehouseCfgFoodUpdate struct {
	Food WarehouseCfgFood `json:"food"`
}

func (m WarehouseCfgFoodUpdate) Validate() error {
	if len(m.Food.TimeRange) > 0 {
		sort.Slice(m.Food.TimeRange, func(i, j int) bool {
			return m.Food.TimeRange[i].From < m.Food.TimeRange[j].From
		})
		var (
			to int64
		)
		for _, tRange := range m.Food.TimeRange {
			if tRange.From > tRange.To {
				return errors.New(errorcode.WarehouseCfgTimeRangeInvalid)
			}
			if to > tRange.From {
				return errors.New(errorcode.WarehouseCfgTimeRangeInvalid)
			}
			to = tRange.To
		}
	}
	return validation.ValidateStruct(&m)
}

// ConvertTOBSON ...
func (m WarehouseCfgFoodUpdate) ConvertTOBSON() updatemodel.WarehouseCfgFoodUpdate {
	res := updatemodel.WarehouseCfgFoodUpdate{
		ForceClosed: m.Food.ForceClosed,
		TimeRange:   make([]updatemodel.TimeRange, 0),
	}
	if len(m.Food.TimeRange) > 0 {
		for _, timeR := range m.Food.TimeRange {
			res.TimeRange = append(res.TimeRange, updatemodel.TimeRange{
				From: timeR.From,
				To:   timeR.To,
			})
		}
	}
	res.IsClosed = m.Food.IsClosed() || m.Food.ForceClosed
	return res
}

// WarehouseCfgSupplierUpdate ...
type WarehouseCfgSupplierUpdate struct {
	Supplier WarehouseCfgSupplier `json:"supplier"`
}

func (m WarehouseCfgSupplierUpdate) Validate() error {
	return validation.ValidateStruct(&m, validation.Field(&m.Supplier))
}

func (m WarehouseCfgSupplierUpdate) ConvertTOBSON() updatemodel.WarehouseCfgSupplierUpdate {
	return updatemodel.WarehouseCfgSupplierUpdate{
		InvoiceDeliveryMethod: m.Supplier.InvoiceDeliveryMethod,
	}
}

// WarehouseCfgPartnerUpdate ...
type WarehouseCfgPartnerUpdate struct {
	Partner WarehouseCfgPartner `json:"partner"`
}

func (m WarehouseCfgPartnerUpdate) Validate() error {
	return validation.ValidateStruct(&m, validation.Field(&m.Partner))
}

func (m WarehouseCfgPartnerUpdate) ConvertToBSON() updatemodel.WarehouseCfgPartnerUpdate {
	return updatemodel.WarehouseCfgPartnerUpdate{
		IdentityCode:   m.Partner.IdentityCode,
		Code:           m.Partner.Code,
		Enabled:        m.Partner.Enabled,
		Authentication: m.Partner.Authentication,
	}
}

// WarehouseCfgOrderUpdate ...
type WarehouseCfgOrderUpdate struct {
	Order WarehouseCfgOrder `json:"order"`
}

func (m WarehouseCfgOrderUpdate) Validate() error {
	return validation.ValidateStruct(&m, validation.Field(&m.Order))
}

func (m WarehouseCfgOrderUpdate) ConvertToBSON() updatemodel.WarehouseCfgOrderUpdate {
	return updatemodel.WarehouseCfgOrderUpdate{
		MinimumValue: m.Order.MinimumValue,
		PaymentMethod: updatemodel.WarehouseCfgPaymentMethodUpdate{
			Cod:          m.Order.PaymentMethod.Cod,
			BankTransfer: m.Order.PaymentMethod.BankTransfer,
		},
		IsLimitNumberOfPurchases: m.Order.IsLimitNumberOfPurchases,
		LimitNumberOfPurchases:   m.Order.LimitNumberOfPurchases,
		NotifyOnNewOrder:         m.Order.NotifyOnNewOrder,
	}
}

// WarehouseCfgDeliveryUpdate ...
type WarehouseCfgDeliveryUpdate struct {
	Delivery WarehouseCfgDelivery `json:"delivery"`
}

func (m WarehouseCfgDeliveryUpdate) Validate() error {
	return validation.ValidateStruct(&m, validation.Field(&m.Delivery))
}

func (m WarehouseCfgDeliveryUpdate) ConvertToBSON() updatemodel.WarehouseCfgDeliveryUpdate {
	return updatemodel.WarehouseCfgDeliveryUpdate{
		DeliveryMethods:      m.Delivery.DeliveryMethods,
		PriorityServiceCodes: m.Delivery.PriorityServiceCodes,
		EnabledSources:       m.Delivery.EnabledSources,
		Types:                m.Delivery.Types,
	}
}

// WarehouseCfgOtherUpdate ...
type WarehouseCfgOtherUpdate struct {
	Other WarehouseCfgOther `json:"other"`
}

func (m WarehouseCfgOtherUpdate) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Other),
	)
}

func (m WarehouseCfgOtherUpdate) ConvertToBSON() updatemodel.WarehouseCfgOtherUpdate {
	return updatemodel.WarehouseCfgOtherUpdate{
		DoesSupportSellyExpress: m.Other.DoesSupportSellyExpress,
	}
}

// WarehouseCfgOrderConfirm ...
type WarehouseCfgOrderConfirm struct {
	OrderConfirm mgwarehouse.ConfigOrderConfirm `json:"orderConfirm"`
}

// Validate ...
func (m WarehouseCfgOrderConfirm) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderConfirm),
	)
}
