package slack_test

import (
	"cd-slack-notification-bot/go/pkg/slack"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseGithubPRInfoFromTitle(t *testing.T) {
	testCases := []struct {
		title    string
		expected *slack.GithubPRInfo
	}{
		{
			title: "\u003chttps://github.com/stablyio/stably-ramp-frontend/pull/420|#420 Improvement/kyc tos post object id to backend \u0026amp; refactor e2e tests\u003e",
			expected: &slack.GithubPRInfo{
				RepoName:  "stably-ramp-frontend",
				RepoOwner: "stablyio",
				PRNumber:  420,
			},
		},
		{
			title: "\u003chttps://github.com/stablyio/terraform-trinity/pull/6|#6 Aave - Create S3 bucket to store webapp\u003e",
			expected: &slack.GithubPRInfo{
				RepoName:  "terraform-trinity",
				RepoOwner: "stablyio",
				PRNumber:  6,
			},
		},
		{
			title: "\u003chttps://github.com/stablyio/go-backend/pull/614|#614 Fix Horizen EON's deposit issue\u003e",
			expected: &slack.GithubPRInfo{
				RepoName:  "go-backend",
				RepoOwner: "stablyio",
				PRNumber:  614,
			},
		},
	}

	for _, tc := range testCases {
		t.Run("Case: "+tc.title, func(t *testing.T) {
			actual, err := slack.ParseGithubPRInfoFromTitle(tc.title)
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Equal(t, *tc.expected, *actual)
		})
	}

	t.Run("Case: Invalid title", func(t *testing.T) {
		title := "Invalid title"
		actual, err := slack.ParseGithubPRInfoFromTitle(title)
		require.NoError(t, err)
		require.Nil(t, actual)

		title = ""
		actual, err = slack.ParseGithubPRInfoFromTitle(title)
		require.NoError(t, err)
		require.Nil(t, actual)
	})
}
