package errorcode

import "git.selly.red/Selly-Server/warehouse/external/response"

const (
	WarehouseIsRequiredName    = "warehouse_is_required_name"
	WarehouseInvalidSupplier   = "warehouse_invalid_supplier"
	WarehouseNotFound          = "warehouse_not_found"
	WarehouseInvalidID         = "warehouse_invalid_id"
	WarehouseLocationFindError = "warehouse_location_find_error"
	WarehouseNotBelongSupplier = "warehouse_not_belong_supplier"
	WarehouseErrorWhenUpdate   = "warehouse_err_when_update"
)

//
var warehouse = []response.Code{

	{
		Key:     WarehouseInvalidSupplier,
		Message: "nhà cung cấp không hợp lệ",
		Code:    201,
	},
	{
		Key:     WarehouseNotFound,
		Message: "kho hàng không tồn tại",
		Code:    202,
	},
	{
		Key:     WarehouseInvalidID,
		Message: "id kho hàng không hợp lệ",
		Code:    203,
	},
	{
		Key:     WarehouseLocationFindError,
		Message: "lỗi khi find location",
		Code:    204,
	},
	{
		Key:     WarehouseNotBelongSupplier,
		Message: "kho hàng không thuộc nhà cung cấp",
		Code:    205,
	},
	{
		Key:     WarehouseErrorWhenUpdate,
		Message: "lỗi khi cập nhập kho hàng",
		Code:    206,
	},
}
