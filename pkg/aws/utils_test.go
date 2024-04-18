package aws_test

import (
	"cd-slack-notification-bot/go/pkg/aws"
	"testing"

	"github.com/stretchr/testify/require"
)

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
	}

	for _, tc := range testCases {
		t.Run("Case: "+tc.targetURL, func(t *testing.T) {
			runID, err := aws.GetAWSCodeBuildRunID(tc.targetURL)
			require.NoError(t, err)
			require.Equal(t, tc.expectedID, runID)
		})
	}
}
