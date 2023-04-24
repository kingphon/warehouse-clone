package responsemodel

import "git.selly.red/Selly-Server/warehouse/external/utils/ptime"

// ResponseSupplierHolidayDetail ...
type ResponseSupplierHolidayDetail struct {
	ID         string                  `json:"_id"`
	Title      string                  `json:"title"`
	From       *ptime.TimeResponse     `json:"from"`
	To         *ptime.TimeResponse     `json:"to"`
	Reason     string                  `json:"reason"`
	Source     string                  `json:"source"`
	Status     string                  `json:"status"`
	IsApplyAll bool                    `json:"isApplyAll"`
	Warehouses []ResponseWarehouseInfo `json:"warehouses"`
}
