package strategies

import (
	"net/url"
	"pht/comments-processor/pht/config"
	"pht/comments-processor/pht/model/dto"
	"pht/comments-processor/pht/services"
	"strconv"
	"strings"
)

type CheckPostStrategy interface {
	CheckPost(post dto.TablePost) (dto.CheckPostResult, error)
}

type ContentCheckPostStrategy struct {
	postsProvider *services.PostsProvider
	contentURL    string
}

func NewContentCheckPostStrategy(postsProvider *services.PostsProvider, config config.ConfigProvider) (*ContentCheckPostStrategy, error) {
	return &ContentCheckPostStrategy{
		postsProvider: postsProvider,
		contentURL:    config.ContentURL(),
	}, nil
}

func (s *ContentCheckPostStrategy) CheckPost(post dto.TablePost) (dto.CheckPostResult, error) {
	postURL, err := url.Parse(s.contentURL)
	if err != nil {
		return dto.CheckPostResult{}, err
	}

	postContent, err := s.postsProvider.GetPost(post.ID)
	if err != nil {
		return dto.CheckPostResult{}, err
	}

	oldCount := post.CommentsCount

	newCommentsCount := postContent.CommentsCount
	var newCommentsCountVal any = newCommentsCount

	query := url.Values{
		"ordering": {"-created_at"},
	}

	if newCommentsCount > 0 {
		if postContent.LastCommentID != nil {
			query.Set("comment_target_id", strconv.Itoa(*postContent.LastCommentID))
		}
	} else if postContent.DisableComments {
		newCommentsCountVal = "н/к"
	}

	postURL = postURL.JoinPath("#", "publicate", strconv.Itoa(post.ID))
	postURL.RawQuery = query.Encode()

	return dto.CheckPostResult{
		Title:            postContent.Title,
		OldCommentsCount: oldCount,
		NewCommentsCount: newCommentsCountVal,
		URL:              strings.Replace(postURL.String(), `/%23/`, `/#/`, 1),
	}, nil
}
