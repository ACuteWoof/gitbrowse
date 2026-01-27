package routes

import (
	"fmt"
	"net/http"

	"git.lewoof.xyz/clone/gitbrowse/config"
	"git.lewoof.xyz/clone/gitbrowse/template"
)

type RepoGitShowRoute struct {
	ConfigGetter func(repo string) config.PageConfig
}

func (route RepoGitShowRoute) Handler(w http.ResponseWriter, req *http.Request) {
	repo := req.PathValue("repo")
	hash := req.PathValue("hash")
	config := route.ConfigGetter(repo)

	fmt.Fprint(w, template.RepoGitShowPage{Commit: hash, Config: &config}.FullPage())
}
