package service

import (
	"git.selly.red/Selly-Modules/natsio/client"
	"git.selly.red/Selly-Modules/natsio/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetSupplierByIDs ...
func GetSupplierByIDs(listID []primitive.ObjectID) ([]*model.ResponseSupplierInfo, error) {
	body := model.GetSupplierRequest{ListID: listID}
	return client.GetSupplier().GetListSupplierInfo(body)
}

// GetStaffByIDs ...
func GetStaffByIDs(staffIds []string) (*model.ResponseListStaffInfo, error) {
	body := model.GetListStaffRequest{
		StaffIds: staffIds,
	}
	return client.GetStaff().GetListStaffInfoByIds(body)
}
