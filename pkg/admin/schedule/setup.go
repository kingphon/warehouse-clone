package schedule

import (
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/service"
)

// Init ...
func Init() {
	syncGCStatus := "*/30 * * * *"
	setIsClosed := "*/10 * * * *"
	updateHolidayWarehouse := "0 */3 * * *"
	updateHolidaysStatus := "0 2 * * *"
	updatePaymentMethodBankTransfer := "0 0 27 1 *"
	if config.IsEnvDevelop() {
		syncGCStatus = "*/5 * * * *"
		setIsClosed = "*/5 * * * *"
	}
	scheduleSvc := service.Schedule{}
	jobs := []*Job{
		{
			Spec: syncGCStatus,
			Name: "Sync Global Care",
			Cmd:  service.OutboundRequest().SyncAllGlobalCare,
		},
		{
			Spec: setIsClosed,
			Name: "Job set is closed",
			Cmd:  scheduleSvc.RunJobUpdateIsClosed,
		},
		{
			Spec: updateHolidayWarehouse,
			Name: "Job update holiday warehouse",
			Cmd:  scheduleSvc.RunJobUpdateHolidayWarehouses,
		},
		{
			Spec: updateHolidaysStatus,
			Name: "Job update holidays status",
			Cmd:  scheduleSvc.RunJobUpdateHolidayStatusForSupplier,
		},
		{
			Spec: updatePaymentMethodBankTransfer,
			Name: "Job update payment method bank transfer",
			Cmd:  scheduleSvc.UpdatePaymentMethodBankTransferWarehouse,
		},
	}
	New(jobs...).Start()
}
