package slack_test

import (
	"cd-slack-notification-bot/go/pkg/slack"
	"cd-slack-notification-bot/go/pkg/testutils"
	"testing"
	"time"

	slackgo "github.com/slack-go/slack"
	"github.com/stretchr/testify/require"
)

func Test_RetrieveChannelHistory(t *testing.T) {
	testutils.SkipCI(t)

	var messages []slackgo.Message

	t.Run("Retrive messages from a channel", func(t *testing.T) {
		client := slack.NewClientDefault()
		var err error
		messages, err = client.RetrieveChannelHistory(
			"C01Q7H30F3L",
			time.Now().Add(-time.Hour*48),
			time.Now(),
		)
		require.NoError(t, err)
		require.NotNil(t, messages)
		require.NotEmpty(t, messages)
	})

	t.Run("Check getting RunID from messages", func(t *testing.T) {
		for _, message := range messages {
			runID, err := slack.ParseRunIDFromCodeBuildMessage(message)
			require.NoError(t, err)

			t.Log(`
		Type: ` + message.Type + `
		SubType: ` + message.SubType + `
		Username: ` + message.Username + `
		User: ` + message.User + `
		BotID: ` + message.BotID + `
		RunID: ` + runID + `
		----------------------
		`)
		}
	})
}
