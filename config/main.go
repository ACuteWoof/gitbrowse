package config

import (
	"strings"
)

func GetIndexConfg() PageConfig {
	return IndexPageConfig
}

func GetRepoConfg(repo string) PageConfig {
	var NewPageConfig PageConfig = RepoPageConfig

	NewPageConfig.Title = strings.ReplaceAll(RepoPageConfig.Title, "<REPO>", repo)
	NewPageConfig.Description = strings.ReplaceAll(RepoPageConfig.Description, "<REPO>", repo)
	for i, style := range NewPageConfig.Styles {
		NewPageConfig.Styles[i] = strings.ReplaceAll(style, "<REPO>", repo)
	}
	NewPageConfig.Thumbnail = strings.ReplaceAll(RepoPageConfig.Thumbnail, "<REPO>", repo)
	NewPageConfig.Favicon = strings.ReplaceAll(RepoPageConfig.Favicon, "<REPO>", repo)
	NewPageConfig.RootDir = strings.ReplaceAll(RepoPageConfig.RootDir, "<REPO>", repo)
	return NewPageConfig
}
