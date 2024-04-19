package slack

import (
	"cd-slack-notification-bot/go/pkg/config"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type Client interface {
	RetrieveChannelHistory(channelID string, from time.Time, to time.Time) ([]slack.Message, error)
	GetMessageLink(channelID string, timestamp string) (string, error)
	ReplyThread(channelID string, threadTimestamp string, message string) (string, error)
}

type clientImplementation struct {
	token string
	api   *slack.Client
}

func NewClient(token string) Client {
	return &clientImplementation{
		token: token,
		api:   slack.New(token),
	}
}

func NewClientDefault() Client {
	return NewClient(
		config.GetConfigDefault().Slack.Token,
	)
}

func (c *clientImplementation) RetrieveChannelHistory(
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

func (c *clientImplementation) GetMessageLink(
	channelID string,
	timestamp string,
) (string, error) {
	api := slack.New(c.token)
	link, err := api.GetPermalink(&slack.PermalinkParameters{
		Channel: channelID,
		Ts:      timestamp,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}
	return link, nil
}

func (c *clientImplementation) ReplyThread(
	channelID string,
	threadTimestamp string,
	text string,
) (string, error) {
	api := slack.New(c.token)
	_, outputTimestamp, _, err := api.SendMessage(
		channelID,
		slack.MsgOptionText(text, false),
		slack.MsgOptionTS(threadTimestamp),
	)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return outputTimestamp, nil
}
