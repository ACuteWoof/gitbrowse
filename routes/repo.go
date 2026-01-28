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
	"github.com/go-git/go-git/v6/plumbing/object"
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

	// file, err := tree.File("README.md")
	possibleReadmes := []string{"README.md", "README.markdown", "README.txt", "README", "readme", "README", "LICENSE", "LICENSE.md", "LICENSE.txt", "license", "license.txt", "license.md"}
	var readme *object.File
	for _, possibleReadme := range possibleReadmes {
		file, err := tree.File(possibleReadme)
		if err == nil {
			readme = file
			break
		}
	}
	if readme == nil {
		http.Redirect(w, req, req.URL.Path+"/branch", http.StatusTemporaryRedirect)
		return
	}

	content, err := readme.Contents()
	errCheck(err)

	var readmeHtml string
	if strings.HasSuffix(readme.Name, ".md") {
		readmeHtml = markdownToHtml(content)
	} else {
		readmeHtml = "<pre>" + html.EscapeString(string(content)) + "</pre>"
	}

	fmt.Fprint(w, template.RepoReadmePage{Readme: readmeHtml, Config: &config}.FullPage())
}
