package config

import (
	"cd-slack-notification-bot/go/pkg/utils"
	"time"
)

type SlackConfig struct {
	Token                          string
	CodeBuildNotificationChannelID string
	RetrieveMessageBatchSize       int
	RetrieveMessageWaitDuration    time.Duration
}

func GetSlackConfig() *SlackConfig {
	const retrieveMessageBatchSize = 1000
	const retrieveMessageWaitDuration = 2 * time.Second
	return &SlackConfig{
		Token:                          utils.GetEnvVarValue("SLACK_TOKEN", false),
		CodeBuildNotificationChannelID: utils.GetEnvVarValue("SLACK_CODEBUILD_NOTIFICATION_CHANNEL_ID", false),
		RetrieveMessageBatchSize:       retrieveMessageBatchSize,
		RetrieveMessageWaitDuration:    retrieveMessageWaitDuration,
	}
}
