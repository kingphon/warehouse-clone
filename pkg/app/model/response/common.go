package responsemodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type Upsert struct {
	ID string `json:"_id"`
}

// User ...
type User struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}

// ResponseCreate ...
type ResponseCreate struct {
	ID string `json:"_id"`
}

// ResponseUpdate ...
type ResponseUpdate struct {
	ID string `json:"_id"`
}

// ResponseChangeStatus ...
type ResponseChangeStatus struct {
	ID     string `json:"_id"`
	Status string `json:"status"`
}

// ResponseSupplierInfo ...
type ResponseSupplierInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ResponseUserInfo ...
type ResponseUserInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	SupplierID string `json:"supplierId"`
}
