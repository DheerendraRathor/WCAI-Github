package main

import (
    "github.com/google/go-github/github"
    "net/http"
    "fmt"

    // ORM setup
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"

    // Custom packages
    "./httpClient"
    "./wcai"
    "./config"
    //"context"
    "sync"
    "strings"
    "context"
)


var config configuration.WcaiConfiguration


func GetGithubClient() *github.Client {
    myHttpClient := &http.Client{
        Transport: &httpClient.GithubTransport{
            RoundTripper: http.DefaultTransport,
            ClientId: config.ClientId,
            ClientSecret: config.ClientSecret,
        },
    }
    return github.NewClient(myHttpClient)
}

func UpdateReposToDb(db *gorm.DB, user string, profile wcai.OnlineProfile, repositories []*github.Repository,
                    wg *sync.WaitGroup)  {

    defer func() {
        wg.Done()
    }()

    for _, repository := range repositories {
        db.Set("gorm:insert_option", "ON CONFLICT DO NOTHING").Create(&wcai.GithubRepository{
            BaseModel: wcai.BaseModel{Id: repository.GetID()},
            Name: repository.GetName(),
            Stars: repository.GetStargazersCount(),
            Forks: repository.GetForksCount(),
            IsForked: false,
            Views: 0,
            Clones: 0,
        })

        ownership := "owner"
        if strings.ToLower(*repository.Owner.Login) != strings.ToLower(user) {
            ownership = "member"
        }

        db.Set("gorm:insert_option", "ON CONFLICT DO NOTHING").Create(
            &wcai.GithubRepositoryOwnership{
                RepoId: repository.GetID(),
                UserId: profile.UserId,
                GithubUsername: user,
                UserType: ownership,
            },
        )

    }

    db.Model(&profile).Update("processed", true)
}

func FetchRepoAndPushToDb(db *gorm.DB, profile wcai.OnlineProfile, wg *sync.WaitGroup) {

    user, err := wcai.GetUserNameFromGithubUrl(profile.Link)
    if err != nil {
        fmt.Println(err)
        return
    }

    repositories, _ := wcai.GetRepositoriesForUser(*GetGithubClient(), user)

    wg.Add(1)
    go UpdateReposToDb(db, user, profile, repositories, wg)
}

func ListRepositories(db *gorm.DB, profiles []wcai.OnlineProfile)  {

    var wg sync.WaitGroup
    for _, profile := range profiles {
        FetchRepoAndPushToDb(db, profile, &wg)
    }
    wg.Wait()
}

func main()  {
    config = configuration.GetConfiguration()
    connString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=require password=%s",
        config.DbHost, config.DbUser, config.DbName, config.DbPassword)
    db, _ := gorm.Open("postgres", connString)
    db.DB().SetMaxOpenConns(16)
    defer db.Close()

    var githubProfiles []wcai.OnlineProfile

    db.Where("website = ? and processed = ?", "github", false).Find(&githubProfiles)

    ListRepositories(db, githubProfiles)

    githubClient := *GetGithubClient()
    limits, _, _ := githubClient.RateLimits(context.Background())

    fmt.Println(limits.Core.Limit, limits.Core.Remaining, limits.Core.Reset)

}
