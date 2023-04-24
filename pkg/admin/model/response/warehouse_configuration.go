package responsemodel

import mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"

// ResponseWarehouseConfigurationDetail ...
type ResponseWarehouseConfigurationDetail struct {
	ID               string                          `json:"id"`
	Warehouse        string                          `json:"warehouse"`
	Supplier         ResponseWarehouseSupplierConfig `json:"supplier"`
	Order            ResponseWarehouseOrderConfig    `json:"order"`
	Partner          ResponseWarehousePartnerConfig  `json:"partner"`
	Delivery         ResponseWarehouseDeliveryConfig `json:"delivery"`
	Other            ResponseWarehouseOtherConfig    `json:"other"`
	Food             ResponseWarehousePartnerFood    `json:"food"`
	AutoConfirmOrder mgwarehouse.ConfigOrderConfirm  `json:"autoConfirmOrder"`
}

// ResponseWarehousePartnerFood ...
type ResponseWarehousePartnerFood struct {
	ForceClosed bool        `json:"forceClosed"`
	IsClosed    bool        `json:"isClosed"`
	TimeRange   []TimeRange `json:"timeRange"`
}

// TimeRange ...
type TimeRange struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

// ResponseWarehouseSupplierConfig ...
type ResponseWarehouseSupplierConfig struct {
	InvoiceDeliveryMethod string `json:"invoiceDeliveryMethod"`
}

// ResponseWarehouseOrderConfig ...
type ResponseWarehouseOrderConfig struct {
	MinimumValue             float64                              `json:"minimumValue"`
	PaymentMethod            ResponseWarehousePaymentMethodConfig `json:"paymentMethod"`
	IsLimitNumberOfPurchases bool                                 `json:"isLimitNumberOfPurchases"`
	LimitNumberOfPurchases   int64                                `json:"limitNumberOfPurchases"`
	NotifyOnNewOrder         mgwarehouse.ConfigNotifyOnNewOrder   `json:"notifyOnNewOrder"`
}

// ResponseWarehousePaymentMethodConfig ...
type ResponseWarehousePaymentMethodConfig struct {
	Cod          bool `json:"cod"`
	BankTransfer bool `json:"bankTransfer"`
}

// ResponseWarehouseDeliveryConfig ...
type ResponseWarehouseDeliveryConfig struct {
	DeliveryMethods      []string `json:"deliveryMethods"`
	PriorityServiceCodes []string `json:"priorityServiceCodes"`
	EnabledSources       []int    `json:"enabledSources"`
	Types                []string `json:"types"`
}

// ResponseWarehousePartnerConfig ...
type ResponseWarehousePartnerConfig struct {
	IdentityCode   string `json:"identityCode"`
	Code           string `json:"code"`
	Enabled        bool   `json:"enabled"`
	Authentication string `json:"authentication"`
}

// ResponseWarehouseOtherConfig ...
type ResponseWarehouseOtherConfig struct {
	DoesSupportSellyExpress bool `json:"doesSupportSellyExpress"`
}
