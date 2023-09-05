package commentbiz

import (
	"context"

	"github.com/tronglv92/accounts/common"
	commentmodel "github.com/tronglv92/accounts/module/comment/model"
)

type CreateCommentStorage interface {
	CreateComment(ctx context.Context, data *commentmodel.CommentCreate) (*commentmodel.CommentCreate, error)
}

type createCommentBiz struct {
	createCommentStorage CreateCommentStorage
}

func NewCreateCommentBiz(createCommentStorage CreateCommentStorage) *createCommentBiz {
	return &createCommentBiz{
		createCommentStorage: createCommentStorage,
	}
}
func (business *createCommentBiz) CreateComment(ctx context.Context, data *commentmodel.CommentCreate) (*commentmodel.CommentCreate, error) {

	data.Fullfill()

	//4. create new comment
	newComment, err := business.createCommentStorage.CreateComment(ctx, data)

	if err != nil {
		return nil, common.ErrCannotCreateEntity(commentmodel.CollectionComment, err)
	}

	return newComment, nil
}
