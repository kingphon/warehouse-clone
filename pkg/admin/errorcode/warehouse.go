package errorcode

import "git.selly.red/Selly-Server/warehouse/external/response"

const (
	WarehouseIsRequiredName                = "warehouse_is_required_name"
	WarehouseIsRequiredCode                = "warehouse_is_required_code"
	WarehouseIsRequiredSupplier            = "warehouse_is_required_supplier"
	WarehouseIsRequiredContactName         = "warehouse_is_required_contact_name"
	WarehouseIsRequiredContactPhone        = "warehouse_is_required_contact_phone"
	WarehouseIsRequiredContactAddress      = "warehouse_is_required_contact_address"
	WarehouseIsRequiredContactEmail        = "warehouse_is_required_contact_email"
	WarehouseIsRequiredLocationProvince    = "warehouse_is_required_location_province"
	WarehouseIsRequiredLocationDistrict    = "warehouse_is_required_location_district"
	WarehouseIsRequiredLocationWard        = "warehouse_is_required_location_ward"
	WarehouseIsRequiredLocationAddress     = "warehouse_is_required_location_address"
	WarehouseIsRequiredLocationFullAddress = "warehouse_is_required_location_fullAddress"
	WarehouseErrorWhenInsert               = "warehouse_err_when_insert"
	WarehouseInvalidSupplier               = "warehouse_invalid_supplier"
	WarehouseNotFound                      = "warehouse_not_found"
	WarehouseInvalidID                     = "warehouse_invalid_id"
	WarehouseInvalidStatus                 = "warehouse_invalid_status"
	WarehouseExistedName                   = "warehouse_existed_name"
	WarehouseLocationFindError             = "warehouse_location_find_error"
	SupplierContractInvalidStatus          = "supplier_contract_invalid_status"
	WarehouseInvalidProvince               = "warehouse_invalid_province"
	WarehouseInvalidDistrict               = "warehouse_invalid_district"
	WarehouseInvalidWard                   = "warehouse_invalid_ward"
	WarehouseNotBelongSupplier             = "warehouse_not_belong_supplier"
	WarehouseErrorWhenUpdate               = "warehouse_err_when_update"
)

var warehouse = []response.Code{
	{
		Key:     WarehouseIsRequiredName,
		Message: "tên không được trống",
		Code:    211,
	},
	{
		Key:     WarehouseIsRequiredCode,
		Message: "code không được trống",
		Code:    212,
	}, {
		Key:     WarehouseIsRequiredSupplier,
		Message: "nhà cung cấp không được trống",
		Code:    213,
	}, {
		Key:     WarehouseIsRequiredContactName,
		Message: "tên liên lạc không được trống",
		Code:    214,
	}, {
		Key:     WarehouseIsRequiredContactPhone,
		Message: "số điện thoại liên lạc không được trống",
		Code:    214,
	}, {
		Key:     WarehouseIsRequiredContactAddress,
		Message: "địa chỉ liên lạc không được trống",
		Code:    215,
	}, {
		Key:     WarehouseIsRequiredContactEmail,
		Message: "email liên lạc không được trống",
		Code:    216,
	}, {
		Key:     WarehouseIsRequiredLocationProvince,
		Message: "thành phố không được trống",
		Code:    217,
	}, {
		Key:     WarehouseIsRequiredLocationDistrict,
		Message: "quận huyện không được trống",
		Code:    218,
	}, {
		Key:     WarehouseIsRequiredLocationWard,
		Message: "phường/xã không được trống",
		Code:    219,
	}, {
		Key:     WarehouseIsRequiredLocationAddress,
		Message: "địa chỉ không được trống",
		Code:    220,
	}, {
		Key:     WarehouseIsRequiredLocationFullAddress,
		Message: "địa chỉ đầy đủ không được trống",
		Code:    221,
	}, {
		Key:     WarehouseErrorWhenInsert,
		Message: "lỗi khi tạo kho hàng",
		Code:    221,
	}, {
		Key:     WarehouseInvalidSupplier,
		Message: "nhà cung cấp không hợp lệ",
		Code:    222,
	}, {
		Key:     WarehouseNotFound,
		Message: "kho hàng không tồn tại",
		Code:    223,
	}, {
		Key:     WarehouseInvalidID,
		Message: "id kho hàng không hợp lệ",
		Code:    224,
	}, {
		Key:     WarehouseInvalidStatus,
		Message: "trạng thái kho hàng ko hợp lệ",
		Code:    225,
	}, {
		Key:     WarehouseExistedName,
		Message: "tên kho hàng đã tồn tại",
		Code:    226,
	}, {
		Key:     WarehouseLocationFindError,
		Message: "lỗi khi find location",
		Code:    227,
	}, {
		Key:     SupplierContractInvalidStatus,
		Message: "trạng thái hợp đồng nhà cung cấp không hợp lệ",
		Code:    228,
	}, {
		Key:     WarehouseInvalidProvince,
		Message: "mã tỉnht/hành phố không hợp lệ",
		Code:    229,
	}, {
		Key:     WarehouseInvalidDistrict,
		Message: "mã quận/huyện không hợp lệ",
		Code:    230,
	},
	{
		Key:     WarehouseInvalidWard,
		Message: "mã xã/phương không hợp lệ",
		Code:    231,
	},
	{
		Key:     WarehouseNotBelongSupplier,
		Message: "kho hàng không thuộc nhà cung cấp",
		Code:    232,
	},
	{
		Key:     WarehouseErrorWhenUpdate,
		Message: "lỗi khi cập nhập kho hàng",
		Code:    233,
	},
}
