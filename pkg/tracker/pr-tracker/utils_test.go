package prtracker_test

import (
	prtracker "cd-slack-notification-bot/go/pkg/tracker/pr-tracker"
	"cd-slack-notification-bot/go/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

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
