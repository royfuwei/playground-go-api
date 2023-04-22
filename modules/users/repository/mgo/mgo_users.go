package mgo

import (
	"context"
	"fmt"
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

type mgoUsersRepository struct {
	baseRepo       *mongodb.BaseRepository
	collectionName string
}

func NewMgoUsersRepository(client *mongo.Client) domain.UsersRepository {
	collectionName := "users"
	collection := client.Database("auth").Collection(collectionName)
	baseRepo := mongodb.NewBaseRepository(collection)
	m := &mgoUsersRepository{
		baseRepo,
		collectionName,
	}
	if err := m.createIndices(); err != nil {
		glog.Fatal(err)
	}
	return m
}

func (m *mgoUsersRepository) FindFilterAndUpdate(filter primitive.M, data *domain.User) (*domain.User, error) {
	modiTime := primitive.NewDateTimeFromTime(time.Now())
	data.ModiTime = modiTime
	update := bson.M{
		"$set": data,
	}
	opts := &options.FindOneAndUpdateOptions{
		Upsert: tools.Bool(true),
	}
	result, err := m.baseRepo.FindOneAndUpdate(filter, update, opts)
	var user *domain.User
	if err := result.Decode(&user); err != nil {
		glog.Error(err)
		return nil, err
	}
	return user, err
}

func (m *mgoUsersRepository) Add(data *domain.User) (*domain.User, error) {
	modiTime := primitive.NewDateTimeFromTime(time.Now())
	data.ModiTime = modiTime
	insertResult, err := m.baseRepo.Add(data, 5*time.Second)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	user := &domain.User{
		ID:         insertResult.InsertedID.(primitive.ObjectID),
		Account:    data.Account,
		NickName:   data.NickName,
		Telephone:  data.Telephone,
		Roles:      data.Roles,
		Password:   data.Password,
		ExpiryDate: data.ExpiryDate,
		ModiTime:   modiTime,
	}
	return user, nil
}

func (m *mgoUsersRepository) FindOneByFilter(filter bson.M) (*domain.User, error) {
	var user *domain.User
	result, err := m.baseRepo.FindOneByFilter(filter)
	if err != nil {
		return nil, err
	}

	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (m *mgoUsersRepository) IsExistByFilter(filter bson.M) (bool, error) {
	_, err := m.FindOneByFilter(filter)
	fmt.Println("err: ", err)
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 檢查欄位的值是否存在
func (m *mgoUsersRepository) IsExistByField(fieldName string, fieldValue string) (bool, error) {
	isExist, err := m.baseRepo.IsExistByField(fieldName, fieldValue, 5*time.Second)
	if err != nil {
		glog.Error(err)
		return false, err
	}
	return isExist, nil
}

func (m *mgoUsersRepository) find(filter bson.M, opts ...*options.FindOptions) (users []*domain.User, total int64, err error) {
	cur, total, err := m.baseRepo.Find(filter, 5*time.Second, opts...)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var result *domain.User
		if err := cur.Decode(&result); err != nil {
			return nil, 0, err
		}
		users = append(users, result)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (m *mgoUsersRepository) FindById(id string) (user *domain.User, err error) {
	result, err := m.baseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := result.Decode(&user); err != nil {
		glog.Error(err)
		return nil, err
	}
	return user, nil
}

// 建立mongodb index
func (m *mgoUsersRepository) createIndices() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collections, _ := m.baseRepo.Collection.Database().ListCollectionNames(ctx, bson.D{})
	for _, collection := range collections {
		if collection == m.collectionName {
			m.baseRepo.Collection.Database().Collection(m.collectionName).Indexes().DropAll(ctx)
		}
	}

	models := []mongo.IndexModel{
		{
			Keys: bson.M{
				"modiTime": -1,
			},
			Options: &options.IndexOptions{
				Name:       tools.String("modiTimeIdx"),
				Background: tools.Bool(true),
			},
		},
		{
			Keys: bson.M{
				"nickname": 1,
			},
			Options: &options.IndexOptions{
				Name:       tools.String("nicknameIdx"),
				Background: tools.Bool(true),
				Unique:     tools.Bool(false),
			},
		},
	}
	if _, err := m.baseRepo.Collection.Indexes().CreateMany(ctx, models); err != nil {
		return err
	}
	return nil
}
