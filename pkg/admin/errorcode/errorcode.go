package errorcode

import "git.selly.red/Selly-Server/warehouse/external/response"

// Init list error codes
// Code from 200-299
// audit: 100-199
// warehouse: 200-399
// Configurations: 400- 599
// outboundRequest: 600 - 799
// supplerHoliday: 800 - 999
func Init() {
	// Init common code first
	response.Init()
	response.AddListCodes(audit)
	response.AddListCodes(warehouse)
	response.AddListCodes(warehouseConfiguration)
	response.AddListCodes(outboundRequest)
	response.AddListCodes(supplierHoliday)
}
