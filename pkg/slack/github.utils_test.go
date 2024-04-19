package slack_test

import (
	"cd-slack-notification-bot/go/pkg/slack"
	"cd-slack-notification-bot/go/pkg/utils"
	"testing"

	slackgo "github.com/slack-go/slack"
	"github.com/stretchr/testify/require"
)

func Test_ParseGithubPRInfoFromTitle(t *testing.T) {
	testCases := []struct {
		title    string
		expected *slack.GithubPRSubInfo
	}{
		{
			title: "\u003chttps://github.com/stablyio/stably-ramp-frontend/pull/420|#420 Improvement/kyc tos post object id to backend \u0026amp; refactor e2e tests\u003e",
			expected: &slack.GithubPRSubInfo{
				RepoName:  "stably-ramp-frontend",
				RepoOwner: "stablyio",
				PRNumber:  420,
			},
		},
		{
			title: "\u003chttps://github.com/stablyio/terraform-trinity/pull/6|#6 Aave - Create S3 bucket to store webapp\u003e",
			expected: &slack.GithubPRSubInfo{
				RepoName:  "terraform-trinity",
				RepoOwner: "stablyio",
				PRNumber:  6,
			},
		},
		{
			title: "\u003chttps://github.com/stablyio/go-backend/pull/614|#614 Fix Horizen EON's deposit issue\u003e",
			expected: &slack.GithubPRSubInfo{
				RepoName:  "go-backend",
				RepoOwner: "stablyio",
				PRNumber:  614,
			},
		},
	}

	for _, tc := range testCases {
		t.Run("Case: "+tc.title, func(t *testing.T) {
			actual, err := slack.ParseGithubPRSubInfoFromTitle(tc.title)
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Equal(t, *tc.expected, *actual)
		})
	}

	t.Run("Case: Invalid title", func(t *testing.T) {
		title := "Invalid title"
		actual, err := slack.ParseGithubPRSubInfoFromTitle(title)
		require.NoError(t, err)
		require.Nil(t, actual)

		title = ""
		actual, err = slack.ParseGithubPRSubInfoFromTitle(title)
		require.NoError(t, err)
		require.Nil(t, actual)
	})
}

func Test_ParseGithubPRInfoFromPROpenedMessage(t *testing.T) {
	messages := []slackgo.Message{}
	err := utils.LoadFromFile("./responses/github.json", &messages)
	require.NoError(t, err)
	require.Len(t, messages, 7)

	githubInfos := []*slack.GithubPRInfo{}
	for _, message := range messages {
		runID, err := slack.ParseGithubPRInfoFromPROpenedMessage(message)
		require.NoError(t, err)
		githubInfos = append(githubInfos, runID)
	}
	require.Len(t, githubInfos, len(messages))

	expectedGithubInfos := []*slack.GithubPRInfo{
		{
			PRNumber:        615,
			RepoName:        "go-backend",
			RepoOwner:       "stablyio",
			ThreadTimestamp: "1713476246.235249",
			Timestamp:       "1713476246.235249",
		},
		{
			PRNumber:        420,
			RepoName:        "stably-ramp-frontend",
			RepoOwner:       "stablyio",
			ThreadTimestamp: "1713454012.666599",
			Timestamp:       "1713454012.666599",
		},
		{
			PRNumber:        6,
			RepoName:        "terraform-trinity",
			RepoOwner:       "stablyio",
			ThreadTimestamp: "1713335235.568829",
			Timestamp:       "1713335235.568829",
		},
		{
			PRNumber:        432,
			RepoName:        "stably-ramp-frontend",
			RepoOwner:       "stablyio",
			ThreadTimestamp: "1713320791.744569",
			Timestamp:       "1713320791.744569",
		},
		nil,
		nil,
		nil,
	}

	require.Equal(t, expectedGithubInfos, githubInfos)
}
