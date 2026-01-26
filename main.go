package main

import (
	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/routes"
	"net/http"
	// "os"
)

func main() {
	startHttpServer()
}

func startHttpServer() {
	SetupSingleUserHandlers()
	http.HandleFunc("/static/", HandleStatic)
	http.ListenAndServe(":8088", nil)
}

func HandleStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func SetupSingleUserHandlers() {
	http.HandleFunc("/", routes.IndexRoute{RepoRoute: "/browse/", ConfigGetter: config.GetIndexConfg}.Handler)
	http.HandleFunc("/browse/{repo}/", routes.RepoReadmeRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	http.HandleFunc("/browse/{repo}/branch/", routes.RepoBranchesRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	http.HandleFunc("/browse/{repo}/branch/{branch}/tree/{filepath...}", routes.RepoBranchTreeRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	http.HandleFunc("/browse/{repo}/branch/{branch}/commit", routes.RepoBranchLogRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	http.HandleFunc("/browse/{repo}/tag/", routes.RepoTagsRoute{ConfigGetter: config.GetRepoConfg}.Handler)
	// http.HandleFunc("/browse/{repo}/commit/{commit}/", routes.RepoCommitRoute{ConfigGetter: config.GetRepoConfg}.Handler)
}
