package service

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
)

//
// PUBLIC METHODS
//

// DetailByWarehouseID ...
func (s warehouseConfigurationImplement) DetailByWarehouseID(ctx context.Context, id primitive.ObjectID) (*responsemodel.ResponseWarehouseConfigurationDetail, error) {
	var (
		d    = dao.WarehouseConfiguration()
		cond = bson.M{"warehouse": id}
	)

	warehouseCfg := d.FindOneByCondition(ctx, cond)
	if warehouseCfg.ID.IsZero() {
		return nil, errors.New(errorcode.WarehouseCfgNotFound)
	}

	return s.detail(ctx, warehouseCfg), nil
}

//
// PRIVATE METHODS
//

// detail ...
func (s warehouseConfigurationImplement) detail(ctx context.Context, d mgwarehouse.Configuration) *responsemodel.ResponseWarehouseConfigurationDetail {
	res := &responsemodel.ResponseWarehouseConfigurationDetail{
		ID:        d.ID.Hex(),
		Warehouse: d.Warehouse.Hex(),
		Supplier: responsemodel.ResponseWarehouseSupplierConfig{
			InvoiceDeliveryMethod: d.Supplier.InvoiceDeliveryMethod,
		},
		Order: responsemodel.ResponseWarehouseOrderConfig{
			MinimumValue: d.Order.MinimumValue,
			PaymentMethod: responsemodel.ResponseWarehousePaymentMethodConfig{
				Cod:          d.Order.PaymentMethod.Cod,
				BankTransfer: d.Order.PaymentMethod.BankTransfer,
			},
			IsLimitNumberOfPurchases: d.Order.IsLimitNumberOfPurchases,
			LimitNumberOfPurchases:   d.Order.LimitNumberOfPurchases,
			NotifyOnNewOrder:         d.Order.NotifyOnNewOrder,
		},
		Partner: responsemodel.ResponseWarehousePartnerConfig{
			IdentityCode:   d.Partner.IdentityCode,
			Code:           d.Partner.Code,
			Enabled:        d.Partner.Enabled,
			Authentication: d.Partner.Authentication,
		},
		Delivery: responsemodel.ResponseWarehouseDeliveryConfig{
			DeliveryMethods:      d.Delivery.DeliveryMethods,
			PriorityServiceCodes: d.Delivery.PriorityServiceCodes,
			EnabledSources:       d.Delivery.EnabledSources,
			Types:                d.Delivery.Types,
		},
		Other: responsemodel.ResponseWarehouseOtherConfig{
			DoesSupportSellyExpress: d.Other.DoesSupportSellyExpress,
		},
		Food: responsemodel.ResponseWarehousePartnerFood{
			IsClosed:    d.Food.ForceClosed || d.Food.IsClosed,
			ForceClosed: d.Food.ForceClosed,
			TimeRange:   make([]responsemodel.TimeRange, 0),
		},
		AutoConfirmOrder: d.AutoConfirmOrder,
	}
	if len(d.Food.TimeRange) > 0 {
		for _, t := range d.Food.TimeRange {
			res.Food.TimeRange = append(res.Food.TimeRange, responsemodel.TimeRange{
				From: t.From,
				To:   t.To,
			})
		}
	}
	return res
}
