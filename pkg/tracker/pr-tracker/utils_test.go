package prtracker_test

import (
	prtracker "cd-slack-notification-bot/go/pkg/tracker/pr-tracker"
	"cd-slack-notification-bot/go/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetPRsWithSuccessfulCD(t *testing.T) {
	state := &prtracker.State{}
	err := utils.LoadFromFile("./samples/state.json", state)
	require.NoError(t, err)

	expected := map[string]*prtracker.PRInfo{
		"1711751025.657609": {
			PRNumber: 587,
			Statuses: []prtracker.StatusInfo{
				{
					State:        "success",
					CodeBuildURL: "https://us-west-2.console.aws.amazon.com/codebuild/home?region=us-west-2#/builds/go-backend-cd:06e6ee2b-fe15-4293-801f-d9cd6aa48eb5/view/new",
				},
			},
			ThreadTimestamp: "1711751025.657609",
		},
		"1712023933.086279": {
			PRNumber: 594,
			Statuses: []prtracker.StatusInfo{
				{
					State:        "success",
					CodeBuildURL: "https://us-west-2.console.aws.amazon.com/codebuild/home?region=us-west-2#/builds/go-backend-cd:eed208a6-c33b-49ee-91f2-56bdec33b9e0/view/new",
				},
			},
			ThreadTimestamp: "1712023933.086279",
		},
	}

	successfulPRs, err := state.GetPRsWithSuccessfulCD()
	require.NoError(t, err)
	require.Equal(t, expected, successfulPRs)
}

func Test_GetRunIDFromPRInfo(t *testing.T) {
	prInfos := []*prtracker.PRInfo{}
	err := utils.LoadFromFile("samples/pr_infos.json", &prInfos)
	require.NoError(t, err)
	require.Len(t, prInfos, 3)

	runIDs := []string{}
	for _, prInfo := range prInfos {
		runID, err := prtracker.GetRunIDFromPRInfo(prInfo)
		require.NoError(t, err)
		runIDs = append(runIDs, runID)
	}
	require.Len(t, runIDs, len(prInfos))

	expected := []string{
		"go-backend-cd:acf71aea-9890-475d-be4a-3aede7867304",
		"",
		"go-backend-cd:8776f54f-ce45-4229-99a9-aa62a0f750ee",
	}
	require.Equal(t, expected, runIDs)
}
