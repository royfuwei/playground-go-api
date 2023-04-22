package mongodb

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ObjIds 字串轉ObjectID
func ObjIds(ids []string) ([]primitive.ObjectID, error) {
	var objIds []primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objIds = append(objIds, objID)
	}
	return objIds, nil
}

// IsMongoDupKey 檢查error是否為重複類型的錯誤
func IsMongoDupKey(err error) bool {
	wce, ok := err.(mongo.CommandError)
	if !ok {
		return false
	}
	return wce.Code == 11000 || wce.Code == 11001 || wce.Code == 12582 || wce.Code == 16460 && strings.Contains(wce.Message, " E11000 ")
}
