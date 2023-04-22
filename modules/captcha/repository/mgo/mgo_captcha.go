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

type mgoCaptchaRepository struct {
	baseRepo         *mongodb.BaseRepository
	collectionName   string
	mgoCaptchaTTLIdx string
}

func NewMgoCaptchaRepository(client *mongo.Client) domain.CaptchaRepository {
	collectionName := "captcha"
	collection := client.Database("auth").Collection((collectionName))
	baseRepo := mongodb.NewBaseRepository(collection)
	mgoCaptchaTTLIdx := "ttlIdx"
	m := &mgoCaptchaRepository{
		baseRepo,
		collectionName,
		mgoCaptchaTTLIdx,
	}
	if err := m.setCaptchaTTLIdx(); err != nil {
		glog.Fatal(err)
	}
	return m
}

func (m *mgoCaptchaRepository) setCaptchaTTLIdx() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collections, _ := m.baseRepo.Collection.Database().ListCollectionNames(ctx, bson.D{})
	for _, collection := range collections {
		if collection == m.collectionName {
			m.baseRepo.Collection.Database().Collection(m.collectionName).Indexes().DropOne(ctx, m.mgoCaptchaTTLIdx)
		}
	}

	model := mongo.IndexModel{
		Keys: bson.M{
			"expiryDate": 1,
		},
		Options: &options.IndexOptions{
			Name:               tools.String(m.mgoCaptchaTTLIdx),
			ExpireAfterSeconds: tools.IntToInt32(0),
		},
	}
	if _, err := m.baseRepo.Collection.Indexes().CreateOne(ctx, model); err != nil {
		return err
	}
	return nil
}

func (m *mgoCaptchaRepository) AddBaseCaptcha(expiresTime time.Duration) (*domain.BaseCaptcha, error) {
	expiryDate := time.Now().Add(expiresTime)
	data := &domain.BaseCaptcha{}
	data.ExpiryDate = primitive.NewDateTimeFromTime(expiryDate)
	insertResult, err := m.baseRepo.Add(data, 5*time.Second)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	baseCaptcha := data
	baseCaptcha.ID = insertResult.InsertedID.(primitive.ObjectID)
	return baseCaptcha, nil
}
func (m *mgoCaptchaRepository) AddSmsCaptcha(data *domain.SmsCaptcha, expiresTime time.Duration) (*domain.SmsCaptcha, error) {
	expiryDate := time.Now().Add(expiresTime)
	data.ExpiryDate = primitive.NewDateTimeFromTime(expiryDate)
	insertResult, err := m.baseRepo.Add(data, 5*time.Second)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	smsCaptcha := data
	smsCaptcha.ID = insertResult.InsertedID.(primitive.ObjectID)
	return smsCaptcha, nil
}

func (m *mgoCaptchaRepository) FindBaseCaptchaById(id string) (baseCaptcha *domain.BaseCaptcha, err error) {
	result, err := m.baseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := result.Decode(&baseCaptcha); err != nil {
		glog.Error(err)
		return nil, err
	}
	return baseCaptcha, nil
}

func (m *mgoCaptchaRepository) FindSmsCaptchaById(id string) (smsCaptcha *domain.SmsCaptcha, err error) {
	result, err := m.baseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := result.Decode(&smsCaptcha); err != nil {
		glog.Error(err)
		return nil, err
	}
	return smsCaptcha, nil
}

func (m *mgoCaptchaRepository) DeleteById(id string) (bool, error) {
	result, err := m.baseRepo.DeleteByID(id, 5*time.Second)
	if err != nil {
		return false, err
	}
	return int(result.DeletedCount) > 0, nil
}

func (m *mgoCaptchaRepository) FindOneSmsCaptcha(captcha, identifier, telephone, telephoneRegion string) (*domain.SmsCaptcha, error) {
	var smsCaptcha *domain.SmsCaptcha
	filter := bson.M{
		"captcha":         captcha,
		"identifier":      identifier,
		"telephone":       telephone,
		"telephoneRegion": telephoneRegion,
	}
	result, err := m.baseRepo.FindOneByFilter(filter)
	if err != nil {
		return nil, err
	}

	if err := result.Decode(&smsCaptcha); err != nil {
		return nil, err
	}
	return smsCaptcha, nil
}
