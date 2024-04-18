package github_test

import (
	"cd-slack-notification-bot/go/pkg/github"
	"cd-slack-notification-bot/go/pkg/testutils"
	"cd-slack-notification-bot/go/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Client(t *testing.T) {
	t.Run("GetPullRequestCommits", func(t *testing.T) {
		testutils.SkipCI(t)

		// Test the GetPullRequestCommits method
		client := github.NewClientDefault()
		commits, err := client.GetPullRequestCommits(588)
		require.NoError(t, err)
		require.NotEmpty(t, commits)

		lastCommit := commits[len(commits)-1]

		t.Log("Author commit info:", utils.ToJSONStringPanic(lastCommit.GetAuthorCommitInfo()))
		t.Log("Committer commit info:", utils.ToJSONStringPanic(lastCommit.GetCommitterCommitInfo()))
		t.Log("Message:", lastCommit.GetMessage())
		t.Log("Commit SHA:", lastCommit.GetCommitSHA())
	})
}
