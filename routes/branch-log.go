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

	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/template"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type RepoBranchLogRoute struct {
	ConfigGetter func(repo string) config.PageConfig
}

func (route RepoBranchLogRoute) Handler(w http.ResponseWriter, req *http.Request) {
	repo := req.PathValue("repo")
	branch := req.PathValue("branch")

	config := route.ConfigGetter(repo)

	r, err := git.PlainOpen(config.RootDir)
	errCheck(err)

	exists := false
	bs, err := r.Branches()
	errCheck(err)
	bs.ForEach(func(b *plumbing.Reference) error {
		if b.Name().Short() == branch {
			exists = true
			return nil
		}
		return nil
	})
	if !exists {
		http.Redirect(w, req, config.URLRoot + "/branch", http.StatusTemporaryRedirect)
		return
	}

	refName := plumbing.NewBranchReferenceName(branch)
	ref, err := r.Reference(refName, true)
	fmt.Fprint(w, template.RepoBranchLogPage{Repo: r, Branch: branch, BranchRef: ref, Config: &config}.FullPage())
}
