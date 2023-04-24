package responsemodel

type Upsert struct {
	ID string `json:"_id"`
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

// ResponseSupplierShort ...
type ResponseSupplierShort struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

// ResponseStaffShort ..
type ResponseStaffShort struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

// ResponseInfo ...
type ResponseInfo struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}
