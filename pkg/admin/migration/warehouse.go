package migration

import (
	"context"
	"fmt"
	"time"

	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/warehouse/external/constant"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"

	"git.selly.red/Selly-Modules/mongodb"
	natsclient "git.selly.red/Selly-Modules/natsio/client"
	natsmodel "git.selly.red/Selly-Modules/natsio/model"
	"github.com/labstack/echo/v4"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/database"
)

// Migrate ...
func Migrate(g *echo.Group) {
	g.POST("/warehouse", migrateWarehouses)
	g.GET("/location-refactor", migrateLocationRefactor)

	g.POST("/warehouse-business-type", migrateWarehousesBusinessType)

	g.GET("/migrate-update-status-warehouse-pending-active", MigrateUpdateStatusWarehousePendingInactive)
}

func migrateWarehousesBusinessType(c echo.Context) error {
	col := database.WarehouseCol()
	_, _ = col.UpdateMany(context.Background(), bson.M{}, bson.M{
		"$set": bson.M{
			"businessType": "other",
		},
	})

	return c.JSON(200, "ok")
}

func migrateWarehouses(c echo.Context) error {
	db := database.WarehouseCol().Database()
	oldDB := db.Client().Database("unibag")
	// provinceCol := oldDB.Collection("cities")
	// districtCol := oldDB.Collection("districts")
	wardCol := oldDB.Collection("wards")
	inventoryCol := oldDB.Collection("inventories")
	whCol := db.Collection("warehouses")
	whCfgCol := db.Collection("warehouse-configurations")
	ctx := context.Background()
	cursor, err := inventoryCol.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var docs []inventoryRaw
	if err = cursor.All(ctx, &docs); err != nil {
		return err
	}
	statusMap := map[bool]string{
		true:  "active",
		false: "inactive",
	}

	type province struct {
		Code int    `bson:"code"`
		Slug string `bson:"slug"`
	}
	type district struct {
		Code         int    `bson:"code"`
		Slug         string `bson:"slug"`
		ProvinceCode int    `bson:"cityId"`
		ProvinceSlug string `bson:"city"`
	}
	type ward struct {
		Code         int    `bson:"code"`
		Slug         string `bson:"slug"`
		ProvinceCode int    `bson:"cityId"`
		ProvinceSlug string `bson:"city"`
		DistrictCode int    `bson:"districtId"`
		DistrictSlug string `bson:"district"`
	}
	var (
		// provinces []province
		// districts []district
		wards []ward
	)
	// cursor, err = provinceCol.Find(ctx, bson.M{})
	// if err != nil {
	// 	return err
	// }
	// cursor.All(ctx, &provinces)
	// cursor, err = districtCol.Find(ctx, bson.M{})
	// if err != nil {
	// 	return err
	// }
	// cursor.All(ctx, &districts)
	cursor, err = wardCol.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	cursor.All(ctx, &wards)
	// fmt.Println("Provinces: ", len(provinces))
	// fmt.Println("Districts: ", len(districts))
	fmt.Println("Wards: ", len(wards))

	getWardBySlug := func(provinceSlug, districtSlug, wardSlug string) (w ward) {
		found := funk.Find(wards, func(item ward) bool {
			return item.Slug == wardSlug && item.ProvinceSlug == provinceSlug && item.DistrictSlug == districtSlug
		})
		if found != nil {
			w = found.(ward)
		}
		return w
	}

	var warehouses, whCfgs []interface{}
	for _, doc := range docs {
		loc := doc.Location
		w := getWardBySlug(loc.Province, loc.District, loc.Ward)
		wh := mgwarehouse.Warehouse{
			ID:           doc.ID,
			Name:         doc.Name,
			SearchString: doc.SearchString,
			Slug:         doc.Slug,
			Status:       statusMap[doc.Active],
			Supplier:     doc.Supplier.ID,
			Contact:      doc.Contact,
			Location: mgwarehouse.Location{
				Province:            w.ProvinceCode,
				District:            w.DistrictCode,
				Ward:                w.Code,
				Address:             doc.Location.Address,
				LocationCoordinates: doc.Location.Location,
			},
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		}
		warehouses = append(warehouses, wh)

		cfg := mgwarehouse.Configuration{
			ID:        mongodb.NewObjectID(),
			Warehouse: wh.ID,
			Supplier: mgwarehouse.ConfigSupplier{
				InvoiceDeliveryMethod: doc.InvoiceDeliveryMethod,
			},
			Order: mgwarehouse.ConfigOrder{
				MinimumValue:             doc.MinimumValue,
				PaymentMethod:            doc.PaymentMethods,
				IsLimitNumberOfPurchases: doc.IsLimitNumberOfPurchases,
				LimitNumberOfPurchases:   doc.LimitNumberOfPurchases,
			},
			Partner: mgwarehouse.ConfigPartner{
				IdentityCode:   doc.Partner.IdentityCode,
				Code:           doc.Partner.Code,
				Enabled:        false,
				Authentication: "",
			},
			Delivery: mgwarehouse.ConfigDelivery{
				DeliveryMethods:      doc.DeliveryMethods,
				PriorityServiceCodes: doc.PriorityDeliveryServiceCodes,
				EnabledSources:       doc.EnabledDeliverySources,
				Types:                []string{"selly_delivery"},
			},
			Other: mgwarehouse.ConfigOther{
				DoesSupportSellyExpress: doc.DoesSupportSellyExpress,
			},
		}
		whCfgs = append(whCfgs, cfg)
	}
	_, err = whCol.InsertMany(ctx, warehouses)
	if err != nil {
		fmt.Println("insert wh err:", err)
	}
	_, err = whCfgCol.InsertMany(ctx, whCfgs)
	if err != nil {
		fmt.Println("insert cfg err:", err)
	}
	return c.JSON(200, "ok")
}

func MigrateUpdateStatusWarehousePendingInactive(c echo.Context) error {
	d := dao.Warehouse()
	warehouses := d.FindByCondition(context.Background(), bson.M{"status": constant.WarehouseStatusHoliday})
	if len(warehouses) == 0 {
		return c.JSON(200, "oke")
	}
	listWarehouseUpdatePendingInactive := make([]natsmodel.UpdateStatusWarehousePendingInactive, 0)

	for _, w := range warehouses {
		listWarehouseUpdatePendingInactive = append(listWarehouseUpdatePendingInactive, natsmodel.UpdateStatusWarehousePendingInactive{
			WarehouseID:     w.ID.Hex(),
			PendingInactive: true,
		})
	}

	if err := natsclient.GetWarehouse().UpdateStatusWarehousePendingInactive(natsmodel.UpdateStatusWarehousePendingInactiveRequest{Warehouses: listWarehouseUpdatePendingInactive}); err != nil {
		logger.Error("Error Migrate update status warehouse pending inactive : ", logger.LogData{
			"error": err.Error(),
		})
	}

	return c.JSON(200, "ok")
}

type inventoryRaw struct {
	ID                       primitive.ObjectID              `bson:"_id" json:"_id"`
	Name                     string                          `bson:"name" json:"name"`
	SearchString             string                          `bson:"searchString" json:"-"`
	Code                     int                             `bson:"code" json:"code"`
	Slug                     string                          `bson:"slug" json:"slug"`
	Contact                  mgwarehouse.Contact             `bson:"contact" json:"contact"`
	CanAutoSendEmail         bool                            `bson:"canAutoSendEmail" json:"canAutoSendEmail"`
	CanIssueInvoice          bool                            `bson:"canIssueInvoice" json:"canIssueInvoice"`
	IsLimitNumberOfPurchases bool                            `bson:"isLimitNumberOfPurchases" json:"isLimitNumberOfPurchases"`
	LimitNumberOfPurchases   int64                           `bson:"limitNumberOfPurchases" json:"limitNumberOfPurchases"`
	InvoiceDeliveryMethod    string                          `bson:"invoiceDeliveryMethod" json:"invoiceDeliveryMethod"`
	Active                   bool                            `bson:"active" json:"active"`
	MinimumValue             float64                         `bson:"minimumValue" json:"minimumValue"`
	PaymentMethods           mgwarehouse.ConfigPaymentMethod `bson:"paymentMethods" json:"paymentMethods"`
	Location                 *LocationInventory              `bson:"location" json:"location"`
	Supplier                 struct {
		ID primitive.ObjectID `bson:"_id"`
	} `bson:"supplier,omitempty" json:"supplier"`
	Checksum                     string           `bson:"checksum" json:"checksum"`
	DoesSupportSellyExpress      bool             `bson:"doesSupportSellyExpress" json:"doesSupportSellyExpress"`
	CreatedAt                    time.Time        `bson:"createdAt" json:"createdAt"`
	UpdatedAt                    time.Time        `bson:"updatedAt" json:"updatedAt"`
	Partner                      InventoryPartner `bson:"partner,omitempty" json:"partner"`
	DeliveryMethods              []string         `bson:"deliveryMethods" json:"deliveryMethods"`
	PriorityDeliveryServiceCodes []string         `bson:"priorityDeliveryServiceCodes" json:"priorityDeliveryServiceCodes"`
	EnabledDeliverySources       []int            `bson:"enabledDeliverySources" json:"enabledDeliverySources"`

	LimitedNumberOfProductsPerOrder int `bson:"limitedNumberOfProductsPerOrder" json:"limitedNumberOfProductsPerOrder"`
}

// InventoryPartner ...
type InventoryPartner struct {
	IdentityCode string `json:"identityCode" bson:"identityCode"`
	Code         string `json:"code" bson:"code"`
}

// InventoryPaymentMethods ...
type InventoryPaymentMethods struct {
	Cod          bool `json:"cod" bson:"cod"`
	BankTransfer bool `json:"bankTransfer" bson:"bankTransfer"`
}

// LocationInventory ...
type LocationInventory struct {
	Address      string                  `bson:"address" json:"address"`
	Province     string                  `bson:"province" json:"province"`
	ProvinceName string                  `bson:"provinceName,omitempty" json:"provinceName,omitempty"`
	District     string                  `bson:"district" json:"district"`
	Ward         string                  `bson:"ward" json:"ward"`
	FullAddress  string                  `bson:"fullAddress,omitempty" json:"fullAddress,omitempty"`
	Location     mgwarehouse.Coordinates `bson:"location" json:"location"`
}

// ContactInventory ...
type ContactInventory struct {
	Name    string `bson:"name" json:"name"`
	Phone   string `bson:"phone" json:"phone"`
	Address string `bson:"address" json:"address"`
	Email   string `bson:"email" json:"email"`
}

// MongoLocation ...
type MongoLocation struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}
