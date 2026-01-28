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
	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/template"
	"net/http"
	"os"
)

type IndexRoute struct {
	RepoRoute string;
	ConfigGetter func() config.PageConfig;
}

func (i IndexRoute) Handler(w http.ResponseWriter, _ *http.Request) {
	var dirs []string = getGitDirs(i.ConfigGetter().RootDir, i.RepoRoute)
	config := i.ConfigGetter()
	fmt.Fprint(w, template.IndexPage{Repos: dirs, Config: &config}.FullPage())
}

func getGitDirs(root string, displayRoot string) []string {
	d, err := os.ReadDir(root)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var dirs []string
	for _, dir := range d {
		if dir.IsDir() {
			dirs = append(dirs, displayRoot+dir.Name())
		}
	}
	return dirs
}
