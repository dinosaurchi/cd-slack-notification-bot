package slack

import (
	"cd-slack-notification-bot/go/pkg/config"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type Client struct {
	token string
	api   *slack.Client
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		api:   slack.New(token),
	}
}

func NewClientDefault() *Client {
	return NewClient(
		config.GetConfigDefault().Slack.Token,
	)
}

func (c *Client) RetrieveChannelHistory(
	channelID string,
	from time.Time,
	to time.Time,
) ([]slack.Message, error) {
	if channelID == "" {
		return nil, errors.Errorf("channelID is empty")
	}
	if from.After(to) {
		return nil, errors.Errorf("from time %v is after to time %v", from, to)
	}

	api := slack.New(c.token)

	allMessages := []slack.Message{}

	hasMore := true
	cursor := ""
	for hasMore {
		res, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
			ChannelID:          channelID,
			Inclusive:          true,
			IncludeAllMetadata: true,
			Limit:              config.GetConfigDefault().Slack.RetrieveMessageBatchSize,
			Oldest:             strconv.FormatInt(from.Unix(), 10),
			Latest:             strconv.FormatInt(to.Unix(), 10),
			Cursor:             cursor,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if res.Error != "" {
			return nil, errors.Errorf("Slack API error: %s", res.Error)
		}

		allMessages = append(allMessages, res.Messages...)

		hasMore = res.HasMore
		if hasMore {
			cursor = res.ResponseMetaData.NextCursor
		}

		time.Sleep(config.GetConfigDefault().Slack.RetrieveMessageWaitDuration)
	}

	return allMessages, nil
}
