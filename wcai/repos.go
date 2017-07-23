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

func GetLanguagesForRepository(client github.Client, user string, repo string) (map[string]int, *github.Response, error) {
	return client.Repositories.ListLanguages(context.Background(), user, repo)
}

func IsRepositoryForked(client github.Client, user string, repo string) bool {
	repository, _, _ := client.Repositories.Get(context.Background(), user, repo)
	return repository.Source == nil
}

func GetRepositoryCloneCount(client github.Client, user string, repo string) int {
	traffic, _, _ := client.Repositories.ListTrafficClones(context.Background(), user, repo, &github.TrafficBreakdownOptions{})
	return traffic.GetUniques()
}

func GetRepositoryViewsCount(client github.Client, user string, repo string) int {
	traffic, _, _ := client.Repositories.ListTrafficViews(context.Background(), user, repo, &github.TrafficBreakdownOptions{})
	return traffic.GetUniques()
}
