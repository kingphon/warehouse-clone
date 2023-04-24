package errorcode

import "git.selly.red/Selly-Server/warehouse/external/response"

const (
	ORMessageAutoCancelOR = "hệ thống tự động hủy xuất kho"
	ORMessageCancelOR     = "hủy xuất kho"
)

const (
	OROrderRequired        = "or_order_required"
	ORRequestRequired      = "or_request_required"
	ORMissingInsuranceInfo = "or_missing_insurance_info"
	ORNotFound             = "or_not_found"
	ORInvalidOrderID       = "or_invalid_order_id"
)

var outboundRequest = []response.Code{
	{
		Key:     OROrderRequired,
		Message: "đơn hàng không được trống",
		Code:    600,
	},
	{
		Key:     ORRequestRequired,
		Message: "yêu cầu xuất kho không được trống",
		Code:    601,
	},
	{
		Key:     ORMissingInsuranceInfo,
		Message: "sku thiếu hoặc sai thông tin bảo hiểm, vui lòng kiểm tra lại",
		Code:    602,
	},
	{
		Key:     ORNotFound,
		Message: "không tìm thấy yêu cầu xuất kho",
		Code:    603,
	},
	{
		Key:     ORInvalidOrderID,
		Message: "id đơn hàng không hợp lệ",
		Code:    604,
	},
}
