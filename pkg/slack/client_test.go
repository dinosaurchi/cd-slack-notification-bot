package slack_test

import (
	"cd-slack-notification-bot/go/pkg/slack"
	"cd-slack-notification-bot/go/pkg/testutils"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_RetrieveChannelHistory(t *testing.T) {
	testutils.SkipCI(t)

	client := slack.NewClientDefault()
	messages, err := client.RetrieveChannelHistory(
		"C01Q7H30F3L",
		time.Now().Add(-time.Hour*48),
		time.Now(),
	)
	require.NoError(t, err)
	require.NotNil(t, messages)
	require.NotEmpty(t, messages)

	for _, message := range messages {
		runID, err := slack.ParseRunIDFromMessage(message)
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
}
