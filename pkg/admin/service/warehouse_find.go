package service

import (
	"context"
	"errors"
	"git.selly.red/Selly-Server/warehouse/external/constant"
	"sync"

	"git.selly.red/Selly-Modules/natsio/client"
	"git.selly.red/Selly-Modules/natsio/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/mgquery"
	"git.selly.red/Selly-Server/warehouse/external/utils/parray"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/errorcode"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
)

//
// PUBLIC METHODS
//

// All ...
func (s warehouseImplement) All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseWarehouseAll) {
	var (
		d  = dao.Warehouse()
		wg sync.WaitGroup
	)

	// Condition
	cond := bson.M{}
	q.AssignKeyword(cond)
	q.AssignWarehouseStatus(cond)
	q.AssignWarehouseSupplier(cond)
	q.AssignBusinessType(cond)
	// Start wait group
	wg.Add(2)

	// Find
	go func() {
		defer wg.Done()

		// Prepare data
		result.List = make([]responsemodel.WarehouseBrief, 0)

		// Find options
		findOpts := q.GetFindOptionsWithPage()
		findOpts.SetProjection(bson.M{
			"_id":          1,
			"businessType": 1,
			"name":         1,
			"status":       1,
			"slug":         1,
			"supplier":     1,
			"contact":      1,
			"location":     1,
			"createdAt":    1,
			"updatedAt":    1,
		})

		// Find
		docs := d.FindByCondition(ctx, cond, findOpts)

		// Get list supplierID
		listSupplierID := make([]primitive.ObjectID, len(docs))
		for i, doc := range docs {
			listSupplierID[i] = doc.Supplier
		}

		// get suppliers...
		suppliers, err := s.getSupplierByIDs(ctx, listSupplierID)
		if err != nil {
			err = errors.New(errorcode.WarehouseInvalidSupplier)
			return
		}

		var (
			g = errgroup.Group{}

			listProvince *model.LocationProvinceResponse
			listDistrict *model.LocationDistrictResponse
			listWard     *model.LocationWardResponse
		)

		// Get provinces ...
		g.Go(func() error {
			provinces, err := s.getProvincesByCodes(ctx, docs)
			if err != nil {
				return errors.New(errorcode.WarehouseInvalidProvince)
			}

			listProvince = provinces
			return nil
		})

		// Get Districts ...
		g.Go(func() error {
			districts, err := s.getDistrictsByCodes(ctx, docs)
			if err != nil {
				return errors.New(errorcode.WarehouseInvalidDistrict)
			}

			listDistrict = districts
			return nil
		})

		// Get Ward ...
		g.Go(func() error {
			wards, err := s.getWardsByCodes(ctx, docs)
			if err != nil {
				return errors.New(errorcode.WarehouseInvalidWard)
			}

			listWard = wards
			return nil
		})

		if err := g.Wait(); err != nil {
			return
		}

		for _, doc := range docs {
			var (
				supplier = &model.ResponseSupplierInfo{}
				province = model.LocationProvince{}
				district = model.LocationDistrict{}
				ward     = model.LocationWard{}

				wg = sync.WaitGroup{}
			)

			wg.Add(4)

			// Supplier ...
			go func() {
				if found := parray.Find(suppliers, func(item *model.ResponseSupplierInfo) bool {
					return item.ID == doc.Supplier.Hex()
				}); found != nil {
					supplier = found.(*model.ResponseSupplierInfo)
				}

				wg.Done()
			}()

			// Province
			go func() {
				if foundP := parray.Find(listProvince.Provinces, func(item model.LocationProvince) bool {
					return item.Code == doc.Location.Province
				}); foundP != nil {
					province = foundP.(model.LocationProvince)
				}

				wg.Done()
			}()

			// District
			go func() {
				if foundD := parray.Find(listDistrict.Districts, func(item model.LocationDistrict) bool {
					return item.Code == doc.Location.District
				}); foundD != nil {
					district = foundD.(model.LocationDistrict)
				}

				wg.Done()
			}()

			// Ward
			go func() {
				if foundW := parray.Find(listWard.Wards, func(item model.LocationWard) bool {
					return item.Code == doc.Location.Ward
				}); foundW != nil {
					ward = foundW.(model.LocationWard)
				}

				wg.Done()
			}()

			wg.Wait()
			result.List = append(result.List, s.brief(ctx, doc, supplier, province, district, ward))
		}
	}()

	// Count
	go func() {
		defer wg.Done()
		result.Total = d.CountByCondition(ctx, cond)
	}()

	// Wait
	wg.Wait()

	// Assign limit
	result.Limit = q.Limit

	return
}

// Detail ...
func (s warehouseImplement) Detail(ctx context.Context, id primitive.ObjectID) (*responsemodel.ResponseWarehouseDetail, error) {
	var (
		d    = dao.Warehouse()
		cond = bson.M{"_id": id}
	)

	warehouse := d.FindOneByCondition(ctx, cond)
	if warehouse.ID.IsZero() {
		return nil, errors.New(errorcode.WarehouseNotFound)
	}

	var (
		wg       = sync.WaitGroup{}
		location = &model.ResponseLocationAddress{}
		supplier = &model.ResponseSupplierInfo{
			ID:   "",
			Name: "",
		}
	)

	wg.Add(2)

	go func() {
		// Get location by code
		location, _ = WarehouseLocation{}.GetLocationByCode(ctx, warehouse)
		wg.Done()
	}()

	go func() {
		// Check supplier invalid and get info
		listSupplierID := []primitive.ObjectID{warehouse.Supplier}
		suppliers, err := s.getSupplierByIDs(ctx, listSupplierID)
		if err == nil && len(suppliers) > 0 {
			supplier = &model.ResponseSupplierInfo{
				ID:   suppliers[0].ID,
				Name: suppliers[0].Name,
			}
		}
		wg.Done()
	}()

	wg.Wait()

	return s.detail(ctx, warehouse, location, supplier), nil
}

// FindByCondition ...
func (s warehouseImplement) FindByCondition(ctx context.Context, cond interface{}) (result []mgwarehouse.Warehouse) {
	return dao.Warehouse().FindByCondition(ctx, cond)
}

// FindOneByCondition ...
func (s warehouseImplement) FindOneByCondition(ctx context.Context, cond interface{}) (result mgwarehouse.Warehouse) {
	return dao.Warehouse().FindOneByCondition(ctx, cond)
}

//
// PRIVATE METHODS
//

// brief ...
func (s warehouseImplement) brief(ctx context.Context, doc mgwarehouse.Warehouse, supplier *model.ResponseSupplierInfo, province model.LocationProvince, district model.LocationDistrict, ward model.LocationWard) responsemodel.WarehouseBrief {
	latLng := s.ConvertCoordinatesToLatLng(ctx, doc.Location)

	return responsemodel.WarehouseBrief{
		ID:           doc.ID.Hex(),
		BusinessType: doc.BusinessType,
		Name:         doc.Name,
		Status:       doc.Status,
		Slug:         doc.Slug,
		Supplier: responsemodel.WarehouseSupplier{
			ID:   doc.Supplier.Hex(),
			Name: supplier.Name,
		},
		Location: responsemodel.ResponseWarehouseLocation{
			Province: responsemodel.ResponseWarehouseProvince{
				ID:   province.ID,
				Name: province.Name,
				Code: province.Code,
			},
			District: responsemodel.ResponseWarehouseDistrict{
				ID:   district.ID,
				Name: district.Name,
				Code: district.Code,
			},
			Ward: responsemodel.ResponseWarehouseWard{
				ID:   ward.ID,
				Name: ward.Name,
				Code: ward.Code,
			},
			Address:             doc.Location.Address,
			LocationCoordinates: latLng,
		},
		Contact: responsemodel.ResponseWarehouseContact{
			Name:    doc.Contact.Name,
			Phone:   doc.Contact.Phone,
			Address: doc.Contact.Address,
			Email:   doc.Contact.Email,
		},
		CreatedAt: ptime.TimeResponseInit(doc.CreatedAt),
		UpdatedAt: ptime.TimeResponseInit(doc.UpdatedAt),
	}
}

// detail ...
func (s warehouseImplement) detail(ctx context.Context, d mgwarehouse.Warehouse, location *model.ResponseLocationAddress, supplier *model.ResponseSupplierInfo) *responsemodel.ResponseWarehouseDetail {
	latLng := s.ConvertCoordinatesToLatLng(ctx, d.Location)

	result := &responsemodel.ResponseWarehouseDetail{
		ID:           d.ID.Hex(),
		BusinessType: d.BusinessType,
		Name:         d.Name,
		Supplier: responsemodel.WarehouseSupplier{
			ID:   d.Supplier.Hex(),
			Name: supplier.Name,
		},
		Location: responsemodel.ResponseWarehouseLocation{
			Province: responsemodel.ResponseWarehouseProvince{
				ID:   location.Province.ID,
				Name: location.Province.Name,
				Code: location.Province.Code,
			},
			District: responsemodel.ResponseWarehouseDistrict{
				ID:   location.District.ID,
				Name: location.District.Name,
				Code: location.District.Code,
			},
			Ward: responsemodel.ResponseWarehouseWard{
				ID:   location.Ward.ID,
				Name: location.Ward.Name,
				Code: location.Ward.Code,
			},
			Address:             d.Location.Address,
			LocationCoordinates: latLng,
		},
		Contact: responsemodel.ResponseWarehouseContact{
			Name:    d.Contact.Name,
			Phone:   d.Contact.Phone,
			Address: d.Contact.Address,
			Email:   d.Contact.Email,
		},
		CreatedAt:             ptime.TimeResponseInit(d.CreatedAt),
		UpdatedAt:             ptime.TimeResponseInit(d.UpdatedAt),
		ReasonPendingInactive: d.ReasonPendingInactive,
		Status:                d.Status,
		StatusBeforeHoliday:   d.StatusBeforeHoliday,
	}
	if result.Status == constant.WarehouseStatusHoliday {
		var (
			supplierHolidaySvc = supplierHolidayImplement{s.CurrentStaff}
			cond               = bson.M{
				"$or": []bson.M{
					bson.M{"isApplyAll": true},
					bson.M{"isApplyAll": false, "warehouses": d.ID},
				},
				"supplier": d.Supplier,
				"status":   constant.StatusActive,
				"from":     bson.M{"$lte": ptime.Now()},
				"to":       bson.M{"$gte": ptime.Now()},
			}
		)
		supplierHolidaysByWarehouse := supplierHolidaySvc.FindByCondition(ctx, cond)
		if len(supplierHolidaysByWarehouse) > 0 {
			result.SupplierHolidayFrom = ptime.TimeResponseInit(supplierHolidaysByWarehouse[0].From)
			result.SupplierHolidayTo = ptime.TimeResponseInit(supplierHolidaysByWarehouse[0].To)
		}
	}

	return result

}

func (warehouseImplement) getLocationByCode(ctx context.Context, w mgwarehouse.Warehouse) (*model.ResponseLocationAddress, error) {
	body := model.LocationRequestPayload{
		Province: w.Location.Province,
		District: w.Location.District,
		Ward:     w.Location.Ward,
	}
	return client.GetLocation().GetLocationByCode(body)
}

// getSupplierByID ...
func (warehouseImplement) getSupplierByIDs(ctx context.Context, listID []primitive.ObjectID) ([]*model.ResponseSupplierInfo, error) {
	body := model.GetSupplierRequest{ListID: listID}
	return client.GetSupplier().GetListSupplierInfo(body)
}

// getProvincesByCodes ..
func (warehouseImplement) getProvincesByCodes(ctx context.Context, docs []mgwarehouse.Warehouse) (*model.LocationProvinceResponse, error) {
	// Get list provinceID
	listProvinceID := make([]int, len(docs))
	for i, doc := range docs {
		listProvinceID[i] = doc.Location.Province
	}

	body := model.ProvinceRequestPayload{Codes: listProvinceID}
	return client.GetLocation().GetProvincesByCodes(body)
}

// getDistrictsByCodes ..
func (warehouseImplement) getDistrictsByCodes(ctx context.Context, docs []mgwarehouse.Warehouse) (*model.LocationDistrictResponse, error) {
	// Get list districtID
	listDistrictID := make([]int, len(docs))
	for i, doc := range docs {
		listDistrictID[i] = doc.Location.District
	}

	body := model.DistrictRequestPayload{Codes: listDistrictID}
	return client.GetLocation().GetDistrictsByCodes(body)
}

// getWardByCodes ...
func (warehouseImplement) getWardsByCodes(ctx context.Context, docs []mgwarehouse.Warehouse) (*model.LocationWardResponse, error) {
	// get list wardID
	listWardID := make([]int, len(docs))
	for i, doc := range docs {
		listWardID[i] = doc.Location.Ward
	}

	body := model.WardRequestPayload{Codes: listWardID}
	return client.GetLocation().GetWardsByCodes(body)
}

func (warehouseImplement) ConvertCoordinatesToLatLng(ctx context.Context, l mgwarehouse.Location) responsemodel.ResponseLatLng {
	var latLng = responsemodel.ResponseLatLng{}
	if len(l.LocationCoordinates.Coordinates) < 1 {
		return responsemodel.ResponseLatLng{
			Latitude:  0,
			Longitude: 0,
		}
	}
	latLng = responsemodel.ResponseLatLng{
		Latitude:  l.LocationCoordinates.Coordinates[1],
		Longitude: l.LocationCoordinates.Coordinates[0],
	}
	return latLng
}
