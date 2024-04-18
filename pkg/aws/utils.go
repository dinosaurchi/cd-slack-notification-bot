package aws

import (
	"regexp"
)

// If the targetURL is a valid AWS CodeBuild URL, this function will return the run ID.
// Otherwise, it will return an empty string.
func GetAWSCodeBuildRunID(
	targetURL string,
) (string, error) {
	re := regexp.MustCompile(`amazon\.com/.*/builds/(.+?)/view`)
	match := re.FindStringSubmatch(targetURL)
	if len(match) > 1 {
		return match[1], nil
	}
	return "", nil
}
