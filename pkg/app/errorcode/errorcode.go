package errorcode

import "git.selly.red/Selly-Server/warehouse/external/response"

// Init list error codes
// Code from 200 - ...
// warehouse: 201 - 299
// supplier : 300 - 399
func Init() {
	// Init common code first
	response.Init()
	response.AddListCodes(warehouse)
	response.AddListCodes(supplierHoliday)

}
