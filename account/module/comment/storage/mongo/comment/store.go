package commentstore

import (
	"github.com/tronglv92/accounts/common"
	commentmodel "github.com/tronglv92/accounts/module/comment/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type mgoStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoStore(client *mongo.Client) *mgoStore {
	collection := client.Database(common.DBMongoName).Collection(commentmodel.CollectionComment)
	return &mgoStore{client: client, collection: collection}
}
