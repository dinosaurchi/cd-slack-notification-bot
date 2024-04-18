package config

import "cd-slack-notification-bot/go/pkg/utils"

type GithubConfig struct {
	APIURL    string
	Token     string
	RepoName  string
	RepoOwner string
}

func GetGithubConfig() *GithubConfig {
	return &GithubConfig{
		APIURL:    "https://api.github.com",
		Token:     utils.GetEnvVarValue("GITHUB_TOKEN", false),
		RepoName:  utils.GetEnvVarValue("GITHUB_REPO_NAME", false),
		RepoOwner: utils.GetEnvVarValue("GITHUB_REPO_OWNER", false),
	}
}
