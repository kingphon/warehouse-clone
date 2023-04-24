package service

import (
	"context"
	"sync"

	"git.selly.red/Selly-Modules/mongodb"
	"git.selly.red/Selly-Server/warehouse/external/constant"
	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/mgquery"
	"git.selly.red/Selly-Server/warehouse/external/utils/parray"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
	"github.com/friendsofgo/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//
// PUBLIC METHODS
//

// All ...
func (s supplierHolidayImplement) All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseSupplierHolidayAll) {

	// 1. Assign condition
	cond := bson.D{}
	q.Warehouse.AssignKeyword(&cond)
	q.Warehouse.AssignSupplier(&cond)
	q.Warehouse.AssignWarehouse(&cond)
	q.Warehouse.AssignFromToAt(&cond)
	q.AssignSource(&cond)
	q.AssignStatuses(&cond)

	var (
		d  = dao.SupplierHoliday()
		wg = sync.WaitGroup{}
	)

	wg.Add(2)

	// 2. Get data
	go func() {
		defer wg.Done()

		result.List = make([]responsemodel.ResponseSupplierHolidayBrief, 0)

		findOpts := q.GetFindOptionsWithPage()

		docs := d.FindByCondition(ctx, cond, findOpts)

		suppliers, err := s.getSuppliersBySupplierHolidays(ctx, docs)
		if err != nil {
			return
		}

		warehouses := s.getWarehousesBySupplierHolidays(ctx, docs)

		staffs, err := s.getStaffBySupplierHolidays(ctx, docs)
		if err != nil {
			return
		}

		result.List = s.getListSupplierHolidayBrief(ctx, docs, suppliers, warehouses, staffs)
	}()

	// 3. Assign total
	go func() {
		defer wg.Done()
		result.Total = d.CountByCondition(ctx, cond)
	}()

	wg.Wait()

	// 3. Assign limit
	result.Limit = q.Limit

	return
}

// Detail ...
func (s supplierHolidayImplement) Detail(ctx context.Context, id primitive.ObjectID) (result *responsemodel.ResponseSupplierHolidayDetail, err error) {
	var (
		d    = dao.SupplierHoliday()
		cond = bson.M{"_id": id}
	)

	supplierHoliday := d.FindOneByCondition(ctx, cond)
	if supplierHoliday.ID.IsZero() {
		err = errors.New(errorcode.SupplierHolidayNotFound)
		return
	}

	var supplierHolidayList = []mgwarehouse.SupplierHoliday{supplierHoliday}

	suppliers, err := s.getSuppliersBySupplierHolidays(ctx, supplierHolidayList)
	if err != nil {
		return
	}

	warehouses := s.getWarehousesBySupplierHolidays(ctx, supplierHolidayList)

	staffs, err := s.getStaffBySupplierHolidays(ctx, supplierHolidayList)
	if err != nil {
		return
	}

	result = s.getSupplierHolidayDetailInfo(ctx, supplierHoliday, suppliers, warehouses, staffs)
	return
}

// FindByCondition ...
func (s supplierHolidayImplement) FindByCondition(ctx context.Context, cond interface{}) []mgwarehouse.SupplierHoliday {
	return dao.SupplierHoliday().FindByCondition(ctx, cond)
}

//
// PRIVATE METHODS
//

// getStaffBySupplierHolidays ...
func (supplierHolidayImplement) getStaffBySupplierHolidays(ctx context.Context, supplierHoliday []mgwarehouse.SupplierHoliday) (result []responsemodel.ResponseStaffShort, err error) {
	var staffIds = make([]string, 0)
	for _, s := range supplierHoliday {
		staffIds = append(staffIds, s.CreatedBy.Hex())
	}

	uniqIDs := parray.UniqueArrayStrings(staffIds)

	staffs, err := GetStaffByIDs(uniqIDs)
	if err != nil {
		return
	}

	for _, s := range staffs.Staffs {
		result = append(result, responsemodel.ResponseStaffShort{
			ID:   s.ID,
			Name: s.Name,
		})
	}

	return
}

// getSuppliersBySupplierHolidays ...
func (supplierHolidayImplement) getSuppliersBySupplierHolidays(ctx context.Context, docs []mgwarehouse.SupplierHoliday) (result []responsemodel.ResponseSupplierShort, err error) {
	supplierIDs := make([]primitive.ObjectID, 0)
	for _, d := range docs {
		supplierIDs = append(supplierIDs, d.Supplier)
	}

	suppliers, err := GetSupplierByIDs(mongodb.UniqObjectIds(supplierIDs))
	if err != nil {
		return
	}

	for _, s := range suppliers {
		result = append(result, responsemodel.ResponseSupplierShort{
			ID:   s.ID,
			Name: s.Name,
		})
	}

	return
}

// getWarehousesBySupplierHolidays ...
func (s supplierHolidayImplement) getWarehousesBySupplierHolidays(ctx context.Context, docs []mgwarehouse.SupplierHoliday) (result []responsemodel.ResponseWarehouseShort) {
	result = make([]responsemodel.ResponseWarehouseShort, 0)

	if len(docs) == 0 {
		return
	}

	warehouseIDs := make([]primitive.ObjectID, 0)
	for _, d := range docs {
		if len(d.Warehouses) > 0 {
			for _, id := range d.Warehouses {
				warehouseIDs = append(warehouseIDs, id)
			}
		}
	}

	var (
		warehouseSvc = warehouseImplement{CurrentStaff: s.CurrentStaff}
		whUniqIDs    = mongodb.UniqObjectIds(warehouseIDs)
	)
	listWarehouse := warehouseSvc.FindByCondition(ctx, bson.M{"_id": bson.M{"$in": whUniqIDs}})

	for _, w := range listWarehouse {
		result = append(result, responsemodel.ResponseWarehouseShort{
			ID:   w.ID.Hex(),
			Name: w.Name,
		})
	}

	return
}

// getSupplierHolidays ...
func (s supplierHolidayImplement) getListSupplierHolidayBrief(ctx context.Context, docs []mgwarehouse.SupplierHoliday, suppliers []responsemodel.ResponseSupplierShort, warehouses []responsemodel.ResponseWarehouseShort, staffs []responsemodel.ResponseStaffShort) (result []responsemodel.ResponseSupplierHolidayBrief) {
	var (
		wg    = sync.WaitGroup{}
		total = len(docs)
	)

	result = make([]responsemodel.ResponseSupplierHolidayBrief, total)

	wg.Add(total)
	for i, doc := range docs {
		go func(i int, doc mgwarehouse.SupplierHoliday) {
			defer wg.Done()
			result[i] = s.getSupplierHolidayBriefInfo(ctx, doc, suppliers, warehouses, staffs)
		}(i, doc)
	}

	wg.Wait()
	return
}

// getSupplierSupplierHolidayBrief ...
func (s supplierHolidayImplement) getSupplierHolidayBriefInfo(ctx context.Context, doc mgwarehouse.SupplierHoliday, suppliers []responsemodel.ResponseSupplierShort, warehouses []responsemodel.ResponseWarehouseShort, staffs []responsemodel.ResponseStaffShort) responsemodel.ResponseSupplierHolidayBrief {
	var (
		supplier      responsemodel.ResponseSupplierShort
		warehouseList []responsemodel.ResponseWarehouseShort
		staff         responsemodel.ResponseStaffShort
		wg            = sync.WaitGroup{}
	)

	// Supplier
	wg.Add(1)
	go func() {
		defer wg.Done()

		if found := parray.Find(suppliers, func(item responsemodel.ResponseSupplierShort) bool {
			return item.ID == doc.Supplier.Hex()
		}); found != nil {
			supplier = found.(responsemodel.ResponseSupplierShort)
		}
	}()

	// Warehouses
	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, id := range doc.Warehouses {
			if found := parray.Find(warehouses, func(item responsemodel.ResponseWarehouseShort) bool {
				return item.ID == id.Hex()
			}); found != nil {
				warehouseList = append(warehouseList, found.(responsemodel.ResponseWarehouseShort))
			}
		}

	}()

	// CreatedBy
	wg.Add(1)
	go func() {
		defer wg.Done()

		if found := parray.Find(staffs, func(item responsemodel.ResponseStaffShort) bool {
			return item.ID == doc.CreatedBy.Hex()
		}); found != nil {
			staff = found.(responsemodel.ResponseStaffShort)
		}
	}()

	wg.Wait()

	return s.brief(ctx, doc, supplier, warehouseList, staff)
}

// brief ...
func (s supplierHolidayImplement) brief(ctx context.Context, doc mgwarehouse.SupplierHoliday, supplier responsemodel.ResponseSupplierShort, warehouses []responsemodel.ResponseWarehouseShort, staff responsemodel.ResponseStaffShort) responsemodel.ResponseSupplierHolidayBrief {
	result := responsemodel.ResponseSupplierHolidayBrief{
		ID:         doc.ID.Hex(),
		Warehouses: make([]responsemodel.ResponseWarehouseShort, 0),
		Title:      doc.Title,
		From:       ptime.TimeResponseInit(doc.From),
		To:         ptime.TimeResponseInit(doc.To),
		Reason:     doc.Reason,
		Source:     doc.Source,
		Status:     doc.Status,
		Supplier: responsemodel.ResponseSupplierShort{
			ID:   supplier.ID,
			Name: supplier.Name,
		},
		IsApplyAll: doc.IsApplyAll,
	}

	if len(warehouses) > 0 {
		for _, w := range warehouses {
			result.Warehouses = append(result.Warehouses, responsemodel.ResponseWarehouseShort{
				ID:   w.ID,
				Name: w.Name,
			})
		}
	} else {
		if doc.IsApplyAll {
			var warehouseSvc = warehouseImplement{CurrentStaff: s.CurrentStaff}
			getWarehousesBySuppliers := warehouseSvc.FindByCondition(ctx, bson.M{"supplier": doc.Supplier})
			for _, w := range getWarehousesBySuppliers {
				result.Warehouses = append(result.Warehouses, responsemodel.ResponseWarehouseShort{
					ID:   w.ID.Hex(),
					Name: w.Name,
				})
			}
		}
	}

	if doc.Source == constant.WarehouseSupplierHolidaySourceAdmin {
		result.CreatedBy = responsemodel.ResponseInfo{
			ID:   staff.ID,
			Name: staff.Name,
		}
	} else {
		result.CreatedBy = responsemodel.ResponseInfo{
			ID:   supplier.ID,
			Name: supplier.Name,
		}
	}

	return result
}

// getSupplierHolidayDetailInfo ...
func (s supplierHolidayImplement) getSupplierHolidayDetailInfo(ctx context.Context, doc mgwarehouse.SupplierHoliday, suppliers []responsemodel.ResponseSupplierShort, warehouses []responsemodel.ResponseWarehouseShort, staffs []responsemodel.ResponseStaffShort) *responsemodel.ResponseSupplierHolidayDetail {
	var (
		supplier      responsemodel.ResponseSupplierShort
		warehouseList []responsemodel.ResponseWarehouseShort
		staff         responsemodel.ResponseStaffShort
		wg            = sync.WaitGroup{}
	)

	// Supplier
	wg.Add(1)
	go func() {
		defer wg.Done()

		if found := parray.Find(suppliers, func(item responsemodel.ResponseSupplierShort) bool {
			return item.ID == doc.Supplier.Hex()
		}); found != nil {
			supplier = found.(responsemodel.ResponseSupplierShort)
		}
	}()

	// Warehouses
	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, id := range doc.Warehouses {
			if found := parray.Find(warehouses, func(item responsemodel.ResponseWarehouseShort) bool {
				return item.ID == id.Hex()
			}); found != nil {
				warehouseList = append(warehouseList, found.(responsemodel.ResponseWarehouseShort))
			}
		}

		// CreatedBy
		wg.Add(1)
		go func() {
			defer wg.Done()

			if found := parray.Find(staffs, func(item responsemodel.ResponseStaffShort) bool {
				return item.ID == doc.CreatedBy.Hex()
			}); found != nil {
				staff = found.(responsemodel.ResponseStaffShort)
			}
		}()

	}()

	wg.Wait()

	return s.detail(ctx, doc, supplier, warehouseList, staff)
}

// detail ...
func (s supplierHolidayImplement) detail(ctx context.Context, doc mgwarehouse.SupplierHoliday, supplier responsemodel.ResponseSupplierShort, warehouses []responsemodel.ResponseWarehouseShort, staff responsemodel.ResponseStaffShort) *responsemodel.ResponseSupplierHolidayDetail {
	result := &responsemodel.ResponseSupplierHolidayDetail{
		ID:    doc.ID.Hex(),
		Title: doc.Title,
		Supplier: responsemodel.ResponseSupplierShort{
			ID:   supplier.ID,
			Name: supplier.Name,
		},
		From:       ptime.TimeResponseInit(doc.From),
		To:         ptime.TimeResponseInit(doc.To),
		Reason:     doc.Reason,
		Source:     doc.Source,
		Status:     doc.Status,
		Warehouses: make([]responsemodel.ResponseWarehouseShort, 0),
		IsApplyAll: doc.IsApplyAll,
	}

	if len(warehouses) == 0 && doc.IsApplyAll {
		var warehouseSvc = warehouseImplement{CurrentStaff: s.CurrentStaff}
		getWarehousesBySuppliers := warehouseSvc.FindByCondition(ctx, bson.M{"supplier": doc.Supplier})
		for _, w := range getWarehousesBySuppliers {
			result.Warehouses = append(result.Warehouses, responsemodel.ResponseWarehouseShort{
				ID:   w.ID.Hex(),
				Name: w.Name,
			})
		}
	} else {
		for _, w := range warehouses {
			result.Warehouses = append(result.Warehouses, responsemodel.ResponseWarehouseShort{
				ID:   w.ID,
				Name: w.Name,
			})
		}
	}

	if doc.Source == constant.WarehouseSupplierHolidaySourceAdmin {
		result.CreatedBy = responsemodel.ResponseStaffShort{
			ID:   staff.ID,
			Name: staff.Name,
		}
	} else {
		result.CreatedBy = responsemodel.ResponseStaffShort{
			ID:   supplier.ID,
			Name: supplier.Name,
		}
	}

	return result
}
