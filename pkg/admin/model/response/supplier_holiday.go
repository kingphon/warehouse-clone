package responsemodel

import "git.selly.red/Selly-Server/warehouse/external/utils/ptime"

// ResponseSupplierHolidayAll ...
type ResponseSupplierHolidayAll struct {
	List  []ResponseSupplierHolidayBrief `json:"list"`
	Limit int64                          `json:"limit"`
	Total int64                          `json:"total"`
}

// ResponseSupplierHolidayBrief ...
type ResponseSupplierHolidayBrief struct {
	ID         string                   `json:"_id"`
	Title      string                   `json:"title"`
	From       *ptime.TimeResponse      `json:"from"`
	To         *ptime.TimeResponse      `json:"to"`
	Reason     string                   `json:"reason"`
	Source     string                   `json:"source"`
	Status     string                   `json:"status"`
	IsApplyAll bool                     `json:"isApplyAll"`
	Warehouses []ResponseWarehouseShort `json:"warehouses"`
	Supplier   ResponseSupplierShort    `json:"supplier,omitempty"`
	CreatedBy  ResponseInfo             `json:"createdBy,omitempty"`
}

// ResponseSupplierHolidayDetail ...
type ResponseSupplierHolidayDetail struct {
	ID         string                   `json:"_id"`
	Title      string                   `json:"title"`
	From       *ptime.TimeResponse      `json:"from"`
	To         *ptime.TimeResponse      `json:"to"`
	Reason     string                   `json:"reason"`
	Source     string                   `json:"source"`
	Status     string                   `json:"status"`
	IsApplyAll bool                     `json:"isApplyAll"`
	Warehouses []ResponseWarehouseShort `json:"warehouses"`
	Supplier   ResponseSupplierShort    `json:"supplier,omitempty"`
	CreatedBy  ResponseStaffShort       `json:"createdBy,omitempty"`
}
