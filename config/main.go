package config

import (
	"strings"
)

func GetIndexConfg() PageConfig {
	return IndexPageConfig
}

func GetRepoConfg(repo string) PageConfig {
	var newstyles []string
	for _, style := range RepoPageConfig.Styles {
		newstyles = append(newstyles, strings.ReplaceAll(style, "<REPO>", repo))
	}

	var NewPageConfig PageConfig = PageConfig{
		Title:       strings.ReplaceAll(RepoPageConfig.Title, "<REPO>", repo),
		Description: strings.ReplaceAll(RepoPageConfig.Description, "<REPO>", repo),
		Thumbnail:   strings.ReplaceAll(RepoPageConfig.Thumbnail, "<REPO>", repo),
		Favicon:     strings.ReplaceAll(RepoPageConfig.Favicon, "<REPO>", repo),
		RootDir:     strings.ReplaceAll(RepoPageConfig.RootDir, "<REPO>", repo),
		URLRoot:     strings.ReplaceAll(RepoPageConfig.URLRoot, "<REPO>", repo),
		Styles:      newstyles,
	}
	return NewPageConfig
}
