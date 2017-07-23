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
	"gopkg.in/alecthomas/kingpin.v2"
)

var config configuration.WcaiConfiguration

var (
	updateFlag = kingpin.Flag("update", "Update mode").Default("nil").Enum(
		"nil", "fetch", "fork", "lang", "clones", "views", "limit",
	)
)

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
	kingpin.Parse()
	config = configuration.GetConfiguration()

	connString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=require password=%s",
		config.DbHost, config.DbUser, config.DbName, config.DbPassword)
	db, _ := gorm.Open("postgres", connString)
	db.DB().SetMaxOpenConns(16)
	defer db.Close()

	gc := *GetGithubClient()

	switch *updateFlag {
	case "fetch":
		wcai.ListRepositoriesAndPushToDb(gc, db)
	case "fork":
		wcai.UpdateRepositoryForkedStatusAndTopics(gc, db)
	case "lang":
		wcai.UpdateRepositoryLanguages(gc, db)
	case "clones":
		wcai.UpdateRepositoryClones(gc, db)
	case "views":
		wcai.UpdateRepositoryViews(gc, db)
	case "limit":
		limits := wcai.GetCoreRateLimits(gc)
		fmt.Println(limits.Limit, limits.Remaining, limits.Reset)
	default:
		kingpin.Usage()
	}
}
