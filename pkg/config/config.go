package config

type Config struct {
	Slack  *SlackConfig
	Github *GithubConfig
}

func GetConfigDefault() *Config {
	return &Config{
		Slack:  GetSlackConfig(),
		Github: GetGithubConfig(),
	}
}
