package testutils

import (
	"os"
	"testing"
)

/*
SkipCI will skip a test if run in CI environment
	- CI=true go test    : will skip the test
	- CI=false go test   : wil not skip the test
	- CI="" go test      : wil not skip the test
	- go test            : wil not skip the test
*/

func SkipCI(t *testing.T) {
	isInCI := os.Getenv("CI")
	if isInCI != "" && isInCI != "false" {
		t.Skip("Skipping test in CI environment")
	}
}
