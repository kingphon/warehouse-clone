package errorcode

import "git.selly.red/Selly-Server/warehouse/external/response"

const (
	WarehouseCfgInvalidInvoiceMethod          = "warehouse_configuration_invalid_invoice_method"
	WarehouseCfgInvalidDeliveryMethod         = "warehouse_configuration_invalid_delivery_method"
	WarehouseCfgNotFound                      = "warehouse_configuration_not_found"
	WarehouseCfgTimeRangeInvalid              = "warehouse_configuration_time_range_invalid"
	WarehouseCfgInvalidID                     = "warehouse_configuration_invalid_id"
	WarehouseCfgInvalidDeliveryType           = "warehouse_configuration_invalid_delivery_type"
	WarehouseCfgIsRequiredDeliveryMethods     = "warehouse_cfg_is_required_delivery_methods"
	WarehouseIsRequiredLimitNumberOfPurchases = "warehouse_is_required_limit_number_of_purchases"
	WarehouseCfgIsRequiredEnabledSources      = "warehouse_configuration_is_required_enabled_sources"
	WarehouseCfgIsRequiredTypes               = "warehouse_configuration_is_required_types"
	WarehouseCfgIsRequiredBusinessType        = "warehouse_configuration_is_required_business_type"
	WarehouseCfgInvalidBusinessType           = "warehouse_configuration_invalid_business_type"
)

var warehouseConfiguration = []response.Code{
	{
		Key:     WarehouseCfgInvalidDeliveryMethod,
		Message: "phương thức vận chuyển không hợp lệ",
		Code:    400,
	},
	{
		Key:     WarehouseCfgInvalidInvoiceMethod,
		Message: "phương thức xuất hóa đơn không hợp lệ",
		Code:    401,
	},
	{
		Key:     WarehouseCfgNotFound,
		Message: "Thiết lập kho khàng không tồn tại",
		Code:    402,
	},
	{
		Key:     WarehouseCfgInvalidID,
		Message: "id thiết lập kho hàng không hợp lệ",
		Code:    403,
	}, {
		Key:     WarehouseCfgInvalidDeliveryType,
		Message: "loại vận chuyển hợp lệ",
		Code:    404,
	}, {
		Key:     WarehouseCfgIsRequiredDeliveryMethods,
		Message: "nhà vận chuyển áp dụng không được trống",
		Code:    405,
	}, {
		Key:     WarehouseIsRequiredLimitNumberOfPurchases,
		Message: "sô lượng sản phẩm tạo đơn giới hạn không được trống",
		Code:    406,
	}, {
		Key:     WarehouseCfgIsRequiredEnabledSources,
		Message: "áp dụng nhà vận chuyển không được trống",
		Code:    407,
	}, {
		Key:     WarehouseCfgIsRequiredTypes,
		Message: " loại vận chuyển không được trống",
		Code:    408,
	},
	{
		Key:     WarehouseCfgIsRequiredBusinessType,
		Message: "Loại hình kinh doanh không được trống",
		Code:    409,
	},
	{
		Key:     WarehouseCfgInvalidBusinessType,
		Message: "Loại hình kinh doanh không hợp lệ",
		Code:    410,
	},
	{
		Key:     WarehouseCfgTimeRangeInvalid,
		Message: "Thời gian mở cửa không hợp lệ",
		Code:    411,
	},
}
