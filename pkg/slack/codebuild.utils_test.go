package slack_test

import (
	"cd-slack-notification-bot/go/pkg/slack"
	"testing"

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
