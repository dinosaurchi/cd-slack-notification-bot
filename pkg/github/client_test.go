package github_test

import (
	"cd-slack-notification-bot/go/pkg/github"
	"cd-slack-notification-bot/go/pkg/testutils"
	"cd-slack-notification-bot/go/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetPullRequestCommits(t *testing.T) {
	testutils.SkipCI(t)

	client := github.NewClientDefault()
	commits, err := client.GetPullRequestCommits(588)
	require.NoError(t, err)
	require.NotEmpty(t, commits)

	lastCommit := commits[len(commits)-1]

	t.Log("Author commit info:", utils.ToJSONStringPanic(lastCommit.GetAuthorCommitInfo()))
	t.Log("Committer commit info:", utils.ToJSONStringPanic(lastCommit.GetCommitterCommitInfo()))
	t.Log("Message:", lastCommit.GetMessage())
	t.Log("Commit SHA:", lastCommit.GetCommitSHA())
}

func Test_GetPullRequestInfo(t *testing.T) {
	testutils.SkipCI(t)

	client := github.NewClientDefault()
	pr, err := client.GetPullRequestInfo(588)
	require.NoError(t, err)
	require.NotNil(t, pr)

	t.Log("PR:", utils.ToJSONStringPanic(pr))
}

func Test_GetCommitCDInfo(t *testing.T) {
	testutils.SkipCI(t)

	client := github.NewClientDefault()
	cdInfo, err := client.GetCommitCDInfo("0d0b012142a9a3b8053d64646b68ec75ac93c496")
	require.NoError(t, err)
	require.NotNil(t, cdInfo)

	t.Log("CD info:", utils.ToJSONStringPanic(cdInfo))
}

func Test_GetPullRequestCDInfo(t *testing.T) {
	testutils.SkipCI(t)

	client := github.NewClientDefault()
	cdInfo, err := client.GetPullRequestCDInfo(588)
	require.NoError(t, err)
	require.NotNil(t, cdInfo)

	t.Log("CD info:", utils.ToJSONStringPanic(cdInfo))
}
