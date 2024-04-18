package slack

import (
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type GithubPRInfo struct {
	PRNumber        uint64
	RepoName        string
	RepoOwner       string
	ThreadTimestamp string
	Timestamp       string
}

// Return GitPRInfo struct from the message
// If the message is not PR opened message, return nil
func ParseGithubPRInfoFromPROpenedMessage(
	message slack.Message,
) (*GithubPRInfo, error) {
	for _, attachment := range message.Attachments {
		if attachment.CallbackID == "pr-opened-interaction" {
			res, err := ParseGithubPRSubInfoFromTitle(attachment.Title)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return &GithubPRInfo{
				PRNumber:        res.PRNumber,
				RepoName:        res.RepoName,
				RepoOwner:       res.RepoOwner,
				ThreadTimestamp: message.ThreadTimestamp,
				Timestamp:       message.Timestamp,
			}, nil
		}
	}

	var nilError error
	return nil, nilError
}

type GithubPRSubInfo struct {
	PRNumber  uint64
	RepoName  string
	RepoOwner string
}

func ParseGithubPRSubInfoFromTitle(
	title string,
) (*GithubPRSubInfo, error) {
	// [^/]+ means any character except '/' and '+' means one or more
	re := regexp.MustCompile(`github\.com/([^/\\*+]+)/([^/\\*+]+)/pull/([0-9]+)`)
	match := re.FindStringSubmatch(title)
	const minMatchLength = 4
	if len(match) >= minMatchLength {
		prNumber, err := strconv.ParseUint(match[3], 10, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return &GithubPRSubInfo{
			PRNumber:  prNumber,
			RepoName:  match[2],
			RepoOwner: match[1],
		}, nil
	}

	var nilError error
	return nil, nilError
}
