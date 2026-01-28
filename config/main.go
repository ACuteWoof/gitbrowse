// Gitbrowse: a simple web server for git.
// Copyright (C) 2026 Vithushan
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
		CloneURL:    strings.ReplaceAll(RepoPageConfig.CloneURL, "<REPO>", repo),
		Styles:      newstyles,
	}
	return NewPageConfig
}
