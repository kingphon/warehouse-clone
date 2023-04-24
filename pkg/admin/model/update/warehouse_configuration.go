package updatemodel

import mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"

// WarehouseCfgFoodUpdate ...
type WarehouseCfgFoodUpdate struct {
	ForceClosed bool        `bson:"forceClosed"`
	IsClosed    bool        `bson:"isClosed"`
	TimeRange   []TimeRange `bson:"timeRange"`
}

// TimeRange ...
type TimeRange struct {
	From int64 `bson:"from"`
	To   int64 `bson:"to"`
}

// WarehouseCfgSupplierUpdate ...
type WarehouseCfgSupplierUpdate struct {
	InvoiceDeliveryMethod string `bson:"invoiceDeliveryMethod"`
}

// WarehouseCfgPartnerUpdate ...
type WarehouseCfgPartnerUpdate struct {
	IdentityCode   string `bson:"identityCode"`
	Code           string `bson:"code"`
	Enabled        bool   `bson:"enabled"`
	Authentication string `bson:"authentication"`
}

// WarehouseCfgOrderUpdate ...
type WarehouseCfgOrderUpdate struct {
	MinimumValue             float64                            `bson:"minimumValue"`
	PaymentMethod            WarehouseCfgPaymentMethodUpdate    `bson:"paymentMethod"`
	IsLimitNumberOfPurchases bool                               `bson:"isLimitNumberOfPurchases"`
	LimitNumberOfPurchases   int64                              `bson:"limitNumberOfPurchases"`
	NotifyOnNewOrder         mgwarehouse.ConfigNotifyOnNewOrder `bson:"notifyOnNewOrder"`
}

// WarehouseCfgPaymentMethodUpdate ...
type WarehouseCfgPaymentMethodUpdate struct {
	Cod          bool `bson:"cod"`
	BankTransfer bool `bson:"bankTransfer"`
}

// WarehouseCfgDeliveryUpdate ...
type WarehouseCfgDeliveryUpdate struct {
	DeliveryMethods      []string `bson:"deliveryMethods"`
	PriorityServiceCodes []string `bson:"priorityServiceCodes"`
	EnabledSources       []int    `bson:"enabledSources"`
	Types                []string `bson:"types"`
}

// WarehouseCfgOtherUpdate ...
type WarehouseCfgOtherUpdate struct {
	DoesSupportSellyExpress bool `bson:"doesSupportSellyExpress"`
}
