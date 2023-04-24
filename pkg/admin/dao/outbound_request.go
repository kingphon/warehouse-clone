package dao

import (
	"context"

	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mgwarehouse "git.selly.red/Selly-Server/warehouse/external/model/mg/warehouse"

	"git.selly.red/Selly-Server/warehouse/pkg/admin/database"
)

// OutboundRequestInterface ...
type OutboundRequestInterface interface {
	FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) []mgwarehouse.OutboundRequest
	FindCursor(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) mgwarehouse.OutboundRequest
	FindOneByORCode(ctx context.Context, code string) mgwarehouse.OutboundRequest
	FindOneByORCodeAndOrderID(ctx context.Context, code string, orderID primitive.ObjectID) mgwarehouse.OutboundRequest

	CountByCondition(ctx context.Context, cond interface{}) int64
	CountByOrderID(ctx context.Context, orderID primitive.ObjectID) int64

	InsertOne(ctx context.Context, payload mgwarehouse.OutboundRequest) error
	UpdateOne(ctx context.Context, cond, payload interface{}) (*mongo.UpdateResult, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, payload interface{}) (*mongo.UpdateResult, error)
}

// OutboundRequest ...
func OutboundRequest() OutboundRequestInterface {
	return outboundRequestImplement{}
}

// outboundRequestImplement ...
type outboundRequestImplement struct{}

// FindOneByORCodeAndOrderID ...
func (d outboundRequestImplement) FindOneByORCodeAndOrderID(ctx context.Context, code string, orderID primitive.ObjectID) mgwarehouse.OutboundRequest {
	return d.FindOne(ctx, bson.D{
		{"partner.code", code},
		{"order", orderID},
	})
}

// CountByOrderID ...
func (d outboundRequestImplement) CountByOrderID(ctx context.Context, orderID primitive.ObjectID) int64 {
	return d.CountByCondition(ctx, bson.M{"order": orderID})
}

// CountByCondition ...
func (d outboundRequestImplement) CountByCondition(ctx context.Context, cond interface{}) int64 {
	total, _ := database.OutboundRequestCol().CountDocuments(ctx, cond)
	return total
}

// FindOneByORCode ...
func (d outboundRequestImplement) FindOneByORCode(ctx context.Context, code string) mgwarehouse.OutboundRequest {
	return d.FindOne(ctx, bson.M{"partner.code": code})
}

// FindOne ...
func (d outboundRequestImplement) FindOne(ctx context.Context, cond interface{}, opts ...*options.FindOneOptions) (doc mgwarehouse.OutboundRequest) {
	if err := database.OutboundRequestCol().FindOne(ctx, cond, opts...).Decode(&doc); err != nil {
		logger.Error("dao.outboundRequestImplement.FindOne", logger.LogData{
			"err":  err.Error(),
			"cond": cond,
		})
	}
	return doc
}

// FindCursor ...
func (d outboundRequestImplement) FindCursor(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (
	*mongo.Cursor, error) {
	return database.OutboundRequestCol().Find(ctx, cond, opts...)
}

// FindByCondition ...
func (d outboundRequestImplement) FindByCondition(ctx context.Context, cond interface{}, opts ...*options.FindOptions) (
	docs []mgwarehouse.OutboundRequest) {
	cursor, err := database.OutboundRequestCol().Find(ctx, cond, opts...)
	if err != nil {
		logger.Error("dao.outboundRequestImplement.FindByCondition - find", logger.LogData{
			"err":  err.Error(),
			"cond": cond,
		})
		return nil
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &docs); err != nil {
		logger.Error("dao.outboundRequestImplement.FindByCondition - decode", logger.LogData{
			"err":  err.Error(),
			"cond": cond,
		})
	}
	return docs
}

// InsertOne ...
func (d outboundRequestImplement) InsertOne(ctx context.Context, doc mgwarehouse.OutboundRequest) error {
	_, err := database.OutboundRequestCol().InsertOne(ctx, doc)
	if err != nil {
		logger.Error("dao.outboundRequestImplement.InsertOne", logger.LogData{
			"err": err.Error(),
			"doc": doc,
		})
	}
	return err
}

// UpdateOne ...
func (d outboundRequestImplement) UpdateOne(ctx context.Context, cond, payload interface{}) (*mongo.UpdateResult, error) {
	res, err := database.OutboundRequestCol().UpdateOne(ctx, cond, payload)
	if err != nil {
		logger.Error("dao.outboundRequestImplement.InsertOne", logger.LogData{
			"err":     err.Error(),
			"cond":    cond,
			"payload": payload,
		})
	}
	return res, err
}

// UpdateByID ...
func (d outboundRequestImplement) UpdateByID(ctx context.Context, id primitive.ObjectID, payload interface{}) (*mongo.UpdateResult, error) {
	return d.UpdateOne(ctx, bson.M{"_id": id}, payload)
}
