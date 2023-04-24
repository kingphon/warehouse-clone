package service

import (
	"context"
	"fmt"
	"time"

	"git.selly.red/Selly-Modules/mongodb"

	"git.selly.red/Selly-Server/warehouse/external/constant"

	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"

	"git.selly.red/Selly-Modules/natsio/client"

	"git.selly.red/Selly-Modules/natsio/model"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/thoas/go-funk"

	"git.selly.red/Selly-Modules/logger"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"

	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
)

type Schedule struct {
}

// RunJobUpdateIsClosed ...
func (s Schedule) RunJobUpdateIsClosed() {
	var (
		ctx                = context.Background()
		wModel             []mongo.WriteModel
		listSupplier       = make([]string, 0)
		listUpdateSupplier = make([]model.SupplierIsClosed, 0)
	)
	warehouses := dao.Warehouse().FindByCondition(ctx, bson.M{
		"businessType": requestmodel.BusinessTypeFood,
	})
	for _, wh := range warehouses {
		cfg := dao.WarehouseConfiguration().FindByWarehouseID(ctx, wh.ID)
		if cfg.ID.IsZero() || cfg.Food.ForceClosed {
			continue
		}
		isClosed := CheckIsClosed(cfg.Food)
		if isClosed != cfg.Food.IsClosed {
			cond := bson.M{
				"_id": cfg.ID,
			}
			update := bson.M{
				"$set": bson.M{
					"food.isClosed": isClosed,
				},
			}
			wModel = append(wModel, mongo.NewUpdateOneModel().SetFilter(cond).SetUpdate(update))
			if funk.Contains(listSupplier, wh.Supplier.Hex()) {
				listSupplier = append(listSupplier, wh.Supplier.Hex())
			}
		}
	}
	if len(wModel) > 0 {
		if err := dao.WarehouseConfiguration().BulkWrite(ctx, wModel); err != nil {
			logger.Error("Error dao.WarehouseConfiguration().BulkWrite : ", logger.LogData{
				"error": err.Error(),
			})
		}
	}

	if len(listSupplier) == 0 {
		return
	}
	for _, supplier := range listSupplier {
		warehouses = dao.Warehouse().FindByCondition(ctx, bson.M{
			"supplier": mongodb.ConvertStringToObjectID(supplier),
		})
		listWarehouseID := make([]primitive.ObjectID, 0)
		for _, wh := range warehouses {
			listWarehouseID = append(listWarehouseID, wh.ID)
		}
		configurations := dao.WarehouseConfiguration().FindByCondition(ctx, bson.M{
			"warehouse": bson.M{
				"$in": listWarehouseID,
			},
		})
		if len(configurations) == 0 {
			listUpdateSupplier = append(listUpdateSupplier, model.SupplierIsClosed{
				Supplier: supplier,
				IsClosed: true,
			})
			continue
		}
		var isClosed = true
		for _, cf := range configurations {
			if cf.Food.IsClosed {
				isClosed = false
				break
			}
		}
		listUpdateSupplier = append(listUpdateSupplier, model.SupplierIsClosed{
			Supplier: supplier,
			IsClosed: isClosed,
		})
	}
	go client.GetWarehouse().UpdateIsClosedSupplier(model.UpdateSupplierIsClosedRequest{Suppliers: listUpdateSupplier})
	return
}

// RunJobUpdateHolidayWarehouses ...
func (s Schedule) RunJobUpdateHolidayWarehouses() {
	var (
		ctx          = context.Background()
		warehouseSvc = warehouseImplement{}
		cond         = bson.M{
			"status": bson.M{"$in": []string{
				constant.WarehouseStatusHoliday,
				constant.StatusActive,
			}},
		}
	)

	fmt.Println("Start job update holiday warehouse ...!")

	warehouses := dao.Warehouse().FindByCondition(ctx, cond)
	warehouseSvc.UpdateWarehousesStatusByIDs(warehouses)

	fmt.Println("End job update holiday warehouse ...")
	return
}

// RunJobUpdateHolidayStatusForSupplier ...
func (s Schedule) RunJobUpdateHolidayStatusForSupplier() {
	var supplierHolidaySvc = supplierHolidayImplement{}

	fmt.Println("Start job update holiday status for supplier ...!")
	supplierHolidaySvc.UpdateHolidayStatusForSupplier()
	fmt.Println("End job job update holiday status for supplier ...")
}

func (s Schedule) UpdatePaymentMethodBankTransferWarehouse() {
	var (
		ctx  = context.Background()
		cond = bson.M{
			"changPaymentMethod": "13_01_2023_false",
		}
	)
	d := dao.WarehouseConfiguration()
	payload := bson.M{
		"$set": bson.M{
			"changPaymentMethod":               "27_01_2023_true",
			"order.paymentMethod.bankTransfer": true,
		},
	}
	d.UpdateMany(ctx, cond, payload)
}

func CheckIsClosed(m mgwarehouse.ConfigFood) bool {
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
