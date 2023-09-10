package commentstore

import (
	"context"
	"fmt"

	"github.com/tronglv92/accounts/common"
	commentmodel "github.com/tronglv92/accounts/module/comment/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *mgoStore) ListComments(ctx context.Context,
	paging *common.Paging) ([]commentmodel.Comment, error) {

	filter := bson.D{}
	fmt.Println("filter=", filter)

	// sort := bson.D{{Key: "full_slug", Value: 1}}
	// opts := options.Find().SetSort(sort)

	cursor, err := s.collection.Find(ctx, filter)

	if err != nil {

		return nil, err
	}
	var results []commentmodel.Comment
	if err = cursor.All(context.TODO(), &results); err != nil {

		return nil, err
	}
	fmt.Println("len(results)=", len(results))
	// collection.FindOne()
	return results, nil
}
