package services

import (
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
	if err := f.fillLastCommentID(post); err != nil {
		return err
	}

	return nil
}

func (f *PostFiller) fillLastCommentID(post *model.PostDto) error {
	if post.CommentsCount == 0 || post.LastCommentID != nil {
		return nil
	}

	lastComment, err := f.postCommentsGetter.GetLastPostComment(post.ID)
	if err != nil {
		return err
	}

	if lastComment != nil {
		post.LastCommentID = &lastComment.ID
	}

	return nil
}
