package mongodb

import (
	"context"
	"time"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// BaseRepository 各模組功能Repository的基底， 可以把共用的查詢都整理到這邊來。
type BaseRepository struct {
	Collection *mongo.Collection
}

// NewBaseRepository ...
func NewBaseRepository(collection *mongo.Collection) *BaseRepository {
	return &BaseRepository{Collection: collection}
}

func (b *BaseRepository) FindOneAndUpdate(filter bson.M, update bson.M, opts ...*options.FindOneAndUpdateOptions) (*mongo.SingleResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	returnDocument := options.After
	defaultOpts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDocument,
	}
	opts = append(opts, defaultOpts)
	return b.Collection.FindOneAndUpdate(ctx, filter, update, opts...), nil
}

// FindOneByFilter 共用
func (b *BaseRepository) FindOneByFilter(filter bson.M, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return b.Collection.FindOne(ctx, filter, opts...), nil
}

// FindByID 共用 由ID來找出文檔
func (b *BaseRepository) FindByID(id string, opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		b.restartHandler(err)
		return nil, err
	}
	filter := bson.M{
		"_id": objID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return b.Collection.FindOne(ctx, filter, opts...), nil
}

// DeleteByID 共用 由ID刪除文檔
func (b *BaseRepository) DeleteByID(id string, timeout time.Duration, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		b.restartHandler(err)
		return nil, err
	}
	filter := bson.M{
		"_id": objID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return b.Collection.DeleteOne(ctx, filter, opts...)
}

func (b *BaseRepository) DeleteByFilter(filter bson.M, timeout time.Duration, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return b.Collection.DeleteMany(ctx, filter, opts...)
}

// UpdateByID 共用 由ID更新文檔
func (b *BaseRepository) UpdateByID(id string, update bson.M, timeout time.Duration, opts ...*options.FindOneAndUpdateOptions) (*mongo.SingleResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		b.restartHandler(err)
		return nil, err
	}
	filter := bson.M{
		"_id": objID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	returnDocument := options.After
	defaultOpts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDocument,
	}
	opts = append(opts, defaultOpts)
	return b.Collection.FindOneAndUpdate(ctx, filter, update, opts...), nil
}

// IsExistByField 共用 找文檔欄位是否存在需要的值
func (b *BaseRepository) IsExistByField(fieldName string, fieldValue string, timeout time.Duration, opts ...*options.CountOptions) (bool, error) {
	filter := bson.M{
		fieldName: fieldValue,
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	count, err := b.Collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		b.restartHandler(err)
		return false, err
	}
	isExist := count > 0
	return isExist, nil
}

// Add 共用 新增文檔
func (b *BaseRepository) Add(body interface{}, timeout time.Duration, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return b.Collection.InsertOne(ctx, body, opts...)
}

// BulkUpdateByIDs 共用 批次更新
func (b *BaseRepository) BulkUpdateByIDs(ids []string, update bson.M, timeout time.Duration, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	objIDs, err := ObjIds(ids)
	if err != nil {
		b.restartHandler(err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var models []mongo.WriteModel
	for _, objID := range objIDs {
		model := mongo.NewUpdateOneModel().SetFilter(bson.M{
			"_id": objID,
		}).SetUpdate(update)
		models = append(models, model)
	}
	result, err := b.Collection.BulkWrite(ctx, models, opts...)
	if err != nil {
		b.restartHandler(err)
		return nil, err
	}
	return result, nil
}

// Find 搜尋 documents
// 1.*mongo.Cursor 讓各自repository產生對應的[]results
// 2.ctx 供各自repository 使用
// 3.total filter docs 數量
// 4.err error
func (b *BaseRepository) Find(filter bson.M, timeout time.Duration, opts ...*options.FindOptions) (cur *mongo.Cursor, total int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	total, err = b.Collection.CountDocuments(ctx, filter)
	if err != nil {
		b.restartHandler(err)
		return nil, 0, err
	}
	cur, err = b.Collection.Find(ctx, filter, opts...)
	if err != nil {
		glog.Error(err)
		b.restartHandler(err)
		return nil, 0, err
	}
	return cur, total, nil
}

// mongodb 失去連線的錯誤，目前測試最先會隨機出現以下錯誤
// 1. connection(192.168.2.30:30270[-4]) incomplete read of message header: read tcp 192.168.2.30:51968->192.168.2.30:30270: read: connection reset by peer
// 2. context.DeadlineExceeded
// 3. mongo.ErrClientDisconnected # 目前都是前兩個錯誤會先出來...
func (b *BaseRepository) restartHandler(err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pingErr := b.Collection.Database().Client().Ping(ctx, readpref.Primary())
	if pingErr != nil {
		glog.Errorf("error: %v \n", err)
		glog.Errorf("ping error: %v \n", pingErr)
		glog.Fatal("restart inu-label-service!! \n")
	}
}
