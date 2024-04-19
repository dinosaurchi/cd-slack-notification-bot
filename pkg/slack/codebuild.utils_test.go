package slack_test

import (
	"cd-slack-notification-bot/go/pkg/slack"
	"cd-slack-notification-bot/go/pkg/utils"
	"testing"

	slackgo "github.com/slack-go/slack"
	"github.com/stretchr/testify/require"
)

func Test_GetAWSCodeSuiteRunIDFromMessage(t *testing.T) {
	testCases := []struct {
		message    string
		expectedID string
	}{
		{
			message:    "*\u003chttps://us-west-2.console.aws.amazon.com/codesuite/codebuild/projects/go-backend-cd/build/go-backend-cd:95b976f5-57e6-47f8-8d74-3eddbd2e7ec3/phase?region=us-west-2\u0026amp;referer_source=codestar-notifications\u0026amp;referer_medium=chatbot|:x: AWS CodeBuild Notification | us-west-2 | Account: 475910951137\u003e*",
			expectedID: "go-backend-cd:95b976f5-57e6-47f8-8d74-3eddbd2e7ec3",
		},
		{
			message:    "*\u003chttps://us-west-2.console.aws.amazon.com/codesuite/codebuild/projects/go-backend-cd/build/go-backend-cd:79fbf6fe-cc0e-4e69-a8ba-1290715507fd/phase?region=us-west-2\u0026amp;referer_source=codestar-notifications\u0026amp;referer_medium=chatbot|:x: AWS CodeBuild Notification | us-west-2 | Account: 475910951137\u003e*",
			expectedID: "go-backend-cd:79fbf6fe-cc0e-4e69-a8ba-1290715507fd",
		},
		{
			// Changed region to ap-southeast-1
			message:    "*\u003chttps://ap-southeast-1.console.aws.amazon.com/codesuite/codebuild/projects/go-backend-cd/build/go-backend-cd:95b976f5-57e6-47f8-8d74-3eddbd2e7ec3/phase?region=us-west-2\u0026amp;referer_source=codestar-notifications\u0026amp;referer_medium=chatbot|:x: AWS CodeBuild Notification | us-west-2 | Account: 475910951137\u003e*",
			expectedID: "go-backend-cd:95b976f5-57e6-47f8-8d74-3eddbd2e7ec3",
		},
		{
			// Changed codesuite -> coderun
			message:    "*\u003chttps://us-west-2.console.aws.amazon.com/coderun/codebuild/projects/go-backend-cd/build/go-backend-cd:95b976f5-57e6-47f8-8d74-3eddbd2e7ec3/phase?region=us-west-2\u0026amp;referer_source=codestar-notifications\u0026amp;referer_medium=chatbot|:x: AWS CodeBuild Notification | us-west-2 | Account: 475910951137\u003e*",
			expectedID: "go-backend-cd:95b976f5-57e6-47f8-8d74-3eddbd2e7ec3",
		},
		// Not found Run ID cases
		{
			// Changed to google.com
			message:    "*\u003chttps://us-west-2.console.aws.google.com/codesuite/codebuild/projects/go-backend-cd/build/go-backend-cd:95b976f5-57e6-47f8-8d74-3eddbd2e7ec3/phase?region=us-west-2\u0026amp;referer_source=codestar-notifications\u0026amp;referer_medium=chatbot|:x: AWS CodeBuild Notification | us-west-2 | Account: 475910951137\u003e*",
			expectedID: "",
		},
		{
			// Removed the build/
			message:    "*\u003chttps://us-west-2.console.aws.amazon.com/codesuite/codebuild/projects/go-backend-cd/go-backend-cd:95b976f5-57e6-47f8-8d74-3eddbd2e7ec3/phase?region=us-west-2\u0026amp;referer_source=codestar-notifications\u0026amp;referer_medium=chatbot|:x: AWS CodeBuild Notification | us-west-2 | Account: 475910951137\u003e*",
			expectedID: "",
		},
		{
			// Removed the /phase
			message:    "*\u003chttps://us-west-2.console.aws.amazon.com/codesuite/codebuild/projects/go-backend-cd/build/go-backend-cd:95b976f5-57e6-47f8-8d74-3eddbd2e7ec3/?region=us-west-2\u0026amp;referer_source=codestar-notifications\u0026amp;referer_medium=chatbot|:x: AWS CodeBuild Notification | us-west-2 | Account: 475910951137\u003e*",
			expectedID: "",
		},
		{
			// Replace : with / in the Run ID
			message:    "*\u003chttps://us-west-2.console.aws.amazon.com/codesuite/codebuild/projects/go-backend-cd/build/go-backend-cd/95b976f5-57e6-47f8-8d74-3eddbd2e7ec3/phase?region=us-west-2\u0026amp;referer_source=codestar-notifications\u0026amp;referer_medium=chatbot|:x: AWS CodeBuild Notification | us-west-2 | Account: 475910951137\u003e*",
			expectedID: "",
		},
		{
			// Empty Run ID
			message:    "*\u003chttps://us-west-2.console.aws.amazon.com/codesuite/codebuild/projects/go-backend-cd/build//phase?region=us-west-2\u0026amp;referer_source=codestar-notifications\u0026amp;referer_medium=chatbot|:x: AWS CodeBuild Notification | us-west-2 | Account: 475910951137\u003e*",
			expectedID: "",
		},
	}

	for _, tc := range testCases {
		t.Run("Case: "+tc.message, func(t *testing.T) {
			runID, err := slack.GetAWSCodeSuiteRunIDFromMessage(tc.message)
			require.NoError(t, err)
			require.Equal(t, tc.expectedID, runID)
		})
	}
}

func Test_GetAWSCodeBuildRunID(t *testing.T) {
	testCases := []struct {
		targetURL  string
		expectedID string
	}{
		{
			targetURL:  "https://us-west-2.console.aws.amazon.com/codebuild/home?region=us-west-2#/builds/go-backend-cd:fa393c46-5080-4582-9f5c-d4ced847fd63/view/new",
			expectedID: "go-backend-cd:fa393c46-5080-4582-9f5c-d4ced847fd63",
		},
		{
			targetURL:  "https://us-west-2.console.aws.amazon.com/codebuild/home?region=us-west-2#/builds/go-backend-cd:fa393c46-5080-4582-9f5c-d4ced847fd63/view",
			expectedID: "go-backend-cd:fa393c46-5080-4582-9f5c-d4ced847fd63",
		},
		{
			targetURL:  "https://us-west-2.console.aws.amazon.com/codebuild/home?region=us-west-2#/builds/go-backend-cd:c3cbbf62-b806-4f5a-b6a9-71fab1b6f11f/view/new",
			expectedID: "go-backend-cd:c3cbbf62-b806-4f5a-b6a9-71fab1b6f11f",
		},
		{
			targetURL:  "https://ap-southeast-1.console.aws.amazon.com/codebuild/home?region=ap-southeast-1#/builds/go-backend-cd:c3cbbf62-b806-4f5a-b6a9-71fab1b6f11f/view/new",
			expectedID: "go-backend-cd:c3cbbf62-b806-4f5a-b6a9-71fab1b6f11f",
		},
		{
			// Changed codebuild -> coderun
			targetURL:  "https://us-west-2.console.aws.amazon.com/coderun/home?region=us-west-2#/builds/go-backend-cd:fa393c46-5080-4582-9f5c-d4ced847fd63/view",
			expectedID: "go-backend-cd:fa393c46-5080-4582-9f5c-d4ced847fd63",
		},
		// Not found Run ID cases
		{
			// Changed to google.com
			targetURL:  "https://ap-southeast-1.console.aws.google.com/codebuild/home?region=ap-southeast-1#/builds/go-backend-cd:c3cbbf62-b806-4f5a-b6a9-71fab1b6f11f/view/new",
			expectedID: "",
		},
		{
			// Removed the /builds/
			targetURL:  "https://ap-southeast-1.console.aws.amazon.com/codebuild/home?region=ap-southeast-1#/go-backend-cd:c3cbbf62-b806-4f5a-b6a9-71fab1b6f11f/view/new",
			expectedID: "",
		},
		{
			// Removed the /view
			targetURL:  "https://ap-southeast-1.console.aws.amazon.com/codebuild/home?region=ap-southeast-1#/builds/go-backend-cd:c3cbbf62-b806-4f5a-b6a9-71fab1b6f11f",
			expectedID: "",
		},
		{
			// Replace : with / in the Run ID
			targetURL:  "https://ap-southeast-1.console.aws.amazon.com/codebuild/home?region=ap-southeast-1#/builds/go-backend-cd/c3cbbf62-b806-4f5a-b6a9-71fab1b6f11f/view/new",
			expectedID: "",
		},
		{
			// Empty Run ID
			targetURL:  "https://ap-southeast-1.console.aws.amazon.com/codebuild/home?region=ap-southeast-1#/builds//view/new",
			expectedID: "",
		},
	}

	for _, tc := range testCases {
		t.Run("Case: "+tc.targetURL, func(t *testing.T) {
			runID, err := slack.GetAWSCodeBuildRunID(tc.targetURL)
			require.NoError(t, err)
			require.Equal(t, tc.expectedID, runID)
		})
	}
}

func Test_ParseRunIDFromCodeBuildMessage(t *testing.T) {
	messages := []slackgo.Message{}
	err := utils.LoadFromFile("./responses/codebuild.json", &messages)
	require.NoError(t, err)
	require.Len(t, messages, 17)

	runIDs := []string{}
	for _, message := range messages {
		runID, err := slack.ParseRunIDFromCodeBuildMessage(message)
		require.NoError(t, err)
		runIDs = append(runIDs, runID)
	}
	require.Len(t, runIDs, len(messages))

	expectedRunIDs := []string{
		"go-backend-cd:f739b205-6483-47b5-a75a-75a6f0dbb2be",
		"go-backend-cd:95b976f5-57e6-47f8-8d74-3eddbd2e7ec3",
		"go-backend-cd:79fbf6fe-cc0e-4e69-a8ba-1290715507fd",
		"go-backend-cd:8776f54f-ce45-4229-99a9-aa62a0f750ee",
		"go-backend-cd:35eab56c-d475-4c5c-acee-b59ed94cbf06",
		"go-backend-cd:c9ac5d97-8c9e-4209-8015-e15e5e29ffee",
		"go-backend-cd:acf71aea-9890-475d-be4a-3aede7867304",
		"go-backend-cd:324eeac4-f25f-4dde-9fca-5af449dfd221",
		"go-backend-cd:254d25bc-3b75-4d0e-9aa9-f46ba7c79cdc",
		"go-backend-cd:a3112b4c-b6ca-4154-8537-185a8bc4610f",
		"go-backend-cd:48ae0708-1f3e-45d7-888b-fbc9196c696a",
		"go-backend-cd:75bfb7c4-5909-441b-b44c-07352143f465",
		"go-backend-cd:5637ebdb-9881-45ed-b746-49af9b03123d",
		"go-backend-cd:64c4b741-39f9-4056-831f-339dcff8e7b8",
		"go-backend-cd:ec461b53-829e-42fb-a3b0-9240b8a1fb0d",
		"go-backend-cd:f6bc20b5-057c-4832-a025-72c6e85cb43b",
		"",
	}

	require.Equal(t, expectedRunIDs, runIDs)
}
