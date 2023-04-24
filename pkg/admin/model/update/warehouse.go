package updatemodel

import (
	"time"
)

type WarehouseUpdate struct {
	Name         string            `bson:"name"`
	BusinessType string            `bson:"businessType"`
	SearchString string            `bson:"searchString"`
	Slug         string            `bson:"slug"`
	Contact      WarehouseContact  `bson:"contact"`
	Location     LocationWarehouse `bson:"location"`
	UpdatedAt    time.Time         `bson:"updatedAt"`
}

// LocationWarehouse ...
type LocationWarehouse struct {
	Province            int                  `bson:"province"`
	District            int                  `bson:"district"`
	Ward                int                  `bson:"ward"`
	Address             string               `bson:"address"`
	LocationCoordinates WarehouseCoordinates `bson:"location"`
}

// WarehouseCoordinates ...
type WarehouseCoordinates struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

// WarehouseContact ...
type WarehouseContact struct {
	Name    string `bson:"name"`
	Phone   string `bson:"phone"`
	Address string `bson:"address"`
	Email   string `bson:"email"`
}
