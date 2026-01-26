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
	http.HandleFunc("/static/", HandleStatic)
	SetupSingleUserHandlers()
	http.ListenAndServe(":8088", nil)
}

func HandleStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func SetupSingleUserHandlers() {
	http.HandleFunc("/", routes.IndexRoute{DisplayRoot: "/", ConfigGetter: config.GetIndexConfg}.Handler)
	// http.HandleFunc("/{repo}", routes.RepoRoute{Config: &config.RepoPageConfig}.Handler)
	// http.HandleFunc("/{repo}/tags", routes.RepoTabRoute{Config: &config.RepoTabConfig}.Handler)
	// http.HandleFunc("/{repo}/diff/{commit1}/{commit2}", routes.RepoRoute{Config: &config.RepoTabFileConfig}.Handler)
	// http.HandleFunc("/{repo}/commit/{commit}/tree/{path...}", routes.RepoRoute{Config: &config.RepoTabFileConfig}.Handler)
	// http.HandleFunc("/{repo}/branch/{branch}/tree/{path...}", routes.RepoRoute{Config: &config.RepoTabFileConfig}.Handler)
	// http.HandleFunc("/{repo}/tag/{tag}", routes.RepoRoute{Config: &config.RepoTabFileConfig}.Handler)
}
