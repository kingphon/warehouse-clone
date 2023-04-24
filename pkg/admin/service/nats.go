package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"git.selly.red/Selly-Modules/redisdb"
	"github.com/logrusorgru/aurora"
	"go.mongodb.org/mongo-driver/bson"

	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"

	"git.selly.red/Selly-Modules/natsio/model"
	"github.com/friendsofgo/errors"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"

	"github.com/panjf2000/ants/v2"
)

// Nats ...
type Nats struct{}

// CreateOutboundRequest ...
func (s Nats) CreateOutboundRequest(req *model.OutboundRequestPayload) (*model.OutboundRequestResponse, error) {
	return OutboundRequest().Create(req)
}

// UpdateLogisticInfo ...
func (s Nats) UpdateLogisticInfo(req *model.UpdateOutboundRequestLogisticInfoPayload) error {
	return OutboundRequest().UpdateLogisticInfo(req)
}

// CancelOutboundRequest ...
func (s Nats) CancelOutboundRequest(req *model.CancelOutboundRequest) error {
	return OutboundRequest().Cancel(req)
}

// DistinctWarehouseWithField ...
func (s Nats) DistinctWarehouseWithField(ctx context.Context, req *model.DistinctWithField) ([]interface{}, error) {
	var (
		res []interface{}
		err error
	)
	if req.Filed == "" {
		return nil, errors.New("field is required!")
	}
	res, err = dao.Warehouse().DistinctWithField(ctx, req.Conditions, req.Filed)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CountWarehouseWithCondition ...
func (s Nats) CountWarehouseWithCondition(ctx context.Context, req *model.FindWithCondition) (int64, error) {
	var (
		total int64
	)
	total = dao.Warehouse().CountByCondition(ctx, req.Conditions)
	return total, nil
}

// WarehouseItem ...
type WarehouseItem struct {
	Index int
	Data  mgwarehouse.Warehouse
}

// WarehouseItemResponse ...
type WarehouseItemResponse struct {
	Index int
	Data  model.WarehouseNatsResponse
}

// GetWarehouseWithCondition ...
func (s Nats) GetWarehouseWithCondition(ctx context.Context, req *model.FindWithCondition, id string) ([]model.WarehouseNatsResponse, error) {
	var (
		res      []model.WarehouseNatsResponse
		err      error
		isGetAll = false
		k        = constant.RedisKeyWarehouseAll
	)
	if req != nil && len(req.Opts) == 0 {
		m, ok := req.Conditions.(bson.D)
		if ok && len(m) == 0 {
			isGetAll = true
		}
	}
	if isGetAll {
		ok := redisdb.GetJSON(ctx, k, &res)
		if ok {
			return res, nil
		}
	}

	start := time.Now().UTC()
	warehouses := dao.Warehouse().FindByCondition(ctx, req.Conditions, req.Opts...)
	if len(warehouses) == 0 {
		return res, err
	}

	res = make([]model.WarehouseNatsResponse, len(warehouses))

	if len(warehouses) < 100 {
		for i, warehouse := range warehouses {
			res[i], err = s.getWareHouseResponse(ctx, warehouse)
			if err != nil {
				log.Println("getWareHouseResponse", err)
			}
		}
		return res, nil
	}

	var (
		wg   sync.WaitGroup
		ch   = make(chan WarehouseItemResponse)
		done = make(chan bool)
	)

	go func() {
		for data := range ch {
			res[data.Index] = data.Data
		}
		done <- true
	}()

	p, _ := ants.NewPoolWithFunc(50, func(i interface{}) {
		defer wg.Done()
		data := i.(WarehouseItem)
		wr, _ := s.getWareHouseResponse(ctx, data.Data)
		ch <- WarehouseItemResponse{
			Index: data.Index,
			Data:  wr,
		}
	})
	defer p.Release()
	for index, w := range warehouses {
		wg.Add(1)
		p.Invoke(WarehouseItem{
			Index: index,
			Data:  w,
		})
	}
	wg.Wait()

	close(ch)
	<-done
	dur := time.Since(start).Seconds()
	if dur > 1 {
		fmt.Println(aurora.Green(fmt.Sprintf("Nats.GetWarehouseWithCondition find db: Time find.all %f", dur)))
		fmt.Printf("condition: %v\n", req.Conditions)
		fmt.Println("-------")
	}

	if isGetAll {
		redisdb.SetTTL(ctx, k, res, time.Second*30)
	}
	return res, nil
}

// GetOneWarehouse ...
func (s Nats) GetOneWarehouse(ctx context.Context, req *model.FindOneCondition) (model.WarehouseNatsResponse, error) {
	var (
		res model.WarehouseNatsResponse
		err error
	)
	start := time.Now().UTC()
	warehouse := dao.Warehouse().FindOneByCondition(ctx, req.Conditions)
	if warehouse.ID.IsZero() {
		return model.WarehouseNatsResponse{}, errors.New("warehouse not found")
	}
	res, err = s.getWareHouseResponse(ctx, warehouse)
	if err != nil {
		return model.WarehouseNatsResponse{}, err
	}
	dur := time.Since(start).Seconds()
	if dur > 1 {
		fmt.Println(aurora.Green(fmt.Sprintf("Nats.GetOneWarehouse find db: Time find.all %f", dur)))
		fmt.Printf("condition: %v\n", req.Conditions)
		fmt.Println("-------")
	}
	return res, nil
}

func (s Nats) getWareHouseResponse(ctx context.Context, warehouse mgwarehouse.Warehouse) (model.WarehouseNatsResponse, error) {
	configurations := dao.WarehouseConfiguration().FindByWarehouseID(ctx, warehouse.ID)
	if configurations.ID.IsZero() {
		return model.WarehouseNatsResponse{}, errors.New("w_configuration not found")
	}
	return responsemodel.GetWarehouseNatsResponse(warehouse, configurations, ""), nil
}
