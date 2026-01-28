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

package routes

import (
	"fmt"
	"net/http"
	"strings"

	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/template"
	"github.com/go-git/go-git/v6"
	"html"
)

type RepoReadmeRoute struct {
	ConfigGetter func(repo string) config.PageConfig
}

func (route RepoReadmeRoute) Handler(w http.ResponseWriter, req *http.Request) {
	repo := req.PathValue("repo")
	config := route.ConfigGetter(repo)

	r, err := git.PlainOpen(config.RootDir)
	errCheck(err)

	ref, err := r.Head()
	errCheck(err)

	commit, err := r.CommitObject(ref.Hash())
	errCheck(err)

	tree, err := commit.Tree()
	errCheck(err)

	var infoFiles []template.RepoInfoFile

	possibleInfoFiles := []string{"README.md", "README.txt", "README", "readme.txt", "readme.md", "readme", "NOTICE", "NOTICE.txt", "NOTICE.md", "notice", "notice.txt", "notice.md", "LICENSE", "LICENSE.md", "LICENSE.txt", "license", "license.txt", "license.md", "THIRD_PARTY_LICENSES.txt", "THIRD_PARTY_LICENSES.md", "THIRD_PARTY_LICENSES", "third_party_licenses.txt", "third_party_licenses.md", "third_party_licenses"}
	for _, possibleInfoFile := range possibleInfoFiles {
		file, err := tree.File(possibleInfoFile)
		if err == nil {
			content, err := file.Contents()
			errCheck(err)
			if strings.HasSuffix(possibleInfoFile, ".md") {
				infoFiles = append(infoFiles, template.RepoInfoFile{possibleInfoFile, markdownToHtml(content)})
			} else {
				infoFiles = append(infoFiles, template.RepoInfoFile{possibleInfoFile, html.EscapeString(content)})
			}
		}
	}
	if len(infoFiles) == 0 {
		http.Redirect(w, req, req.URL.Path+"/branch", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprint(w, template.RepoInfoPage{InfoFiles: infoFiles, Config: &config}.FullPage())
}
