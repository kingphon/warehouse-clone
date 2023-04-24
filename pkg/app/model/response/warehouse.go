package responsemodel

import "git.selly.red/Selly-Modules/natsio/model"

// ResponseWarehouseLocation ...
type ResponseWarehouseLocation struct {
	Province            ResponseWarehouseProvince `json:"province"`
	District            ResponseWarehouseDistrict `json:"district"`
	Ward                ResponseWarehouseWard     `json:"ward"`
	Address             string                    `json:"address"`
	LocationCoordinates ResponseLatLng            `json:"locationCoordinates"`
}

// ResponseWarehouseProvince ...
type ResponseWarehouseProvince struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code int    `json:"code"`
}

// ResponseWarehouseDistrict ...
type ResponseWarehouseDistrict struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code int    `json:"code"`
}

// ResponseWarehouseWard ...
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
	ID                              string                  `json:"_id"`
	Name                            string                  `json:"name"`
	Notices                         []model.NewsAppResponse `json:"notices"`
	Code                            int                     `json:"id"`
	CanIssueInvoice                 bool                    `json:"canIssueInvoice"`
	InvoiceDeliveryMethod           string                  `json:"invoiceDeliveryMethod"`
	DoesSupportSellyExpress         bool                    `json:"doesSupportSellyExpress"`
	LimitedNumberOfProductsPerOrder int                     `json:"limitedNumberOfProductsPerOrder"`
}

// ResponseWarehouseInfo ...
type ResponseWarehouseInfo struct {
	ID     string `json:"_id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
