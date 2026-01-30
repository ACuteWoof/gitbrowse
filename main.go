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

package main

import (
	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/routes"
	"net/http"
	"strings"
	// "os"
)

func main() {
	startHttpServer()
}

func startHttpServer() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	staticMux := http.NewServeMux()
	staticMux.Handle("/", fs)

	singleUserHandler(mux) // change to setupMultiUserHandlers(mux) to enable multi user support

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/static")
			staticMux.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/favicon.ico" {
			http.ServeFile(w, r, "./static/favicon.ico")
			return
		}

		mux.ServeHTTP(w, r)
	})

	http.ListenAndServe(":8088", nil)
}

func HandleStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func singleUserHandler(mux *http.ServeMux) {
	mux.HandleFunc("/", routes.IndexRoute{RepoRoute: "/", ConfigGetter: config.GetIndexConfg}.Handler)
	mux.HandleFunc("/{repo}/", routes.RepoReadmeRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	mux.HandleFunc("/{repo}/branch/", routes.RepoBranchesRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	mux.HandleFunc("/{repo}/branch/{branch}/tree/{filepath...}", routes.RepoBranchTreeRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	mux.HandleFunc("/{repo}/branch/{branch}/commit", routes.RepoBranchLogRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	mux.HandleFunc("/{repo}/tag/", routes.RepoTagsRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	mux.HandleFunc("/{repo}/tag/{name}/{fileName}", routes.TagDownloadRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	mux.HandleFunc("/{repo}/show/{hash}", routes.RepoGitShowRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	mux.HandleFunc("/{repo}/grep/", routes.RepoGrepRoute{ConfigGetter: config.GetRepoConfg}.Handler)
}
