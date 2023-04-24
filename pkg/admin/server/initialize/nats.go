package initialize

import (
	"time"

	"git.selly.red/Selly-Modules/natsio"
	"git.selly.red/Selly-Modules/natsio/subject"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/handler"
)

func nats() {
	cfg := config.GetENV().Nats
	err := natsio.Connect(natsio.Config{
		URL:            cfg.URL,
		User:           cfg.Username,
		Password:       cfg.Password,
		RequestTimeout: 2 * time.Minute,
	})
	if err != nil {
		panic(err)
	}

	s := natsio.GetServer()
	c, err := s.NewJSONEncodedConn()
	if err != nil {
		panic(err)
	}
	h := handler.Nats{EncodedConn: c}

	subj := subject.Warehouse
	c.QueueSubscribe(subj.CreateOutboundRequest, subj.CreateOutboundRequest, h.CreateOutboundRequest)
	c.QueueSubscribe(subj.UpdateOutboundRequestLogistic, subj.UpdateOutboundRequestLogistic, h.UpdateOutboundRequestLogistic)
	c.QueueSubscribe(subj.CancelOutboundRequest, subj.CancelOutboundRequest, h.CancelOutboundRequest)
	c.QueueSubscribe(subj.GetConfiguration, subj.GetConfiguration, h.GetWarehouseConfiguration)
	c.QueueSubscribe(subj.SyncORStatus, subj.SyncORStatus, h.SyncORStatus)
	c.QueueSubscribe(subj.WebhookTNC, subj.WebhookTNC, h.NewTNCWebhook)
	c.QueueSubscribe(subj.WebhookGlobalCare, subj.WebhookGlobalCare, h.NewGlobalCareWebhook)
	c.QueueSubscribe(subj.WebhookOnPoint, subj.WebhookOnPoint, h.NewOnPointWebhook)
	c.QueueSubscribe(subj.GetWarehouses, subj.GetWarehouses, h.GetWarehouses)
	c.QueueSubscribe(subj.UpdateORDeliveryStatus, subj.UpdateORDeliveryStatus, h.UpdateOutboundRequestLogistic)

	// Nats server sub
	s.QueueSubscribe(subj.FindOne, subj.FindOne, h.GetOneWarehouse)
	s.QueueSubscribe(subj.FindByCondition, subj.FindByCondition, h.GetWarehouseWithCondition)
	s.QueueSubscribe(subj.Distinct, subj.Distinct, h.DistinctWarehouseWithField)
	s.QueueSubscribe(subj.Count, subj.Count, h.CountWarehouseWithCondition)

}
