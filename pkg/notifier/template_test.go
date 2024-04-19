package notifier_test

import (
	"cd-slack-notification-bot/go/pkg/notifier"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetCDMessage(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		message, err := notifier.GetCDMessage(
			"https://dsadas.com/dadasdsadas",
			[]string{
				"success",
			},
		)
		require.NoError(t, err)
		require.Equal(t, "CD succeeded - https://dsadas.com/dadasdsadas", message)
	})

	t.Run("Failure case", func(t *testing.T) {
		message, err := notifier.GetCDMessage(
			"https://dsadas.com/dadasdsadas",
			[]string{
				"failure",
			},
		)
		require.NoError(t, err)
		require.Equal(t, "CD failed - https://dsadas.com/dadasdsadas", message)
	})
}
