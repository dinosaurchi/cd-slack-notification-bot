package notifier_test

import (
	"cd-slack-notification-bot/go/pkg/matcher"
	"cd-slack-notification-bot/go/pkg/notifier"
	"cd-slack-notification-bot/go/pkg/slack"
	"cd-slack-notification-bot/go/pkg/testutils"
	"cd-slack-notification-bot/go/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RunNotifierCustom(t *testing.T) {
	testutils.SkipCI(t)

	// dex-mm-integration-testing
	const githubPRNotificationChannelID = "C020Z9BBRR6"
	// tech-errors-test
	const codeBuildNotificationChannelID = "C02ELHM809H"

	state := &notifier.State{
		CDNotified: map[string]bool{},
		PRNotified: map[string]bool{},
	}

	slackClient := slack.NewClientDefault()

	// Load matcher state
	matcherState := &matcher.State{}
	err := utils.LoadFromFile("samples/matcher_state.json", matcherState)
	require.NoError(t, err)
	require.NotNil(t, matcherState)

	// Run notifier
	state, err = notifier.RunNotifierCustom(
		state,
		matcherState,
		githubPRNotificationChannelID,
		codeBuildNotificationChannelID,
		slackClient,
	)
	require.NoError(t, err)
	require.Equal(t, &notifier.State{
		CDNotified: map[string]bool{
			"go-backend-cd:35eab56c-d475-4c5c-acee-b59ed94cbf06": true,
			"go-backend-cd:48ae0708-1f3e-45d7-888b-fbc9196c696a": true,
			"go-backend-cd:75bfb7c4-5909-441b-b44c-07352143f465": true,
		},
		PRNotified: map[string]bool{
			"go-backend-cd:35eab56c-d475-4c5c-acee-b59ed94cbf06": true,
			"go-backend-cd:48ae0708-1f3e-45d7-888b-fbc9196c696a": true,
			"go-backend-cd:75bfb7c4-5909-441b-b44c-07352143f465": true,
		},
	}, state)

	// Modifer matcher state
	matcherState.ResolvedRunIDs["go-backend-cd:534534543"] = &matcher.MatchedResult{
		CDThreadTimestamp: "1713520757.407589",
		PRThreadTimestamp: "1713520753.202969",
		PRNumber:          444,
		Statuses:          []string{"success"},
	}

	// Run notifier again
	state, err = notifier.RunNotifierCustom(
		state,
		matcherState,
		githubPRNotificationChannelID,
		codeBuildNotificationChannelID,
		slackClient,
	)
	require.NoError(t, err)
	require.Equal(t, &notifier.State{
		CDNotified: map[string]bool{
			"go-backend-cd:35eab56c-d475-4c5c-acee-b59ed94cbf06": true,
			"go-backend-cd:48ae0708-1f3e-45d7-888b-fbc9196c696a": true,
			"go-backend-cd:75bfb7c4-5909-441b-b44c-07352143f465": true,
			"go-backend-cd:534534543":                            true,
		},
		PRNotified: map[string]bool{
			"go-backend-cd:35eab56c-d475-4c5c-acee-b59ed94cbf06": true,
			"go-backend-cd:48ae0708-1f3e-45d7-888b-fbc9196c696a": true,
			"go-backend-cd:75bfb7c4-5909-441b-b44c-07352143f465": true,
			"go-backend-cd:534534543":                            true,
		},
	}, state)
}
