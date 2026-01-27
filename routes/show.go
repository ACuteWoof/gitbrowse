package routes

import (
	"fmt"
	"net/http"

	"git.lewoof.xyz/gitbrowse/config"
	"git.lewoof.xyz/gitbrowse/template"
)

type RepoGitShowRoute struct {
	ConfigGetter func(repo string) config.PageConfig
}

func (route RepoGitShowRoute) Handler(w http.ResponseWriter, req *http.Request) {
	repo := req.PathValue("repo")
	hash := req.PathValue("hash")
	config := route.ConfigGetter(repo)

	fmt.Fprintf(w, template.RepoGitShowPage{Commit: hash, Config: &config}.FullPage())
}
