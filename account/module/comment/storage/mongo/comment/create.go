package commentstore

import (
	"context"

	commentmodel "github.com/tronglv92/accounts/module/comment/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *mgoStore) CreateComment(ctx context.Context, data *commentmodel.CommentCreate) (*commentmodel.CommentCreate, error) {

	// collection := s.client.Database(common.DBMongoName).Collection(common.CommentsCollection)
	comment, err := s.collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	if oid, ok := comment.InsertedID.(primitive.ObjectID); ok {
		data.ID = oid
	}

	return data, nil
}
