package main

import (
	"fmt"
	"net/http"

	"github.com/DheerendraRathor/WCAI-github/config"
	"github.com/DheerendraRathor/WCAI-github/httpClient"
	"github.com/DheerendraRathor/WCAI-github/wcai"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var config configuration.WcaiConfiguration

func GetGithubClient() *github.Client {
	myHttpClient := &http.Client{
		Transport: &httpClient.GithubTransport{
			RoundTripper: http.DefaultTransport,
			ClientId:     config.ClientId,
			ClientSecret: config.ClientSecret,
		},
	}
	return github.NewClient(myHttpClient)
}

func main() {
	config = configuration.GetConfiguration()

	connString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=require password=%s",
		config.DbHost, config.DbUser, config.DbName, config.DbPassword)
	db, _ := gorm.Open("postgres", connString)
	db.DB().SetMaxOpenConns(16)
	defer db.Close()

	gc := *GetGithubClient()

	//wcai.ListRepositoriesAndPushToDb(gc, db)
	//wcai.UpdateRepositoryForkedStatusAndTopics(gc, db)
	//wcai.UpdateRepositoryLanguages(gc, db)

	limits := wcai.GetCoreRateLimits(gc)

	fmt.Println(limits.Limit, limits.Remaining, limits.Reset)

}
