package errorcode

import "git.selly.red/Selly-Server/warehouse/external/response"

const (
	SupplierHolidayInvalidSupplierId  = "supplier_holiday_invalid_supplier_id"
	SupplierHolidayInvalidSource      = "supplier_holiday_invalid_source"
	SupplierHolidayInvalidStatus      = "supplier_holiday_invalid_status"
	SupplierHolidayInvalidID          = "supplier_holiday_invalid_id"
	SupplierHolidayErrorWhenCreate    = "supplier_holiday_error_when_create"
	SupplierHolidayNotFound           = "supplier_holiday_not_found"
	SupplierHolidayErrorWhenUpdate    = "supplier_holiday_error_when_update"
	SupplierHolidayInvalidWarehouseId = "supplier_holiday_invalid_warehouse_id"
	SupplierHolidayIsRequireTitle     = "supplier_holiday_is_required_title"
	SupplierHolidayIsRequireReason    = "supplier_holiday_is_required_reason"
	SupplierHolidayInValidFromToAt    = "supplier_holiday_invalid_from_to_at"
	SupplierHolidayIsExitedHoliday    = "supplier_holiday_is_exited_holiday"
)

// supplierHoliday ...
var supplierHoliday = []response.Code{
	{
		Key:     SupplierHolidayInvalidSupplierId,
		Message: "nhà cung cấp không hợp lệ ",
		Code:    800,
	},
	{
		Key:     SupplierHolidayInvalidSource,
		Message: "source không hợp lệ",
		Code:    801,
	},
	{
		Key:     SupplierHolidayInvalidStatus,
		Message: "trạng thái không hợp lệ",
		Code:    802,
	},
	{
		Key:     SupplierHolidayInvalidID,
		Message: "id kì nghỉ lễ không hợp lệ",
		Code:    803,
	},
	{
		Key:     SupplierHolidayErrorWhenCreate,
		Message: "lỗi khi tạo mới",
		Code:    804,
	},
	{
		Key:     SupplierHolidayNotFound,
		Message: "kì nghỉ lễ không tồn tại",
		Code:    805,
	},
	{
		Key:     SupplierHolidayErrorWhenUpdate,
		Message: "lỗi khi cập nhật",
		Code:    806,
	},
	{
		Key:     SupplierHolidayInvalidWarehouseId,
		Message: "id kho hàng không hợp lệ",
		Code:    807,
	},
	{
		Key:     SupplierHolidayIsRequireTitle,
		Message: "tiêu đề không được trống",
		Code:    808,
	},
	{
		Key:     SupplierHolidayIsRequireReason,
		Message: "lý do không được trống",
		Code:    809,
	},
	{
		Key:     SupplierHolidayInValidFromToAt,
		Message: "thời gian nghỉ lễ không hợp lệ",
		Code:    810,
	},
	{
		Key:     SupplierHolidayIsExitedHoliday,
		Message: "Nhà cung cấp đã có kì nghỉ lễ",
		Code:    811,
	},
}
