package mgo

import (
	"context"
	"playground-go-api/config"
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

type mgoRefreshTokensRepository struct {
	baseRepo              *mongodb.BaseRepository
	expireAfterSeconds    time.Duration
	mgoRefreshTokenTTLIdx string
	collectionName        string
}

func NewMgoRefreshTokensRepository(client *mongo.Client) domain.RefreshTokensRepository {
	collectionName := "refreshtokens"
	collection := client.Database("auth").Collection(collectionName)
	baseRepo := mongodb.NewBaseRepository(collection)
	expireAfterSeconds := time.Duration(config.Cfgs.RefreshTokenExpired) * time.Second

	m := &mgoRefreshTokensRepository{
		baseRepo:              baseRepo,
		expireAfterSeconds:    expireAfterSeconds,
		mgoRefreshTokenTTLIdx: "ttlIdx",
		collectionName:        collectionName,
	}
	if err := m.setRefreshTokenTTLIdx(); err != nil {
		glog.Fatal(err)
	}
	return m
}

func (m *mgoRefreshTokensRepository) setRefreshTokenTTLIdx() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collections, _ := m.baseRepo.Collection.Database().ListCollectionNames(ctx, bson.D{})
	for _, collection := range collections {
		if collection == m.collectionName {
			m.baseRepo.Collection.Database().Collection(m.collectionName).Indexes().DropOne(ctx, m.mgoRefreshTokenTTLIdx)
		}
	}

	model := mongo.IndexModel{
		Keys: bson.M{
			"createAt": 1,
		},
		Options: &options.IndexOptions{
			Name:               tools.String(m.mgoRefreshTokenTTLIdx),
			ExpireAfterSeconds: tools.IntToInt32(int(m.expireAfterSeconds.Seconds())),
		},
	}
	if _, err := m.baseRepo.Collection.Indexes().CreateOne(ctx, model); err != nil {
		return err
	}
	return nil

}

func (m *mgoRefreshTokensRepository) Add(data *domain.RefreshToken) (*domain.RefreshToken, error) {
	createAt := primitive.NewDateTimeFromTime(time.Now())
	data.CreateAt = createAt
	insertResult, err := m.baseRepo.Add(data, 5*time.Second)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	refreshToken := &domain.RefreshToken{
		ID:           insertResult.InsertedID.(primitive.ObjectID),
		UserId:       data.UserId,
		RefreshToken: data.RefreshToken,
		CreateAt:     createAt,
	}
	return refreshToken, nil
}

func (m *mgoRefreshTokensRepository) find(filter bson.M, opts ...*options.FindOptions) (refreshTokens []*domain.RefreshToken, total int64, err error) {
	cur, total, err := m.baseRepo.Find(filter, 5*time.Second, opts...)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var result *domain.RefreshToken
		if err := cur.Decode(&result); err != nil {
			return nil, 0, err
		}
		refreshTokens = append(refreshTokens, result)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return refreshTokens, total, nil
}

func (m *mgoRefreshTokensRepository) DeleteById(id string) (bool, error) {
	result, err := m.baseRepo.DeleteByID(id, 5*time.Second)
	if err != nil {
		return false, err
	}
	return int(result.DeletedCount) > 0, nil
}

func (m *mgoRefreshTokensRepository) FindById(id string) (refreshToken *domain.RefreshToken, err error) {
	result, err := m.baseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := result.Decode(&refreshToken); err != nil {
		glog.Error(err)
		return nil, err
	}
	return refreshToken, nil
}

func (m *mgoRefreshTokensRepository) FindByIds(ids ...string) (refreshTokens []*domain.RefreshToken, total int64, err error) {
	objIds, err := mongodb.ObjIds(ids)
	if err != nil {
		return nil, total, err
	}
	filterF := bson.M{
		"_id": bson.M{
			"$in": objIds,
		},
	}
	return m.find(filterF)
}

func (m *mgoRefreshTokensRepository) FindByUids(uids ...string) (refreshTokens []*domain.RefreshToken, total int64, err error) {
	objIds, err := mongodb.ObjIds(uids)
	if err != nil {
		return nil, total, err
	}
	filterF := bson.M{
		"userId": bson.M{
			"$in": objIds,
		},
	}
	return m.find(filterF)
}

func (m *mgoRefreshTokensRepository) FindOneByFilter(filter bson.M) (*domain.RefreshToken, error) {
	var refreshToken *domain.RefreshToken
	result, err := m.baseRepo.FindOneByFilter(filter)
	if err != nil {
		return nil, err
	}

	if err := result.Decode(&refreshToken); err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (m *mgoRefreshTokensRepository) DeleteByIds(ids ...string) (bool, int, error) {
	objIds, err := mongodb.ObjIds(ids)
	if err != nil {
		return false, 0, err
	}
	filter := bson.M{
		"_id": bson.M{
			"$in": objIds,
		},
	}
	results, err := m.baseRepo.DeleteByFilter(filter, 5*time.Second)
	if err != nil {
		return false, 0, err
	}
	return results.DeletedCount > 0, int(results.DeletedCount), nil
}

func (m *mgoRefreshTokensRepository) FindByUidsAndDeleteMany(uids ...string) (bool, int, error) {
	userObjIds, err := mongodb.ObjIds(uids)
	if err != nil {
		return false, 0, err
	}
	filterF := bson.M{
		"userId": bson.M{
			"$in": userObjIds,
		},
	}

	refreshTokens, total, err := m.find(filterF)
	if err != nil {
		return false, 0, err
	}
	if total <= 0 {
		return false, 0, nil
	}

	var objIds []primitive.ObjectID
	for _, refreshToken := range refreshTokens {
		objIds = append(objIds, refreshToken.ID)
	}
	filterD := bson.M{
		"_id": bson.M{
			"$in": objIds,
		},
	}
	results, err := m.baseRepo.DeleteByFilter(filterD, 5*time.Second)
	if err != nil {
		return false, 0, err
	}
	deleteCount := int(results.DeletedCount)
	return deleteCount > 0, deleteCount, nil
}
