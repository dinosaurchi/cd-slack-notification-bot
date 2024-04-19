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

	// Set new CD and PR thread's timestamps
	matcherState = prepareTestThreads(t,
		slackClient,
		matcherState,
		githubPRNotificationChannelID,
		codeBuildNotificationChannelID,
	)

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

	// Add new RunID and its thread pair
	const newRunID = "go-backend-cd:534534543"
	matcherState.ResolvedRunIDs[newRunID] = &matcher.MatchedResult{
		CDThreadTimestamp: "", // will be set later
		PRThreadTimestamp: "", // will be set later
		Statuses:          []string{"success"},
		PRNumber:          444,
	}
	var prTS, cdTS string
	matcherState, prTS, cdTS = addNewThreadPair(t,
		slackClient,
		matcherState,
		newRunID,
		githubPRNotificationChannelID,
		codeBuildNotificationChannelID,
	)
	require.NotNil(t, matcherState)
	require.NotEmpty(t, prTS)
	require.NotEmpty(t, cdTS)

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

func prepareTestThreads(
	t *testing.T,
	slackClient slack.Client,
	matcherState *matcher.State,
	githubPRNotificationChannelID string,
	codeBuildNotificationChannelID string,
) *matcher.State {
	for runID := range matcherState.ResolvedRunIDs {
		var prThreadTimestamp, cdThreadTimestamp string
		matcherState, prThreadTimestamp, cdThreadTimestamp = addNewThreadPair(t,
			slackClient,
			matcherState,
			runID,
			githubPRNotificationChannelID,
			codeBuildNotificationChannelID,
		)

		require.NotNil(t, matcherState)
		require.NotEmpty(t, prThreadTimestamp)
		require.NotEmpty(t, cdThreadTimestamp)
	}

	return matcherState
}

func addNewThreadPair(
	t *testing.T,
	slackClient slack.Client,
	matcherState *matcher.State,
	runID string,
	githubPRNotificationChannelID string,
	codeBuildNotificationChannelID string,
) (*matcher.State, string, string) {
	// Prepare test threads on GitHub PR notification channel
	t.Logf("Preparing test threads for RunID %s", runID)
	message := "Test PR message with RunID " + runID
	prThreadTimestamp, err := slackClient.CreateThread(
		githubPRNotificationChannelID,
		message,
	)
	require.NoError(t, err)
	matcherState.ResolvedRunIDs[runID].PRThreadTimestamp = prThreadTimestamp

	// Prepare test threads on CodeBuild notification channel
	t.Logf("Preparing test threads for RunID %s", runID)
	message = "Test CD message with RunID " + runID
	cdThreadTimestamp, err := slackClient.CreateThread(
		codeBuildNotificationChannelID,
		message,
	)
	require.NoError(t, err)
	matcherState.ResolvedRunIDs[runID].CDThreadTimestamp = cdThreadTimestamp

	return matcherState, prThreadTimestamp, cdThreadTimestamp
}
