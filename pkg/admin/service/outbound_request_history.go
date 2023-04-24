package service

import (
	"context"

	"git.selly.red/Selly-Modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/external/utils/ptime"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/dao"
	requestmodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/request"
	responsemodel "git.selly.red/Selly-Server/warehouse/pkg/admin/model/response"
)

// OutboundRequestHistoryInterface ...
type OutboundRequestHistoryInterface interface {
	GetList(ctx context.Context, q requestmodel.OutboundRequestHistoryQuery) ([]responsemodel.OutboundRequestHistory, error)

	SaveByOR(ctx context.Context, or mgwarehouse.OutboundRequest) (mgwarehouse.OutboundRequestHistory, error)
}

func OutboundRequestHistory() OutboundRequestHistoryInterface {
	return outboundRequestHistoryImplement{}
}

type outboundRequestHistoryImplement struct{}

// GetList ...
func (s outboundRequestHistoryImplement) GetList(ctx context.Context, q requestmodel.OutboundRequestHistoryQuery) (
	[]responsemodel.OutboundRequestHistory, error) {
	requestID, _ := mongodb.NewIDFromString(q.Request)
	cond := bson.M{"request": requestID}
	docs := dao.OutboundRequestHistory().FindByCondition(ctx, cond)
	res := make([]responsemodel.OutboundRequestHistory, len(docs))
	for i, doc := range docs {
		res[i] = s.getResponse(doc)
	}
	return res, nil
}

// SaveByOR ...
func (s outboundRequestHistoryImplement) SaveByOR(ctx context.Context, or mgwarehouse.OutboundRequest) (
	mgwarehouse.OutboundRequestHistory, error) {
	h := mgwarehouse.OutboundRequestHistory{
		ID:        primitive.NewObjectID(),
		Order:     or.Order,
		Request:   or.ID,
		Status:    or.Status,
		CreatedAt: ptime.Now(),
	}
	err := dao.OutboundRequestHistory().InsertOne(ctx, h)
	return h, err
}

func (s outboundRequestHistoryImplement) getResponse(doc mgwarehouse.OutboundRequestHistory) responsemodel.OutboundRequestHistory {
	return responsemodel.OutboundRequestHistory{
		ID:        doc.ID.Hex(),
		Status:    doc.Status,
		CreatedAt: ptime.TimeResponseInit(doc.CreatedAt),
	}
}
