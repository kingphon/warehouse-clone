package locale

import (
	"fmt"
)

const (
	SupplierHolidayTitle  = "Kỳ nghỉ lễ của %s"
	SupplierHolidayReason = "Tạm ngưng bán đến ngày %s"
)

// GetSupplierHolidayTitle ...
func GetSupplierHolidayTitle(supplierName string) string {
	return fmt.Sprintf(SupplierHolidayTitle, supplierName)
}

// GetSupplierHolidayReason ...
func GetSupplierHolidayReason(toAt string) string {
	return fmt.Sprintf(SupplierHolidayReason, toAt)
}
