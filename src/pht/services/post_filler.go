package services

import (
	"github.com/samber/lo"
	"pht/comments-processor/pht/model"
)

type PostFiller struct {
	postCommentsGetter PostCommentsGetter
}

func NewPostFiller(postCommentsGetter PostCommentsGetter) *PostFiller {
	return &PostFiller{
		postCommentsGetter: postCommentsGetter,
	}
}

func (f *PostFiller) FillPost(post *model.PostDto) error {
	if err := f.fillLastCommentId(post); err != nil {
		return err
	}

	return nil
}

func (f *PostFiller) fillLastCommentId(post *model.PostDto) error {
	if post.CommentsCount == 0 || post.LastCommentId != nil {
		return nil
	}

	comments, err := f.postCommentsGetter.GetPostMostRecentComments(post.Id)
	if err != nil {
		return err
	}

	if v, ok := lo.First(comments); ok {
		post.LastCommentId = &v.Id
	}

	return nil
}
