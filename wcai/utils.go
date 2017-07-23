package wcai

import (
	"errors"
	"fmt"
	"regexp"
)

var githubUserNameRegex = `[a-z\d](?:[a-z\d]|-[a-z\d]){0,38}`

var githubUrlRegex = regexp.MustCompile(fmt.Sprintf(`(?i)(?:[^:]*://)?(?:www\.)?github\.(?:io|com)/(%s).*$`, githubUserNameRegex))

var githubIoRegex = regexp.MustCompile(fmt.Sprintf(`(?i)(?:[^:]*://)?(%s)[.@]github\.(?:io|com).*$`, githubUserNameRegex))

func GetUserNameFromGithubUrl(url string) (string, error) {
	matches := githubUrlRegex.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1], nil
	}

	matches = githubIoRegex.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", errors.New(fmt.Sprintf("%s is not a valid url", url))
}
