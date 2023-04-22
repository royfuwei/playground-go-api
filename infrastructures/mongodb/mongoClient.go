package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoClient(addr string) (*mongo.Client, context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mgoAddr := fmt.Sprintf("mongodb://%s", addr)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mgoAddr))
	if err != nil {
		glog.Fatal(err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		glog.Fatal(err)
	}
	glog.Infof("[connected] Connected to mongo: %s", mgoAddr)
	return client, ctx
}
