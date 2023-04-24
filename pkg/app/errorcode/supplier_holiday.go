package errorcode

import "git.selly.red/Selly-Server/warehouse/external/response"

const (
	SupplierHolidayInvalidSupplierId                   = "supplier_holiday_invalid_supplier_id"
	SupplierHolidayInvalidSource                       = "supplier_holiday_invalid_source"
	SupplierHolidayInvalidStatus                       = "supplier_holiday_invalid_status"
	SupplierHolidayInvalidID                           = "supplier_holiday_invalid_id"
	SupplierHolidayErrorWhenCreate                     = "supplier_holiday_error_when_create"
	SupplierHolidayNotFound                            = "supplier_holiday_not_found"
	SupplierHolidayErrorWhenUpdate                     = "supplier_holiday_error_when_update"
	SupplierHolidayInvalidWarehouseId                  = "supplier_holiday_invalid_warehouse_id"
	SupplierHolidayIsRequireTitle                      = "supplier_holiday_is_required_title"
	SupplierHolidayIsRequireReason                     = "supplier_holiday_is_required_reason"
	SupplierHolidayInValidFromToAt                     = "supplier_holiday_invalid_from_to_at"
	SupplierHolidayIsExitedHoliday                     = "supplier_holiday_is_exited_holiday"
	SupplierHolidayCanNotCreateHolidayForOtherSupplier = "supplier_holiday_can_not_create_holiday_for_other_supplier"
	SupplierHolidayMusHaveAtLeastOneWarehouse          = "supplier_holiday_have_at_least_one_warehouse"
)

// supplierHoliday ...
var supplierHoliday = []response.Code{
	{
		Key:     SupplierHolidayInvalidSupplierId,
		Message: "nhà cung cấp không hợp lệ ",
		Code:    300,
	},
	{
		Key:     SupplierHolidayInvalidSource,
		Message: "source không hợp lệ",
		Code:    301,
	},
	{
		Key:     SupplierHolidayInvalidStatus,
		Message: "trạng thái không hợp lệ",
		Code:    302,
	},
	{
		Key:     SupplierHolidayInvalidID,
		Message: "id kì nghỉ lễ không hợp lệ",
		Code:    303,
	},
	{
		Key:     SupplierHolidayErrorWhenCreate,
		Message: "lỗi khi tạo mới",
		Code:    304,
	},
	{
		Key:     SupplierHolidayNotFound,
		Message: "kì nghỉ lễ không tồn tại",
		Code:    305,
	},
	{
		Key:     SupplierHolidayErrorWhenUpdate,
		Message: "lỗi khi cập nhật",
		Code:    306,
	},
	{
		Key:     SupplierHolidayInvalidWarehouseId,
		Message: "id kho hàng không hợp lệ",
		Code:    307,
	},
	{
		Key:     SupplierHolidayIsRequireTitle,
		Message: "tiêu đề không được trống",
		Code:    308,
	},
	{
		Key:     SupplierHolidayIsRequireReason,
		Message: "lý do không được trống",
		Code:    309,
	},
	{
		Key:     SupplierHolidayInValidFromToAt,
		Message: "thời gian nghỉ lễ không hợp lệ",
		Code:    310,
	},
	{
		Key:     SupplierHolidayIsExitedHoliday,
		Message: "Nhà cung cấp đã có kì nghỉ lễ",
		Code:    311,
	},
	{
		Key:     SupplierHolidayCanNotCreateHolidayForOtherSupplier,
		Message: "Nhà cung cấp không thể tạo kì nghỉ lễ cho nhà cung cấp khác",
		Code:    313,
	},
	{
		Key:     SupplierHolidayMusHaveAtLeastOneWarehouse,
		Message: "Phải có ít nhất 1 kho được thiết lập ",
		Code:    314,
	},
}
