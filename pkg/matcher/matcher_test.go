package matcher_test

import (
	"cd-slack-notification-bot/go/pkg/matcher"
	cdtracker "cd-slack-notification-bot/go/pkg/tracker/cd-tracker"
	prtracker "cd-slack-notification-bot/go/pkg/tracker/pr-tracker"
	"cd-slack-notification-bot/go/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RunMatcher(t *testing.T) {
	state := &matcher.State{
		ResolvedRunIDs: map[string]*matcher.MatchedResult{},
	}

	// Load CD tracker state
	cdState := &cdtracker.State{}
	err := utils.LoadFromFile("samples/cd_tracker_state.json", cdState)
	require.NoError(t, err)

	// Load PR tracker state
	prState := &prtracker.State{}
	err = utils.LoadFromFile("samples/pr_tracker_state.json", prState)
	require.NoError(t, err)

	// Run matcher
	state, err = matcher.RunMatcher(state, cdState, prState)
	require.NoError(t, err)
	require.NotNil(t, state)

	// Load State 1
	expectedState1 := &matcher.State{}
	err = utils.LoadFromFile("samples/matcher_state_1.json", expectedState1)
	require.NoError(t, err)
	require.Equal(t, expectedState1, state)

	// Run matcher again
	state, err = matcher.RunMatcher(state, cdState, prState)
	require.NoError(t, err)
	require.NotNil(t, state)

	// Load State 2
	expectedState2 := &matcher.State{}
	err = utils.LoadFromFile("samples/matcher_state_2.json", expectedState2)
	require.NoError(t, err)
	require.Equal(t, expectedState2, state)

	// Add more CD info
	const newRunID = "go-backend-cd:4783432423432"
	cdState.RunIDToCDs[newRunID] = &cdtracker.CDInfo{
		ThreadTimestamp: "163234234.0001",
		RunID:           newRunID,
	}
	// Add more PR info
	const newThreadTimestamp = "163234666.05465"
	prState.PRs[newThreadTimestamp] = &prtracker.PRInfo{
		ThreadTimestamp: newThreadTimestamp,
		PRNumber:        123,
		Statuses: []prtracker.StatusInfo{
			{
				CodeBuildURL: "https://us-west-2.console.aws.amazon.com/codebuild/home?region=us-west-2#/builds/" + newRunID + "/view/new",
				State:        "success",
			},
		},
	}

	// Run matcher again
	state, err = matcher.RunMatcher(state, cdState, prState)
	require.NoError(t, err)
	require.NotNil(t, state)

	// Load State 3
	expectedState3 := &matcher.State{}
	err = utils.LoadFromFile("samples/matcher_state_3.json", expectedState3)
	require.NoError(t, err)
	require.Equal(t, expectedState3, state)
}
