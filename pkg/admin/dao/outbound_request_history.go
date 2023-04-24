package dao

import (
	"context"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/mongo/options"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/database"
)

// OutboundRequestHistoryInterface ...
type OutboundRequestHistoryInterface interface {
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) []mgwarehouse.OutboundRequestHistory
	FindOne(ctx context.Context, cond interface{}) mgwarehouse.OutboundRequestHistory

	InsertOne(ctx context.Context, payload mgwarehouse.OutboundRequestHistory) error
}

// OutboundRequestHistory ...
func OutboundRequestHistory() OutboundRequestHistoryInterface {
	return outboundRequestHistoryImplement{}
}

// outboundRequestHistoryImplement ...
type outboundRequestHistoryImplement struct{}

// FindOne ...
func (d outboundRequestHistoryImplement) FindOne(ctx context.Context, cond interface{}) (doc mgwarehouse.OutboundRequestHistory) {
	if err := database.OutboundRequestHistoryCol().FindOne(ctx, cond).Decode(&doc); err != nil {
		logger.Error("dao.outboundRequestHistoryImplement.FindOne", logger.LogData{
			"err":  err.Error(),
			"cond": cond,
		})
	}
	return doc
}

// FindByCondition ...
func (d outboundRequestHistoryImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (
	docs []mgwarehouse.OutboundRequestHistory) {
	cursor, err := database.OutboundRequestHistoryCol().Find(ctx, cond, opts...)
	if err != nil {
		logger.Error("dao.outboundRequestHistoryImplement.FindByCondition - find", logger.LogData{
			"err":  err.Error(),
			"cond": cond,
		})
		return nil
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &docs); err != nil {
		logger.Error("dao.outboundRequestHistoryImplement.FindByCondition - decode", logger.LogData{
			"err":  err.Error(),
			"cond": cond,
		})
	}
	return docs
}

// InsertOne ...
func (d outboundRequestHistoryImplement) InsertOne(ctx context.Context, doc mgwarehouse.OutboundRequestHistory) error {
	_, err := database.OutboundRequestHistoryCol().InsertOne(ctx, doc)
	if err != nil {
		logger.Error("dao.outboundRequestHistoryImplement.InsertOne", logger.LogData{
			"err": err.Error(),
			"doc": doc,
		})
	}
	return err
}
