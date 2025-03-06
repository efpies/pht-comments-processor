package services

import (
	"github.com/samber/lo"
	"pht/comments-processor/pht/model/dto"
)

type PostFiller struct {
	postCommentsGetter PostCommentsGetter
}

func NewPostFiller(postCommentsGetter PostCommentsGetter) *PostFiller {
	return &PostFiller{
		postCommentsGetter: postCommentsGetter,
	}
}

func (f *PostFiller) FillPost(post *dto.Post) error {
	if err := f.fillLastCommentID(post); err != nil {
		return err
	}

	return nil
}

func (f *PostFiller) fillLastCommentID(post *dto.Post) error {
	if post.CommentsCount == 0 || post.LastCommentID != nil {
		return nil
	}

	comments, err := f.postCommentsGetter.GetPostMostRecentComments(post.ID)
	if err != nil {
		return err
	}

	if v, ok := lo.First(comments); ok {
		post.LastCommentID = &v.ID
	}

	return nil
}
