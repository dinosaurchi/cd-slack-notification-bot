package github

type Commit struct {
	SHA       string               `json:"sha"`
	NodeID    string               `json:"node_id"`
	Commit    *commitDetails       `json:"commit"`
	Author    *commitUserLoginInfo `json:"author,omitempty"`
	Committer *commitUserLoginInfo `json:"committer,omitempty"`
}

type commitDetails struct {
	Author       *commitUser        `json:"author"`
	Committer    *commitUser        `json:"committer"`
	Message      string             `json:"message"`
	CommentCount int                `json:"comment_count"`
	Verification commitVerification `json:"verification"`
}

type commitUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type commitUserLoginInfo struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

type commitVerification struct {
	Verified  bool   `json:"verified"`
	Reason    string `json:"reason"`
	Signature string `json:"signature,omitempty"`
	Payload   string `json:"payload,omitempty"`
}

type UserCommitInfo struct {
	Login   string `json:"login"`
	LoginID int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Date    string `json:"date"`
}

func (c *Commit) GetAuthorCommitInfo() *UserCommitInfo {
	return &UserCommitInfo{
		Login:   c.Author.Login,
		LoginID: c.Author.ID,
		Name:    c.Commit.Author.Name,
		Email:   c.Commit.Author.Email,
		Date:    c.Commit.Author.Date,
	}
}

func (c *Commit) GetCommitterCommitInfo() *UserCommitInfo {
	return &UserCommitInfo{
		Login:   c.Committer.Login,
		LoginID: c.Committer.ID,
		Name:    c.Commit.Committer.Name,
		Email:   c.Commit.Committer.Email,
		Date:    c.Commit.Committer.Date,
	}
}

func (c *Commit) GetMessage() string {
	return c.Commit.Message
}

func (c *Commit) GetCommitSHA() string {
	return c.SHA
}
