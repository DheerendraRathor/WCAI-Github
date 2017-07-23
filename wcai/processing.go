package wcai

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
)

func updateReposToDb(db *gorm.DB, user string, profile OnlineProfile, repositories []*github.Repository,
	wg *sync.WaitGroup) {

	defer func() {
		wg.Done()
	}()

	tx := db.Begin()

	for _, repository := range repositories {
		tx.Set("gorm:insert_option", "ON CONFLICT (id) DO UPDATE SET owner = EXCLUDED.owner").Create(&GithubRepository{
			BaseModel: BaseModel{Id: repository.GetID()},
			Name:      repository.GetName(),
			Owner:     *repository.Owner.Login,
			Stars:     repository.GetStargazersCount(),
			Forks:     repository.GetForksCount(),
			IsForked:  false,
			Views:     0,
			Clones:    0,
		})

		ownership := "owner"
		if strings.ToLower(*repository.Owner.Login) != strings.ToLower(user) {
			ownership = "member"
		}

		tx.Set("gorm:insert_option", "ON CONFLICT DO NOTHING").Create(
			&GithubRepositoryOwnership{
				RepoId:         repository.GetID(),
				UserId:         profile.UserId,
				GithubUsername: user,
				UserType:       ownership,
			},
		)

	}

	tx.Model(&profile).Update("processed", true)
	tx.Commit()
}

func fetchRepoAndPushToDb(gc github.Client, db *gorm.DB, profile OnlineProfile, wg *sync.WaitGroup) {

	user, err := GetUserNameFromGithubUrl(profile.Link)
	if err != nil {
		fmt.Println(err)
		return
	}

	repositories, _ := GetRepositoriesForUser(gc, user)

	wg.Add(1)
	go updateReposToDb(db, user, profile, repositories, wg)
}

func ListRepositoriesAndPushToDb(gc github.Client, db *gorm.DB) {

	var profiles []OnlineProfile

	db.Where("website = ? and processed = ?", "github", true).Find(&profiles)

	var wg sync.WaitGroup
	for _, profile := range profiles {
		fetchRepoAndPushToDb(gc, db, profile, &wg)
	}
	wg.Wait()
}

func updateRepoForkStatusAndTopics(repo GithubRepository, fullRepo *github.Repository, wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done()
	repo.IsForked = IsRepositoryForked(fullRepo)
	topics := fullRepo.Topics

	tx := db.Begin()
	updates := map[string]interface{}{
		"is_forked":       repo.IsForked,
		"_fork_processed": true,
	}

	tx.Model(&repo).Updates(updates)

	db.Delete(GithubTopic{}, "repo_id = ?", repo.Id)

	for _, topic := range topics {
		tx.Create(&GithubTopic{
			Name:   topic,
			RepoId: repo.Id,
		})
	}
	tx.Commit()
	fmt.Printf("Done for repo: %s/%s\n", repo.Owner, repo.Name)
}

func UpdateRepositoryForkedStatusAndTopics(gc github.Client, db *gorm.DB) {
	var repositories []GithubRepository
	limits := GetCoreRateLimits(gc)
	fmt.Println("Remaining Limit:", limits.Remaining)

	db.Where("owner <> ? and _fork_processed = ?", "", false).Limit(limits.Remaining).Find(&repositories)

	var wg sync.WaitGroup
	for _, repo := range repositories {
		fullRepo, err := GetRepositoryForUser(gc, repo.Owner, repo.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		wg.Add(1)
		go updateRepoForkStatusAndTopics(repo, fullRepo, &wg, db)
	}

	wg.Wait()
}

func updateRepositoryLanguages(db *gorm.DB, repo GithubRepository, languages map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()

	tx := db.Begin()
	for name, value := range languages {
		tx.Create(&GithubLanguage{
			Name:   name,
			Count:  value,
			RepoId: repo.Id,
		})
	}

	tx.Commit()
	fmt.Printf("Done for repo: %s/%s\n", repo.Owner, repo.Name)
}

func UpdateRepositoryLanguages(gc github.Client, db *gorm.DB) {
	var repositories []GithubRepository
	limits := GetCoreRateLimits(gc)
	fmt.Println("Remaining Limit:", limits.Remaining)

	db.Where("owner <> ? and _languages_processed = ?", "", false).Limit(limits.Remaining).Find(&repositories)

	var wg sync.WaitGroup
	for _, repo := range repositories {
		languages, _, err := GetLanguagesForRepository(gc, repo.Owner, repo.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}

		wg.Add(1)
		go updateRepositoryLanguages(db, repo, languages, &wg)
	}

}
