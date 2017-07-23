package wcai

type BaseModel struct {
	Id int `gorm:"primary_key"`
}

type OnlineProfile struct {
	BaseModel
	UserId    string
	Website   string
	Link      string
	Processed bool
}

func (OnlineProfile) TableName() string {
	return "wcai_profile_onlineprofile"
}

type GithubRepository struct {
	BaseModel
	Name               string
	Owner              string
	Stars              int
	Forks              int
	IsForked           bool
	Views              int
	Clones             int
	ForkProcessed      bool `gorm:"column:_fork_processed"`
	ViewsProcessed     bool
	ClonesProcessed    bool
	LanguagesProcessed bool `gorm:"column:_languages_processed"`
}

func (GithubRepository) TableName() string {
	return "wcai_profile_githubrepository"
}

type GithubRepositoryOwnership struct {
	BaseModel
	RepoId         int
	UserId         string
	GithubUsername string
	UserType       string
}

func (GithubRepositoryOwnership) TableName() string {
	return "wcai_profile_githubrepositoryownership"
}

type GithubLanguage struct {
	BaseModel
	RepoId int
	Name   string
	Count  int
}

func (GithubLanguage) TableName() string {
	return "wcai_profile_githublanguage"
}

type GithubTopic struct {
	BaseModel
	RepoId int
	Name   string
}

func (GithubTopic) TableName() string {
	return "wcai_profile_githubtopic"
}
