package commentbiz

import (
	"context"

	commentmodel "github.com/tronglv92/accounts/module/comment/model"

	"github.com/tronglv92/accounts/common"
)

type ListCommentStorage interface {
	ListComments(
		ctx context.Context,
		paging *common.Paging) ([]commentmodel.Comment, error)
}

type listCommentBiz struct {
	listCommentStorage ListCommentStorage
}

func NewListCommentBiz(listCommentStorage ListCommentStorage) *listCommentBiz {
	return &listCommentBiz{
		listCommentStorage: listCommentStorage,
	}
}
func (business *listCommentBiz) ListComments(
	ctx context.Context,
	paging *common.Paging) ([]commentmodel.Comment, error) {
	results, err := business.listCommentStorage.ListComments(ctx, paging)

	if err != nil {
		return nil, common.ErrCannotListEntity(commentmodel.CollectionComment, err)
	}

	return results, nil
}
