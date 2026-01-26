package config

import (
	"strings"
)

func GetIndexConfg() PageConfig {
	return IndexPageConfig
}

func GetRepoConfg(repo string) PageConfig {
	RepoPageConfig.Title = strings.ReplaceAll(RepoPageConfig.Title, "<REPO>", repo)
	RepoPageConfig.Description = strings.ReplaceAll(RepoPageConfig.Description, "<REPO>", repo)
	for i, style := range RepoPageConfig.Styles {
		RepoPageConfig.Styles[i] = strings.ReplaceAll(style, "<REPO>", repo)
	}
	RepoPageConfig.Thumbnail = strings.ReplaceAll(RepoPageConfig.Thumbnail, "<REPO>", repo)
	RepoPageConfig.Favicon = strings.ReplaceAll(RepoPageConfig.Favicon, "<REPO>", repo)
	return RepoPageConfig
}
