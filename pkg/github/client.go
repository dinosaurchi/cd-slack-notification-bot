package github

import (
	"cd-slack-notification-bot/go/pkg/config"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type Client struct {
	repoOwner   string
	repoName    string
	githubToken string
}

func NewClient(repoOwner, repoName, githubToken string) *Client {
	return &Client{
		repoOwner:   repoOwner,
		repoName:    repoName,
		githubToken: githubToken,
	}
}

// Perform a GET request to the target URL
func (c *Client) getRequest(
	targetURL string,
	payload io.Reader,
) ([]byte, error) {
	req, err := http.NewRequest("GET", targetURL, payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.githubToken)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return body, nil
}

func NewClientDefault() *Client {
	return NewClient(
		config.GetConfigDefault().Github.RepoOwner,
		config.GetConfigDefault().Github.RepoName,
		config.GetConfigDefault().Github.Token,
	)
}

func (c *Client) GetPullRequestInfo(
	pullRequestNumber uint64,
) (*PullRequest, error) {
	targetURL, err := url.JoinPath(
		APIURL,
		"repos",
		c.repoOwner,
		c.repoName,
		"pulls",
		strconv.FormatUint(pullRequestNumber, 10),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := c.getRequest(targetURL, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Convert body to PullRequest struct
	var prInfo *PullRequest
	err = json.Unmarshal(body, &prInfo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return prInfo, nil
}

func (c *Client) GetPullRequestCommits(
	pullRequestNumber uint64,
) ([]*Commit, error) {
	// Example: https://api.github.com/repos/owner/myrepo/pulls/123/commits
	targetURL, err := url.JoinPath(
		APIURL,
		"repos",
		c.repoOwner,
		c.repoName,
		"pulls",
		strconv.FormatUint(pullRequestNumber, 10),
		"commits",
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := c.getRequest(targetURL, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Convert body to list of Commit structs
	var commits []*Commit
	err = json.Unmarshal(body, &commits)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return commits, nil
}

// Assume that the CD was configured to be triggered by a commit
// and run on AWS CodeBuild
func (c *Client) GetCommitCDInfo(
	commitSHA string,
) (*CDInfo, error) {
	targetURL, err := url.JoinPath(
		APIURL,
		"repos",
		c.repoOwner,
		c.repoName,
		"commits",
		commitSHA,
		"status",
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := c.getRequest(targetURL, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Convert body to CDInfo struct
	var cdInfo *CDInfo
	err = json.Unmarshal(body, &cdInfo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cdInfo, nil
}

func (c *Client) GetPullRequestCDInfo(
	pullRequestNumber uint64,
) (*CDInfo, error) {
	commits, err := c.GetPullRequestCommits(pullRequestNumber)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(commits) == 0 {
		return nil, errors.Errorf("No commits found for PR %d", pullRequestNumber)
	}

	// Get CD info of the last commit
	lastCommit := commits[len(commits)-1]
	return c.GetCommitCDInfo(lastCommit.GetCommitSHA())
}
