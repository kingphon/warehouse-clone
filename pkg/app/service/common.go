package service

import (
	"git.selly.red/Selly-Modules/natsio/client"
	"git.selly.red/Selly-Modules/natsio/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CheckTokenSupplierUser ...
func CheckTokenSupplierUser(token string, permissions []string) (*model.ResponseCheckTokenSupplierUser, error) {
	body := model.CheckTokenSupplierUserPayload{
		Token:       token,
		Permissions: permissions,
	}

	return client.GetSupplierUser().CheckTokenSupplierUser(body)
}

// GetSupplierByIDs ...
func GetSupplierByIDs(listID []primitive.ObjectID) ([]*model.ResponseSupplierInfo, error) {
	body := model.GetSupplierRequest{ListID: listID}
	return client.GetSupplier().GetListSupplierInfo(body)
}
