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
	Name     string
	Stars    int
	Forks    int
	IsForked bool
	Views    int
	Clones   int
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
	RepoId uint
	Name   string
	Count  uint `gorm:"default:0"`
}

func (GithubLanguage) TableName() string {
	return "wcai_profile_githublanguage"
}

type GithubTopic struct {
	BaseModel
	RepoId uint
	UserId uint
	Name   string
}

func (GithubTopic) TableName() string {
	return "wcai_profile_githubtopic"
}
