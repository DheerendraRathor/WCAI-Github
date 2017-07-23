package wcai

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

func GetRepositoriesForUser(client github.Client, user string) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{Type: "all"}
	var repos []*github.Repository
	var lastErr error = nil
	for {
		tempRepos, resp, err := client.Repositories.List(context.Background(), user, opt)
		if err != nil {
			fmt.Println(err)
			lastErr = err
			break
		}
		repos = append(repos, tempRepos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return repos, lastErr
}

func GetRepositoryForUser(client github.Client, user string, repo string) (*github.Repository, error) {
	repository, _, err := client.Repositories.Get(context.Background(), user, repo)
	return repository, err
}

func GetLanguagesForRepository(client github.Client, user string, repo string) (map[string]int, *github.Response, error) {
	return client.Repositories.ListLanguages(context.Background(), user, repo)
}

func IsRepositoryForked(repository *github.Repository) bool {
	return repository.Source == nil
}

func GetTopicsForRepository(repository *github.Repository) []string {
	return repository.Topics
}

func GetRepositoryCloneCount(client github.Client, user string, repo string) int {
	traffic, _, _ := client.Repositories.ListTrafficClones(context.Background(), user, repo, &github.TrafficBreakdownOptions{})
	return traffic.GetUniques()
}

func GetRepositoryViewsCount(client github.Client, user string, repo string) int {
	traffic, _, _ := client.Repositories.ListTrafficViews(context.Background(), user, repo, &github.TrafficBreakdownOptions{})
	return traffic.GetUniques()
}

func GetCoreRateLimits(client github.Client) *github.Rate {
	limits, _, _ := client.RateLimits(context.Background())
	return limits.Core
}
