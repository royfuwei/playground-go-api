package mgo

import (
	"context"
	"playground-go-api/domain"
	"playground-go-api/infrastructures/mongodb"
	"playground-go-api/infrastructures/tools"
	"time"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mgoConfigsRepository struct {
	baseRepo *mongodb.BaseRepository
}

func NewMgoConfigsRepository(client *mongo.Client) domain.ConfigsRepository {
	collection := client.Database("auth").Collection("configs")
	baseRepo := mongodb.NewBaseRepository(collection)
	return &mgoConfigsRepository{
		baseRepo,
	}
}

func (m *mgoConfigsRepository) FindOneByFilter(filter primitive.M) (*domain.Config, error) {
	var config *domain.Config
	result, err := m.baseRepo.FindOneByFilter(filter)
	if err != nil {
		return nil, err
	}
	if err := result.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (m *mgoConfigsRepository) FindByName(name string) (config *domain.Config, err error) {
	filter := bson.M{
		"name": name,
	}
	var configs []*domain.Config
	cur, total, err := m.baseRepo.Find(filter, 5*time.Second)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var result *domain.Config
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		configs = append(configs, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	if total > 0 {
		config = configs[0]
		return config, nil
	}
	return config, nil
}

// 初始化 initial/configs.json 寫入db
func (m *mgoConfigsRepository) InitFindNameAndUpdate(data *domain.Config) (*domain.Config, error) {
	modiTime := primitive.NewDateTimeFromTime(time.Now())
	data.ModiTime = modiTime
	filter := bson.M{
		"name": data.Name,
	}
	update := bson.M{
		"$set": data,
	}
	opts := &options.FindOneAndUpdateOptions{
		Upsert: tools.Bool(true),
	}
	result, err := m.baseRepo.FindOneAndUpdate(filter, update, opts)
	var config *domain.Config
	if err := result.Decode(&config); err != nil {
		glog.Error(err)
		return nil, err
	}
	return config, err
}
